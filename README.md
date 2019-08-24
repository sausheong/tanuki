<div id="logo" align="center">
    <img src="tanuki.png" alt="Tanuki" width="150"/>
</div>

# Tanuki

Tanuki is a polyglot web framework that allows developers to write a web application or service in multiple programming languages. The basic concept is simple -- an acceptor (which is a HTTP server) receives HTTP requests and redirects it to a handler (either an executable binary/script or a TCP socket server) to handle. These handlers can be written in different programming languages and run independently from each other.

This is :warning: _experimental_ software at this point in time. Use at your own risk! :warning:

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

This will start the Tanuki application at port `8080`. You can change the port number or the IP address or hostname as well, just follow the instructions in the `tanuki` command line tool.

Before you can run any Tanuki application you need to set up the handlers.

## Handlers

There are two types of Tanuki handlers:

### Executable binary or script file

Executable binaries or scripts (or _bins_) are files that are, well, executable. You can normally run them on the command line. As a Tanuki bin handler, it must be able to take process a JSON string (as a command line argument) and return another JSON string to STDOUT. The input JSON contains HTTP request information and the output JSON has HTTP response information.

### TCP socket server

A listener is Tanuki handler that runs as a TCP socket servers. There are two types of listeners:

1. A _local_ listener, which is started by Tanuki and runs in the same host as Tanuki
2. A _remote_ listener, which is _not_ started by Tanuki and Tanuki assumes it is reachable (i.e. it is your responsibility to make sure it's reachable)

Remote listeners allow Tanuki to be distributed to multiple hosts.

### Handler configuration

Handlers are configured in a `handlers.yaml` file by default. This is how a handler configuration file looks like:

```yaml
--- # handlers
- method : get
  route  : /_/hello/ruby
  type   : bin
  path   : bin/hello-ruby

- method : get
  route  : /_/hello/go
  type   : bin
  path   : bin/hello-go

- method : get
  route  : /_/hello/ruby/listener
  type   : listener
  local  : true 
  path   : listeners/hello-ruby

- method : get
  route  : /_/hello/go/listener
  type   : listener
  local  : false
  path   : localhost:55771
```
The configuration is pretty straightforward. The `method` parameter determines the type of HTTP method to be routed, the `route` determines the URL route, the `type` is the type of handler, either a /bin/ or a /listener/ and the `path` is either a path to the file, or a URL of the remote listener, which includes the hostname and the port. Finally, the `local` parameter determines if it’s a local or a remote listener.

## Handler input and output

### HTTP requests

The only input into the handlers (either bins or listeners) is the HTTP request JSON. The below is the Go structs for the JSON (I'm showing you the Go struct because Tanuki is written in Go).

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

### HTTP response

The only response that Tanuki accepts from the handlers is a HTTP response JSON. Here's the Go struct for the JSON.

```go
// ResponseInfo corresponds to a HTTP 1.1 response
type ResponseInfo struct {
	Status int                 `json:"status"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}
```

Here's an example of the response JSON to be sent back to Tanuki.

```json
{
  "status": 200,
  "header": {
    "Content-Length": ["15"],
    "Content-Type": ["text/plain; charset=utf-8"]
  },
  "body": "hello sausheong"
}
```
## Examples

Let's look at some examples. They are pretty simple and almost trivial and does one thing. When the user browses the URL:

```
https://localhost:8081/_/hello/world?name=sausheong
```

The handlers should return:

```
hello sausheong
```


### Executable binaries or scripts

Bins are the simplest handler to write and can be written in any programming language than can parse JSON and take in a command line argument and print a JSON string to STDOUT. This is an example of a simple bin handler written in Ruby.

```ruby
#!/usr/bin/env ruby
require 'json'

request = JSON.parse ARGV[0]
response = {
    status: 200,
    header: {},
    body: "hello #{request['Params']['name'][0]}"
}
puts response.to_json
```

The code above shows how the first command line argument to the script `ARGV[0]` is parsed into JSON, and the parameters are used in the body. The response is created first as a hash and converted into JSON before being printed out to STDOUT.

Here's another, written in bash script. Just because you can doesn't mean you should write handlers in bash though.

```bash
#!/usr/bin/env bash

name=$(echo $1|jq '.["Params"]["name"][0]'|tr -d \")
cat <<- _EOF_
{
    "status": 200, 
    "header": {}, 
    "body": "hello $name"
}
_EOF_
```

For this I use the [`jq` tool](https://stedolan.github.io/jq/), a lightweight command line JSON processor. It parses the first command line argument `$1` and extracts the `name` from the `Params` field, which I then use to output back to STDOUT.



### Listeners

Listeners are a bit more complicated to write but they work the same way as bins. As long as your programming language of choice is able to parse JSON, and

When a HTTP request is sent to Tanuki (from a browser or another application), Tanuki creates a TCP connection to the listener and sends in the JSON HTTP request. Important to remember that the JSON request will end with a newline `\n` so when writing the listener you should use this to detect the end of the JSON string.

In the same way when you write your listener and construct a JSON HTTP response, you must end the JSON with a newline `\n`.

The listener must return a JSON HTTP response as above through the same connection.

This is an example of a simple listener written in Python. The file name is `listeners/hello-python` and the listener will be triggered when a GET request is sent to the URL `/_/hello/python/listener?name=sausheong`

```python
#!/usr/bin/env python

import socket
import sys
import json

HOST = '0.0.0.0'
PORT = sys.argv[1]

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((HOST, int(PORT)))
    s.listen()
    while True:
        conn, addr = s.accept()
        with conn:
            data = conn.recv(1024)
            request = json.loads(data)
            response = {
                'status': 200,
                'header': {},
                'body': "hello " + request['Params']['name'][0]
            }
            conn.sendall(str.encode(json.dumps(response)+"\n", 'utf-8' ))
```



## Why Tanuki?

It seems a bit of work to write a single web application in different programming languages, so why do it? It's really all about future-proofing.

### End programming wars

In any team with multiple capabilities, it's inevitable there's always going to be some programming language or library or tool war on which  technology to use. With Tanuki, you can opt for _all of the above_ since you can (if you want) write every handler in a different programming language or use different libraries or even different versions.

Why you would want to do that is to really tap on the different strengths of the various languages and libraries as needed.

### Get out of dependency hell

Most professional developers have one time or another in their careers needed to take over _legacy_ projects and code. Sometimes it could even be your own legacy code! For example if you have written a large web application using version 2.1 of a web framework a few years ago and now you need to upgrade to version 5.2 this could be pretty painful. Every library you used then would have been new versions and some of them would be incompatible with each other. Very often it's a multi-month or even year project to upgrade the whole application.

With Tanuki you can easily upgrade different parts of the same application over time. And because each handler is a separate piece of software on its own, you can opt to keep using the different versions of the libraries in them!

### Modular replacement

If the project was written in a programming language or uses a library or tool you don't know, you probably have to spent lots of time learning the new technologies, or try to re-write the whole thing using your own technology stack.

However, if the web application was written in Tanuki, you can opt to switch out parts of the project to change, and add new features in your own stack and keep the ones that don't need to change intact!

### Easy testing

Because each handler has standard input and output in a fixed format, testing the handlers can be pretty independent. In fact you can do test-driven development by writing the tests first and start writing the handlers until they all pass. You can also write the tests in one language and test the handlers written in other languages!

### Independent development

If you are building a large application or service and a number of developers are working on it at the same time is that you can farm out handlers to different people to write in different technology stacks, as long as they pass the tests! You can even deploy them on different machines as long as they are reachable by Tanuki.

### Distributed handling

Tanuki supports distributed handling and you can deploy your remote listener handlers to different servers, as long as the remote listener is reachable by Tanuki. This will help in scaling out, especially for resource intensive handlers.

## What's the trade-off?

You can't win 'em all. There's always a trade-off somewhere. Whether the trade-off is worth it depends very much on the web application or service itself. Here are some draw-backs of using Tanuki (besides the point that this is still an :warning: _experimental_ :warning: software!).

### Worse performance

Performance is definitely affected. In the case of bins, every time the action is called, the script is called to execute and return. Listeners are better because they are already running TCP servers, but the effort to create a socket connection and send the request JSON across has overhead. Things can get progressively worst if the listener is remote.

Nonetheless, if the processing done is significant, the percentage of the overhead could be small in comparison.

### No sharing of code or in-memory data

Since every handler (bin or listener) is a separate process, neither code nor data can be shared between them. To share data, you will need an intermediate store, for e.g. a database or a key-value store etc. To share code, one way of getting around it to is to use libraries. 

### Added complexities

Because the different handlers need to totally independent, there are additional overheads and complexities. For example, you need to parse the JSON request yourself and craft a JSON response almost in the raw. Besides adding to processing needs (and therefore lower performance) it also adds complexities in the code. And more complex code is overall harder to maintain.

On the other hand, some of the complexities can be grouped into reusable libraries, which helps to reduce the overheads.

### Not the same with other web frameworks

If you're familiar with other web frameworks, many of them work the same way. Tanuki is quite different from other frameworks because of it the way it works. You might even argue that calling it a web framework is a bit of a stretch but I would beg to differ, from [this definition from Wikipedia](https://en.wikipedia.org/wiki/Web_framework).

> A web framework (WF) or web application framework (WAF) is a software framework that is designed to support the development of web applications including web services, web resources, and web APIs. Web frameworks provide a standard way to build and deploy web applications on the World Wide Web. Web frameworks aim to automate the overhead associated with common activities performed in web development. 

Nonetheless because it works differently, first time developers on Tanuki would find a steeper learning curve. For example, if you're used to having a authentication and authorization drop-in library for other web frameworks, in Tanuki you might need to build that in yourself. To do this you might need to set and get the cookies to keep the state when you go from action to action.

## FAQs

### I'm still skeptical, why spend the effort?

Of course. Tanuki is not for every web application or service, and there is additional effort to do something that could be relatively simple to write. The real benefits of Tanuki comes with over the lifespan of the software. If you're planning to write something that needs to be maintained over a period of time by different developers or groups of developers, Tanuki can be a good choice (because of the reasons given above).

### Isn't this the same as CGI (Common Gateway Interface)?

Yes and no. Tanuki is partially inspired by [CGI](https://en.wikipedia.org/wiki/Common_Gateway_Interface) but there are some significant differences. In CGI, data to the CGI scripts are passed using environment variables, while in Tanuki it is sent into the bin through the command line argument as a JSON request. The Tanuki JSON request is also much more closer to the actual HTTP request.

Tanuki uses TCP socket servers as listeners to reduce the overhead of creating one process per requeset. This method is very similar to that of [FastCGI](https://en.wikipedia.org/wiki/FastCGI) but has a much simpler implementation. Also for FastCGI, a web application or service is built using a single programming language only. While FastCGI APIs exist for multiple languages, the intent was never to write one web application or service using more than one.

So for the implementations the mechanisms are close but not the same (Tanuki is simpler), neither CGI or FastCGI is meant for multiple programming languages like Tanuki is for.

### You were working on Polyglot earlier, what happened to it?

[Polyglot](https://github.com/sausheong/polyglot) was a previous incarnation of the same idea that I've been mulling over for many years and Tanuki is the latest (there were others in between but not released). There are several differences between Polyglot and Tanuki. Polyglot uses a fast message queue (ZeroMQ) as an intermediary to the handlers. This provides flexibility and scalability but it also adds quite a bit complexity. Because of the message queue, testing and debugging handlers can be difficult. Also, to create handlers the programming language needs to have the libraries to talk to the message queue. Also, securing the message queue can be quite a task because of the distributed nature of a Polyglot web application.

In designing Tanuki, I've deliberately gotten away from any message queues to reduce the implementation complexity.  In Tanuki, the handlers are simple programs that are either command line binaries or scripts, or TCP socket servers. This reduces the complexity in creating handlers and also reduces requirements to the most bare-bone -- as long as the language can parse a string, you can write a handler!

Nonetheless, there were good lessons learnt from my experience in Polyglot. However, I've stopped working on Polyglot already.

## Credits

Thanks to Jasmine Lim for the seriously cute mascot!

