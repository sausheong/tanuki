package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sausheong/tanuki/data"
	"github.com/spf13/cobra"
)

var handlers Handlers

var port *int
var ip *net.IP
var static *string
var readTimeout *int64
var writeTimeout *int64
var handlerConfig *string

func init() {
	rootCmd.AddCommand(acceptorCmd)
	port = acceptorCmd.Flags().Int("port", 8080, "server port number")
	ip = acceptorCmd.Flags().IP("host", net.IPv4(0, 0, 0, 0), "host IP address")
	static = acceptorCmd.Flags().String("static", "static", "directory for static files")
	readTimeout = acceptorCmd.Flags().Int64("readtimeout", 10, "server read time-out")
	writeTimeout = acceptorCmd.Flags().Int64("writetimeout", 600, "server write time-out")
	handlerConfig = acceptorCmd.Flags().String("handlerConfig", "handlers.yaml", "handler configuration file")
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
	var err error
	handlers, err = getHandlers(*handlerConfig)
	if err != nil {
		danger("Cannot load handlers", err)
		fmt.Println("Cannot load handlers configuration, please check the handlers.yaml file", err)
		return
	}
	fmt.Println("handlers:", handlers)

	startLocalListeners()
	router := httprouter.New()

	// all routes to the accept func, developer should differentiate at the handler level
	router.GET("/_/*p", accept)
	router.POST("/_/*p", accept)
	router.PUT("/_/*p", accept)
	router.DELETE("/_/*p", accept)
	router.PATCH("/_/*p", accept)
	router.HEAD("/_/*p", accept)
	router.OPTIONS("/_/*p", accept)

	router.ServeFiles("/_s/*filepath", http.Dir(*static))

	host := join((*ip).String(), ":", strconv.Itoa(*port))
	server := &http.Server{
		Addr:           host,
		Handler:        router,
		ReadTimeout:    time.Duration(*readTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(*writeTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Tanuki started at", host, time.Now().String())
	server.ListenAndServe()
}

// performs the main processing for the acceptor
func accept(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Print("Tanuki accepting ", request.Method, " request ", request.URL, " - ")
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
	info("Request - ", string(reqJSON))
	// ------------
	// send request
	// ------------
	var output []byte

	// hand off to the correct handlers
	if handler, ok := handlers.getHandler(strings.ToLower(request.Method), request.URL.Path); ok {
		switch handlerType := handler.Type; handlerType {
		case "bin":
			// bins are executable binary files. Tanuki walks through the bin/ directory to look for
			// executable binaries and adds them to a list. Each binary takes in a request JSON through the
			// command line argument and returns a response JSON through STDOUT. When a route matches the
			// name of the binary, Tanuki will call the binary with the request JSON as the argument and
			// parses the return output as a response JSON
			output, err = exec.Command(handler.Path, string(reqJSON)).Output()
			if err != nil {
				danger("Cannot execute bin", err)
			}
			info("Binary called", request.Method, join("(", handler.Path, ") - ", time.Since(start).String()))

		case "listener":
			// listeners are TCP socket servers. Tanuki walks through files in the listener/ directory, and starts
			// each listener, adding it to a hash, with the routeID of the listener as the key and the port number
			// of the TCP server as the value. Each listener takes in a request JSON (terminated by a newline \n)
			// and returns a response JSON through the same connection. Listeners are supposed to be more performant
			// because they can be multi-threaded and also already started up, unlike binaries
			var conn net.Conn
			if handler.Local {
				conn, err = net.Dial("tcp", ":"+handler.Port)
				if err != nil {
					reply(writer, 404, []byte("Cannot connect to listener"))
					danger("Cannot connect to local listener", err)
					fmt.Println(time.Since(start).String())
					return
				}
			} else {
				conn, err = net.Dial("tcp", handler.Path)
				if err != nil {
					reply(writer, 404, []byte("Cannot connect to listener"))
					danger("Cannot connect to remote listener", err)
					fmt.Println(time.Since(start).String())
					return
				}
			}
			fmt.Fprintf(conn, string(reqJSON)+"\n")
			// listen for reply
			output, err = bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				fmt.Println("Cannot read from listener", err)
			}
			info("Listener called", request.Method, join(handler.Path, time.Since(start).String()))

		default:
			reply(writer, 404, []byte("Handler type not found"))
			info("Handler type not found", request.Method, request.URL.Path, handler.Path)
			fmt.Println(time.Since(start).String())
			return
		}

	} else {
		reply(writer, 404, []byte("Tanuki action not found"))
		info("Action not found", request.Method, request.URL.Path, handler.Path)
		fmt.Println(time.Since(start).String())
		return
	}

	// ----------------
	// receive response
	// ----------------
	// parse the JSON output
	var response data.ResponseInfo
	info("Response:", string(output))
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		reply(writer, 500, []byte("Cannot unmarshal response JSON - "+err.Error()))
		danger("Cannot unmarshal response JSON", err, request.Method, request.URL.Path, " - ", time.Since(start).String())
		danger(string(output))
		fmt.Println(time.Since(start).String())
		return
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
	fmt.Println(time.Since(start).String())
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

func startLocalListeners() {
	for i, handler := range handlers {
		if handler.Type == "listener" {
			if handler.Local {
				// get a free port for the listener
				port, err := getFreePort()
				if err != nil {
					fmt.Println("Cannot get port", err)
				}
				// start the listener and pass it the port number
				go exec.Command(handler.Path, strconv.Itoa(port)).Run()
				handlers[i].Port = strconv.Itoa(port)
				fmt.Println("handler started:", handler.Path, port)
			} else {
				_, err := net.Dial("tcp", handler.Path)
				if err != nil {
					danger("Cannot connect to", handler.Path, err)
				}
			}
		}
	}
}
