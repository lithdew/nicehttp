// Package nicehttp contains helper utilities for downloading files/making requests with valyala/fasthttp.

package nicehttp

import (
	"github.com/valyala/fasthttp"
	"io"
	"runtime"
)

// defaultClient is a nicehttp.Client with sane configuration defaults.
var defaultClient = Client{
	// Allow for parallel chunk-based downloading.
	AcceptsRanges: true,

	// Default to the number of available CPUs.
	NumWorkers: runtime.NumCPU(),

	// 10 MiB chunks.
	ChunkSize: 10 * 1024 * 1024,

	// Redirect 16 times at most.
	MaxRedirectCount: 16,
}

// Do sends a HTTP request prescribed in req and populates its results into res. It additionally handles redirects
// unlike the de-facto Do(req, res) method in fasthttp.
func Do(req *fasthttp.Request, res *fasthttp.Response) error {
	return defaultClient.Do(req, res)
}

// QueryHeaders queries headers from url via a HTTP HEAD request, and populates dst with its contents.
func QueryHeaders(dst *fasthttp.ResponseHeader, url string) error {
	return defaultClient.QueryHeaders(dst, url)
}

// DownloadBytes downloads the contents of url, and returns them as a byte slice.
func DownloadBytes(dst []byte, url string) ([]byte, error) {
	return defaultClient.DownloadBytes(dst, url)
}

// DownloadFile downloads of url, and writes its contents to a newly-created file titled filename.
func DownloadFile(filename, url string) error {
	return defaultClient.DownloadFile(filename, url)
}

// DownloadSerially contents of url and writes it to w.
func DownloadSerially(w io.Writer, url string) error {
	return defaultClient.DownloadSerially(w, url)
}

// DownloadInChunks downloads file at url comprised of length bytes in chunks using multiple workers, and stores it in
// writer w.
func DownloadInChunks(w io.WriterAt, url string, length int) error {
	return defaultClient.DownloadInChunks(w, url, length)
}
