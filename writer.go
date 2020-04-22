package nicehttp

import (
	"github.com/lithdew/bytesutil"
	"io"
)

var (
	_ io.Writer = (*WriterAtOffset)(nil)
	_ Writer    = (*WriteBuffer)(nil)
)

// Writer implements io.Writer and io.WriterAt.
type Writer interface {
	io.Writer
	io.WriterAt
}

// WriterAtOffset implements io.Writer for a given io.WriterAt at an offset.
type WriterAtOffset struct {
	dst    io.WriterAt
	offset int64
}

// NewWriterAtOffset instantiates a new writer at a specified offset.
func NewWriterAtOffset(dst io.WriterAt, offset int64) *WriterAtOffset {
	return &WriterAtOffset{dst: dst, offset: offset}
}

// Write implements io.Writer.
func (w WriterAtOffset) Write(b []byte) (int, error) {
	return w.dst.WriteAt(b, w.offset)
}

// WriteBuffer implements io.Writer and io.WriterAt on an optionally-provided byte slice.
type WriteBuffer struct {
	dst []byte
}

// NewWriteBuffer instantiates a new write buffer around dst. dst may be nil.
func NewWriteBuffer(dst []byte) *WriteBuffer {
	return &WriteBuffer{dst: dst}
}

// Write implements io.Writer.
func (b *WriteBuffer) Write(p []byte) (int, error) {
	b.dst = append(b.dst, p...)
	return len(p), nil
}

// WriteAt implements io.WriterAt.
func (b *WriteBuffer) WriteAt(p []byte, off int64) (int, error) {
	if min := int(off) + len(p); min > len(b.dst) {
		b.dst = bytesutil.ExtendSlice(b.dst, min)
	}
	return copy(b.dst[off:], p), nil
}

// Bytes returns the underlying byte slice.
func (b *WriteBuffer) Bytes() []byte {
	return b.dst
}
