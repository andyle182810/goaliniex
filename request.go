package goaliniex

import (
	"io"
	"net/http"
)

type request struct {
	Method    string
	Endpoint  string
	Params    any
	UserEmail string
	Header    http.Header
	Body      io.Reader
	FullURL   string
	Public    bool
}
