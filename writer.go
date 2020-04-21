package nicehttp

import "io"

// WriterAtOffset implements io.Writer for a given io.WriterAt at an offset.
type WriterAtOffset struct {
	Src    io.WriterAt
	Offset int64
}

// Write implements io.Writer.
func (w WriterAtOffset) Write(b []byte) (int, error) {
	return w.Src.WriteAt(b, w.Offset)
}
