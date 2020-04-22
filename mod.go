// Package nicehttp contains helper utilities for downloading files/making requests with valyala/fasthttp.

package nicehttp

import (
	"github.com/valyala/fasthttp"
	"io"
	"time"
)

// defaultClient is a nicehttp.Client with sane configuration defaults.
var defaultClient = NewClient()

// Do sends a HTTP request prescribed in req and populates its results into res. It additionally handles redirects
// unlike the de-facto Do(req, res) method in fasthttp.
func Do(req *fasthttp.Request, res *fasthttp.Response) error {
	return defaultClient.Do(req, res)
}

// DoTimeout sends a HTTP request prescribed in req and populates its results into res. It additionally handles
// redirects unlike the de-facto Do(req, res) method in fasthttp. It overrides the default timeout set.
func DoTimeout(req *fasthttp.Request, res *fasthttp.Response, timeout time.Duration) error {
	return defaultClient.DoTimeout(req, res, timeout)
}

// DoDeadline sends a HTTP request prescribed in req and populates its results into res. It additionally handles
// redirects unlike the de-facto Do(req, res) method in fasthttp. It overrides the default timeout set with a deadline.
func DoDeadline(req *fasthttp.Request, res *fasthttp.Response, deadline time.Time) error {
	return defaultClient.DoDeadline(req, res, deadline)
}

// QueryHeaders learns from url its content length, and if it accepts parallel chunk fetching.
func QueryHeaders(url string) (contentLength int, acceptsRanges bool) {
	return defaultClient.QueryHeaders(url)
}

// QueryHeadersTimeout learns from url its content length, and if it accepts parallel chunk fetching.
func QueryHeadersTimeout(url string, timeout time.Duration) (contentLength int, acceptsRanges bool) {
	return defaultClient.QueryHeadersTimeout(url, timeout)
}

// QueryHeadersDeadline learns from url its content length, and if it accepts parallel chunk fetching.
func QueryHeadersDeadline(url string, deadline time.Time) (contentLength int, acceptsRanges bool) {
	return defaultClient.QueryHeadersDeadline(url, deadline)
}

// Download downloads the contents of url and writes its contents to w.
func Download(w Writer, url string, contentLength int, acceptsRanges bool) error {
	return defaultClient.Download(w, url, contentLength, acceptsRanges)
}

// DownloadTimeout downloads the contents of url and writes its contents to w.
func DownloadTimeout(w Writer, url string, contentLength int, acceptsRanges bool, timeout time.Duration) error {
	return defaultClient.DownloadTimeout(w, url, contentLength, acceptsRanges, timeout)
}

// DownloadDeadline downloads the contents of url and writes its contents to w.
func DownloadDeadline(w Writer, url string, contentLength int, acceptsRanges bool, deadline time.Time) error {
	return defaultClient.DownloadDeadline(w, url, contentLength, acceptsRanges, deadline)
}

// DownloadBytes downloads the contents of url, and returns them as a byte slice.
func DownloadBytes(dst []byte, url string) ([]byte, error) {
	return defaultClient.DownloadBytes(dst, url)
}

// DownloadBytesTimeout downloads the contents of url, and returns them as a byte slice.
func DownloadBytesTimeout(dst []byte, url string, timeout time.Duration) ([]byte, error) {
	return defaultClient.DownloadBytesTimeout(dst, url, timeout)
}

// DownloadBytesDeadline downloads the contents of url, and returns them as a byte slice.
func DownloadBytesDeadline(dst []byte, url string, deadline time.Time) ([]byte, error) {
	return defaultClient.DownloadBytesDeadline(dst, url, deadline)
}

// DownloadFile downloads of url, and writes its contents to a newly-created file titled filename.
func DownloadFile(filename, url string) error {
	return defaultClient.DownloadFile(filename, url)
}

// DownloadFileTimeout downloads of url, and writes its contents to a newly-created file titled filename.
func DownloadFileTimeout(filename, url string, timeout time.Duration) error {
	return defaultClient.DownloadFileTimeout(filename, url, timeout)
}

// DownloadFileDeadline downloads of url, and writes its contents to a newly-created file titled filename.
func DownloadFileDeadline(filename, url string, deadline time.Time) error {
	return defaultClient.DownloadFileDeadline(filename, url, deadline)
}

// DownloadSerially contents of url and writes it to w.
func DownloadSerially(w io.Writer, url string) error {
	return defaultClient.DownloadSerially(w, url)
}

// DownloadSeriallyTimeout contents of url and writes it to w.
func DownloadSeriallyTimeout(w io.Writer, url string, timeout time.Duration) error {
	return defaultClient.DownloadSeriallyTimeout(w, url, timeout)
}

// DownloadSeriallyDeadline contents of url and writes it to w.
func DownloadSeriallyDeadline(w io.Writer, url string, deadline time.Time) error {
	return defaultClient.DownloadSeriallyDeadline(w, url, deadline)
}

// DownloadInChunks downloads file at url comprised of length bytes in chunks using multiple workers, and stores it in
// writer w.
func DownloadInChunks(w io.WriterAt, url string, length int) error {
	return defaultClient.DownloadInChunks(w, url, length)
}

// DownloadInChunksTimeout downloads file at url comprised of length bytes in chunks using multiple workers, and stores
// it in writer w.
func DownloadInChunksTimeout(w io.WriterAt, url string, length int, timeout time.Duration) error {
	return defaultClient.DownloadInChunksTimeout(w, url, length, timeout)
}

// DownloadInChunksDeadline downloads file at url comprised of length bytes in chunks using multiple workers, and
// stores it in writer w.
func DownloadInChunksDeadline(w io.WriterAt, url string, length int, deadline time.Time) error {
	return defaultClient.DownloadInChunksDeadline(w, url, length, deadline)
}
