package nicehttp

import (
	"errors"
	"fmt"
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"time"
)

// Client wraps over fasthttp.Client a couple of useful helper functions.
type Client struct {
	// The underlying instance which nicehttp.Client wraps around.
	Instance fasthttp.Client

	// Decide whether or not URLs that accept being downloaded in parallel chunks are handled with multiple workers.
	AcceptsRanges bool

	// The number of workers that are to be spawned for downloading chunks in parallel.
	NumWorkers int

	// Size of individual byte chunks downloaded.
	ChunkSize int

	// Max number of redirects to follow before a request is marked to have failed.
	MaxRedirectCount int

	// Max timeout for a single download/fetch.
	Timeout time.Duration
}

// Do sends a HTTP request prescribed in req and populates its results into res. It additionally handles redirects
// unlike the de-facto Do(req, res) method in fasthttp.
func (c *Client) Do(req *fasthttp.Request, res *fasthttp.Response) error {
	return c.DoTimeout(req, res, c.Timeout)
}

// DoTimeout sends a HTTP request prescribed in req and populates its results into res. It additionally handles
// redirects unlike the de-facto Do(req, res) method in fasthttp. It overrides the default timeout set.
func (c *Client) DoTimeout(req *fasthttp.Request, res *fasthttp.Response, timeout time.Duration) error {
	return c.DoDeadline(req, res, time.Now().Add(timeout))
}

// DoDeadline sends a HTTP request prescribed in req and populates its results into res. It additionally handles
// redirects unlike the de-facto Do(req, res) method in fasthttp. It overrides the default timeout set with a deadline.
func (c *Client) DoDeadline(req *fasthttp.Request, res *fasthttp.Response, deadline time.Time) error {
	for i := 0; i <= c.MaxRedirectCount; i++ {
		if err := c.Instance.DoDeadline(req, res, deadline); err != nil {
			return err
		}

		if !fasthttp.StatusCodeIsRedirect(res.StatusCode()) {
			return nil
		}

		location := res.Header.Peek("Location")
		if len(location) == 0 {
			return errors.New("missing 'Location' header after redirect")
		}

		req.URI().UpdateBytes(location)

		res.Reset()
	}

	return errors.New("redirected too many times")
}

// QueryHeaders queries headers from url via a HTTP HEAD request, and populates dst with its contents.
func (c *Client) QueryHeaders(dst *fasthttp.ResponseHeader, url string) error {
	if dst == nil {
		return errors.New("dst must not be nil")
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod(fasthttp.MethodHead)
	req.SetRequestURI(url)

	if err := c.Do(req, res); err != nil {
		return fmt.Errorf("failed to call HEAD on url %q: %w", url, err)
	}

	res.Header.CopyTo(dst)

	return nil
}

// DownloadBytes downloads the contents of url, and returns them as a byte slice.
func (c *Client) DownloadBytes(dst []byte, url string) ([]byte, error) {
	contentLength, acceptsRanges := c.learn(url)

	w := NewWriteBuffer(bytesutil.ExtendSlice(dst, contentLength))

	if err := c.download(w, url, contentLength, acceptsRanges); err != nil {
		return w.dst, err
	}

	return w.dst, nil
}

// DownloadFile downloads the contents of url, and writes its contents to a newly-created file titled filename.
func (c *Client) DownloadFile(filename, url string) error {
	contentLength, acceptsRanges := c.learn(url)

	w, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open dest file: %w", err)
	}

	if err := w.Truncate(int64(contentLength)); err != nil {
		return fmt.Errorf("failed to truncate file to %d byte(s): %w", contentLength, err)
	}

	return c.download(w, url, contentLength, acceptsRanges)
}

// learn learns from url its content length, and if it accepts parallel chunk fetching.
func (c *Client) learn(url string) (contentLength int, acceptsRanges bool) {
	headers := acquireResponseHeaders()
	defer releaseResponseHeaders(headers)

	if err := c.QueryHeaders(headers, url); err == nil {
		if contentLength = headers.ContentLength(); contentLength <= 0 {
			contentLength = 0
		}

		acceptsRanges = bytesutil.String(headers.Peek("Accept-Ranges")) == "bytes"
	}

	return contentLength, acceptsRanges
}

// download downloads url and writes its contents to w.
func (c *Client) download(w Writer, url string, contentLength int, acceptsRanges bool) error {
	if c.AcceptsRanges && acceptsRanges {
		if contentLength <= 0 {
			return fmt.Errorf("content length is %d - see doc for (*fasthttp.ResponseHeader).ContentLength()", contentLength)
		}

		if err := c.DownloadInChunks(w, url, contentLength); err != nil {
			return err
		}

		return nil
	}

	if err := c.DownloadSerially(w, url); err != nil {
		return err
	}

	return nil
}

// DownloadSerially serially downloads the contents of url and writes it to w.
func (c *Client) DownloadSerially(w io.Writer, url string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)

	if err := c.Do(req, res); err != nil {
		return fmt.Errorf("failed to download %q: %w", url, err)
	}

	return res.BodyWriteTo(w)
}

// DownloadInChunks downloads file at url comprised of length bytes in chunks using multiple workers, and stores it in
// writer w.
func (c *Client) DownloadInChunks(f io.WriterAt, url string, length int) error {
	deadline := time.Now().Add(c.Timeout)

	timeout := fasthttp.AcquireTimer(c.Timeout)
	defer fasthttp.ReleaseTimer(timeout)

	var g errgroup.Group

	// ByteRange represents a byte range.
	type ByteRange struct{ Start, End int }

	ch := make(chan ByteRange, c.NumWorkers)

	// Spawn w workers that will dispatch and execute byte range-inclusive HTTP requests.

	for i := 0; i < c.NumWorkers; i++ {
		i := i

		g.Go(func() error {
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			res := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseResponse(res)

			req.SetRequestURI(url)

			for r := range ch {
				req.Header.SetByteRange(r.Start, r.End)

				if err := c.DoDeadline(req, res, deadline); err != nil {
					return fmt.Errorf("worker %d failed to get bytes range (start: %d, end: %d): %w", i, r.Start, r.End, err)
				}

				if err := res.BodyWriteTo(NewWriterAtOffset(f, int64(r.Start))); err != nil {
					return fmt.Errorf("worker %d failed to write to file at offset %d: %w", i, r.Start, err)
				}
			}

			return nil
		})
	}

	// Fill up ch with byte ranges to be download from url.

	var r ByteRange

Feed:
	for r.End < length {
		r.End += c.ChunkSize
		if r.End > length {
			r.End = length
		}

		select {
		case <-timeout.C:
			break Feed
		case ch <- r:
		}

		r.Start += c.ChunkSize
		if r.Start > length {
			r.Start = length
		}
	}

	close(ch)

	// Wait until all byte ranges have been downloaded, or return early if an error was encountered downloading
	// a chunk.

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to download %q in chunks: %w", url, err)
	}

	return nil
}
