package goaliniex

import (
	"io"
	"net/http"
)

type request struct {
	Method      string
	Endpoint    string
	Params      any
	SigningData []byte
	Header      http.Header
	Body        io.Reader
	FullURL     string
}
