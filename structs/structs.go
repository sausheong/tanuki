package structs

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

// ResponseInfo corresponds to a HTTP 1.1 response
type ResponseInfo struct {
	Status int                 `json:"status"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}
