package tests

var reqJSON string
var expected string

func init() {
	reqJSON = `{"Method":"GET","URL":{"Scheme":"","Opaque":"","Host":"","Path":"/_/hello/world","RawQuery":"name=sausheong","Fragment":""},"Proto":"HTTP/1.1","Header":{"Accept":["text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"],"Accept-Encoding":["gzip, deflate"],"Accept-Language":["en-sg"],"Connection":["keep-alive"],"Cookie":["hello=world; mykey=myvalue"],"Upgrade-Insecure-Requests":["1"],"User-Agent":["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.2 Safari/605.1.15"]},"Body":"","ContentLength":0,"TransferEncoding":null,"Host":"localhost:8080","Params":{"name":["sausheong"]},"Multipart":{},"RemoteAddr":"[::1]:52340","RequestURI":"/_/hello/world?name=sausheong"}`
	expected = `{"status":200,"header":{},"body":"hello sausheong"}` + "\n"
}
