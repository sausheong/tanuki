package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sausheong/tanuki/data"
	"github.com/spf13/cobra"
)

var capturePort *int
var outputFile *string

func init() {
	rootCmd.AddCommand(captureCmd)
	capturePort = captureCmd.Flags().Int("port", 8088, "capture server port number")
	outputFile = captureCmd.Flags().String("output", "request.json", "JSON request captured to file name")
}

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture a HTTP request into a JSON file for testing",
	Long:  `Capture a HTTP request from the browser into a file to be used for testing and debugging. This command starts a HTTP server (defaults at 8088). Use Postman or your browser to send a request to this server to capture it to a JSON file`,
	Run: func(cmd *cobra.Command, args []string) {
		capture()
	},
}

func capture() {
	router := httprouter.New()

	// all routes to the accept func, developer should differentiate at the handler level
	router.GET("/*p", processCapture)
	router.POST("/*p", processCapture)
	router.PUT("/*p", processCapture)
	router.DELETE("/*p", processCapture)
	router.PATCH("/*p", processCapture)
	router.HEAD("/*p", processCapture)
	router.OPTIONS("/*p", processCapture)

	host := join(":", strconv.Itoa(*capturePort))
	server := &http.Server{
		Addr:           host,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Tanuki Capture Server started at", host, time.Now().String())
	server.ListenAndServe()
}

func processCapture(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Print("Tanuki Capture Server accepting ", request.Method, " request ", request.URL, " - ")
	start := time.Now()
	// the multipart contains the multipart data
	multipart := make(map[string][]data.Multipart)

	// parse the multipart form for stuff in the forms if it's a POST
	if request.Method == "POST" {
		request.ParseMultipartForm(3 << 20)
		if request.MultipartForm != nil {
			for mk, mv := range request.MultipartForm.File {
				var parts []data.Multipart
				for _, v := range mv {
					f, err := v.Open()
					if err != nil {
						danger("Cannot read multipart message", err)
					}
					var buf bytes.Buffer
					_, err = io.Copy(&buf, f)
					if err != nil {
						danger("Cannot copy multipart message into buffer", err)
					}
					content := base64.StdEncoding.EncodeToString(buf.Bytes())
					part := data.Multipart{
						Filename:    v.Filename,
						ContentType: v.Header["Content-Type"][0],
						Content:     content,
					}
					parts = append(parts, part)
				}
				multipart[mk] = parts
			}
		}
	}

	// the form contains data from the URL as well as the POST form
	params := make(map[string][]string)
	err := request.ParseForm()
	if err != nil {
		danger("Failed to parse form", err)
	}

	for fk, fv := range request.Form {
		params[fk] = fv
	}

	// create the struct for the JSON
	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	reqInfo := data.RequestInfo{
		Method: request.Method,
		URL: data.URLInfo{
			Scheme:   request.URL.Scheme,
			Opaque:   request.URL.Opaque,
			Host:     request.URL.Host,
			Path:     request.URL.Path,
			RawQuery: request.URL.RawQuery,
			Fragment: request.URL.Fragment,
		},
		Proto:            request.Proto,
		Header:           request.Header,
		Body:             buf.String(),
		ContentLength:    request.ContentLength,
		TransferEncoding: request.TransferEncoding,
		Host:             request.Host,
		Params:           params,
		Multipart:        multipart,
		RemoteAddr:       request.RemoteAddr,
		RequestURI:       request.RequestURI,
	}
	// marshal the RequestInfo struct into JSON
	reqJSON, err := json.Marshal(reqInfo)
	if err != nil {
		danger("Failed to marshal the request into JSON - ", err)
	}

	err = ioutil.WriteFile(*outputFile, []byte(reqJSON), 0644)
	fmt.Println(time.Since(start).String())

	if err != nil {
		fmt.Println("Cannot capture request - ", err)
		reply(writer, 200, []byte("Request capture failed"))
	}

	// respond to the client
	reply(writer, 200, []byte("Request captured!"))
}
