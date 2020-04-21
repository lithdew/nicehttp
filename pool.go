package nicehttp

import (
	"github.com/valyala/fasthttp"
	"sync"
)

var responseHeaderPool sync.Pool

func acquireResponseHeaders() *fasthttp.ResponseHeader {
	h := responseHeaderPool.Get()
	if h == nil {
		h = &fasthttp.ResponseHeader{}
	}
	return h.(*fasthttp.ResponseHeader)
}

func releaseResponseHeaders(h *fasthttp.ResponseHeader) {
	responseHeaderPool.Put(h)
}
