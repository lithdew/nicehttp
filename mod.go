// Package nicehttp contains helper utilities for downloading files/making requests with valyala/fasthttp.

package nicehttp

import (
	"github.com/valyala/fasthttp"
	"io"
	"runtime"
)

// DefaultClient is a nicehttp.Client with sane configuration defaults.
var DefaultClient = Client{
	Preallocate: true,

	AcceptsRanges: true,
	NumWorkers:    runtime.NumCPU(),
	RangeSize:     10 * 1024 * 1024,

	MaxRedirectCount: 16,
}

// Do sends a HTTP request prescribed in req and populates its results into res. It additionally handles redirects
// unlike the de-facto Do(req, res) method in fasthttp.
func Do(req *fasthttp.Request, res *fasthttp.Response) error {
	return DefaultClient.Do(req, res)
}

// QueryHeaders queries headers from url via a HTTP HEAD request, and populates dst with its contents.
func QueryHeaders(dst *fasthttp.ResponseHeader, url string) error {
	return DefaultClient.QueryHeaders(dst, url)
}

// DownloadFile downloads of url, and writes its contents to a newly-created file titled filename.
func DownloadFile(filename, url string) error {
	return DefaultClient.DownloadFile(filename, url)
}

// DownloadInChunks downloads file at url comprised of length bytes in cs byte-sized chunks using w goroutines, and
// store it in file f.
func DownloadInChunks(f io.WriterAt, url string, length, w, cs int) error {
	return DefaultClient.DownloadInChunks(f, url, length, w, cs)
}

// Download contents of url and write it to w.
func Download(w io.Writer, url string) error {
	return DefaultClient.Download(w, url)
}
