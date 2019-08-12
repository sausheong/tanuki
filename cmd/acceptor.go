package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sausheong/tanuki/structs"
	"github.com/spf13/cobra"
)

var bins []string
var port *int
var ip *net.IP
var static *string
var readTimeout *int64
var writeTimeout *int64

func init() {
	rootCmd.AddCommand(acceptorCmd)
	port = acceptorCmd.Flags().Int("port", 8080, "server port number")
	ip = acceptorCmd.Flags().IP("host", net.IPv4(0, 0, 0, 0), "host IP address")
	static = acceptorCmd.Flags().String("static", "static", "directory for static files")
	readTimeout = acceptorCmd.Flags().Int64("readtimeout", 10, "server read time-out")
	writeTimeout = acceptorCmd.Flags().Int64("writetimeout", 600, "server write time-out")
}

var acceptorCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Tanuki acceptor",
	Long: `The Tanuki acceptor receives all HTTP requests to the web application. This command starts the acceptor. Run 
this command only in the Tanuki application root.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	getAllBins()
	router := httprouter.New()

	// currently supports GET and POST only
	// TODO
	router.GET("/_/*p", accept)
	router.POST("/_/*p", accept)

	router.ServeFiles("/_s/*filepath", http.Dir(*static))

	host := join((*ip).String(), ":", strconv.Itoa(*port))
	server := &http.Server{
		Addr:           host,
		Handler:        router,
		ReadTimeout:    time.Duration(*readTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(*writeTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Tanuki started at", host)
	server.ListenAndServe()
}

// performs the main processing for the acceptor
func accept(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Println("Tanuki accepting", request.Method, "request", request.URL)
	// the multipart contains the multipart data
	multipart := make(map[string][]structs.Multipart)

	// parse the multipart form for stuff in the forms if it's a POST
	if request.Method == "POST" {
		request.ParseMultipartForm(3 << 20)
		if request.MultipartForm != nil {
			for mk, mv := range request.MultipartForm.File {
				var parts []structs.Multipart
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
					part := structs.Multipart{
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
	reqInfo := structs.RequestInfo{
		Method: request.Method,
		URL: structs.URLInfo{
			Scheme:   request.URL.Scheme,
			Opaque:   request.URL.Opaque,
			Host:     request.URL.Host,
			Path:     request.URL.Path,
			RawQuery: request.URL.RawQuery,
			Fragment: request.URL.Fragment,
		},
		Proto:            request.Proto,
		Header:           request.Header,
		Body:             request.Body,
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
		danger("Failed to marshal the request into JSON", err)
	}
	// routeID is used to identify which responder to call
	routeID := join(strings.ToLower(request.Method), strings.ReplaceAll(request.URL.Path[2:], "/", "__"))

	// send request
	var output []byte
	// if if it's in the bins, run it
	if exists(bins, routeID) {
		// execute the bin and get a JSON output
		output, err = exec.Command(join("bin/", routeID), string(reqJSON)).Output()
		if err != nil {
			danger("Cannot execute bin", err)
		}
		info("Action called", request.Method, request.URL.Path, join("(", routeID, ")"))
	} else {
		reply(writer, 404, []byte("Tanuki action not found"))
		info("Action not found", request.Method, request.URL.Path, join("(", routeID, ")"))
		return
	}
	// parse the JSON output
	var response structs.ResponseInfo
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		reply(writer, 500, []byte("Cannot unmarshal response JSON - "+err.Error()))
		danger("Cannot unmarshal response JSON", err, request.Method, request.URL.Path)
	}

	// write headers to writer
	for k, v := range response.Header {
		for _, val := range v {
			writer.Header().Add(k, val)
		}
	}

	// see if we need to decode the body first
	var data []byte
	// get content type
	ctype, hasCType := response.Header["Content-Type"]
	if hasCType == true {
		if isTextMimeType(ctype[0]) {
			data = []byte(response.Body)
		} else {
			data, _ = base64.StdEncoding.DecodeString(response.Body)
		}
	} else {
		data = []byte(response.Body) // if not given the content type, assume it's text
	}

	// respond to the client
	reply(writer, response.Status, data)

}

// send response to client
func reply(writer http.ResponseWriter, status int, body []byte) {
	writer.WriteHeader(status)
	writer.Write(body)
}

func isTextMimeType(ctype string) bool {
	if strings.HasPrefix(ctype, "text") ||
		strings.HasPrefix(ctype, "application/json") {
		return true
	}
	return false
}

func getBin(method, path string) (binPath string) {
	binPath = join(strings.ToUpper(method), "__", strings.ReplaceAll(path, "/", "__"))

	return
}

// load all bins into the bins variable
func getAllBins() {
	err := filepath.Walk("bin",
		func(path string, info os.FileInfo, err error) error {
			// not a directory
			if !info.IsDir() {
				// must be an executable file
				if info.Mode()&0100 == os.FileMode(0000100) {
					bins = append(bins, info.Name())
				}
			}

			return nil
		})
	if err != nil {
		danger("Cannot load bins", err)
	}
}
