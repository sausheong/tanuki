# Tanuki

Tanuki is a polyglot web framework that allows developers to write web applications and services in multiple programming languages. The basic concept is simple -- an acceptor (which is a HTTP server) accepts HTTP requests and redirects it to either an executable binary file (or script) or a TCP socket server to handle. These handlers can be written in different programming languages and run independently from each other.

Tanuki is experimental software at this point in time. Use at your own risk.

## Installing Tanuki

Tanuki is written in Go. To install it in your platform, you can install [Go](https://golang.org) and get the source:

```
go get github.com/sausheong/tanuki
```

There will be downloadable binaries for your platform at a later date.

## How to use Tanuki

Once you can have downloaded the source, you can build the command line tool `tanuki`.

```
go build
```

With tha command line tool, can you create a skeleton structure for your new web application or service:

```
./tanuki create poko
```

This will create a new directory named `poko` with the necessary struct and files. To start your new Tanuki application, go to your new application directory and run this:

```
./tanuki start
```

This will start the Tanuki application at port `8080`. You can port number or the IP address or hostname as well, just follow the instructions in the `tanuki` command line tool.

If you have executable binaries or scripts in the correct format in your application `/bin` directory, they will be loaded and you can call them immediately. If you have listeners in the correct format in your application `/listeners` directory they will be started as well by Tanuki.

The file name format for bins and listeners are:

```
<METHOD>__<PATH>
```

The forward slashes in the path `/` will be converted to double underscores `__`. Let's say you have the a Ruby executable script `bin/get__hello__ruby` (notice the double underscores), this means you can open up the browser and go to the URL `/_/hello/ruby` and this will execute the script. Remember, all Tanuki actions start with `/_/`.

## Tanuki actions

### Tanuki HTTP requests

The only input into the handlers (either bins or listeners) is the HTTP request JSON. The below is the Go structs for the JSON.

```go
// RequestInfo corresponds to a HTTP 1.1 request
type RequestInfo struct {
	Method           string                 `json:"Method"`
	URL              URLInfo                `json:"URL"`
	Proto            string                 `json:"Proto"`
	Header           map[string][]string    `json:"Header"`
	Body             string                 `json:"Body"`
	ContentLength    int64                  `json:"ContentLength"`
	TransferEncoding []string               `json:"TransferEncoding"`
	Host             string                 `json:"Host"`
	Params           map[string][]string    `json:"Params"`
	Multipart        map[string][]Multipart `json:"Multipart"`
	RemoteAddr       string                 `json:"RemoteAddr"`
	RequestURI       string                 `json:"RequestURI"`
}

// Multipart corresponds to a multi-part file
type Multipart struct {
	Filename    string `json:"Filename"`
	ContentType string `json:"ContentType"`
	Content     string `json:"Content"`
}

// URLInfo corresponds to a URL
type URLInfo struct {
	Scheme   string `json:"Scheme"`
	Opaque   string `json:"Opaque"`
	Host     string `json:"Host"`
	Path     string `json:"Path"`
	RawQuery string `json:"RawQuery"`
	Fragment string `json:"Fragment"`
}
```

Tanuki provides a `Params` field, which is a hash (or dictionary) with a string as the key and an array of strings as the value. The `Params` field is a convenience for Tanuki developers and contains all the parameters sent by the client, including those in the URL or in the case of a `POST`, also in the forms. 

This is an example of the JSON for the HTTP request, a GET request to Tanuki at the URL `/_/hello/world`. 

```json
{
  "Method": "GET",
  "URL": {
    "Scheme": "",
    "Opaque": "",
    "Host": "",
    "Path": "/_/hello/world",
    "RawQuery": "name=sausheong",
    "Fragment": ""
  },
  "Proto": "HTTP/1.1",
  "Header": {
    "Accept": [
      "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
    ],
    "Accept-Encoding": [
      "gzip, deflate"
    ],
    "Accept-Language": [
      "en-sg"
    ],
    "Connection": [
      "keep-alive"
    ],
    "Cookie": [
      "hello=world; mykey=myvalue"
    ],
    "Upgrade-Insecure-Requests": [
      "1"
    ],
    "User-Agent": [
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.2 Safari/605.1.15"
    ]
  },
  "Body": "",
  "ContentLength": 0,
  "TransferEncoding": null,
  "Host": "localhost:8080",
  "Params": {
    "name": [
      "sausheong"
    ]
  },
  "Multipart": {},
  "RemoteAddr": "[::1]:52340",
  "RequestURI": "/_/hello/world?name=sausheong"
}
```

You might notice here that if you want to get the `name` you can either parse the `RawQuery` or get `RequestURI` or you can simply get it from `Params`. 

### Tanuki HTTP response

The only response that Tanuki accepts from the handlers is a HTTP response JSON. Here's the Go struct for the JSON.

```go
// ResponseInfo corresponds to a HTTP 1.1 response
type ResponseInfo struct {
	Status int                 `json:"status"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}
```



## Why Tanuki?

It seems a bit of work to write a single web application in different programming languages, so why do it? It's really all about future-proofing.

### No more fighting

### Easy upgrade

### Switching technologies

### Modular replacement

### Easy testing


## What's the trade-off?