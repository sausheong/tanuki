# Tanuki

Tanuki is a polyglot web framework that allows developers to write web applications and services in multiple programming languages. The basic concept is simple -- an acceptor (which is a HTTP server) accepts HTTP requests and redirects it to either an executable binary file (or script) or a TCP socket server to handle. These handlers can be written in different programming languages and run independently from each other.

Tanuki is experimental software at this point in time. Use at your own risk.

## Installing Tanuki

Tanuki is written in Go. To install it in your platform, you can install [Go](https://golang.org) and get the source:

`go get github.com/sausheong/tanuki`

There will be downloadable binaries for your platform at a later date.

## How to use Tanuki

Once you can have downloaded the source, you can build the command line tool `tanuki`.

`go build`

With tha command line tool, can you create a skeleton structure for your new web application or service:

`./tanuki create poko`

This will create a new directory named `poko` with the necessary struct and files. To start your new Tanuki application, go to your new application directory and run this:

`./tanuki start`

This will start the Tanuki application at port `8080`. You can port number or the IP address or hostname as well, just follow the instructions in the `tanuki` command line tool.

If you have executable binaries or scripts in the correct format in your application `/bin` directory, they will be loaded and you can call them immediately. If you have listeners in the correct format in your application `/listeners` directory they will be started as well by Tanuki.

The file name format for bins and listeners are:

`<METHOD>__<PATH>`

The forward slashes in the path `/` will be converted to double underscores `__`. Let's say you have the a Ruby executable script `bin/get__hello__ruby` (notice the double underscores), this means you can open up the browser and go to the URL `/_/hello/ruby` and this will execute the script. Remember, all Tanuki actions start with `/_/`.

## Tanuki actions


## Why Tanuki?

It seems a bit of work to write a single web application in different programming languages, so why do it? It's really all about future-proofing.

### No more fighting

### Easy upgrade

### Switching technologies

### Modular replacement

### Easy testing


## What's the trade-off?