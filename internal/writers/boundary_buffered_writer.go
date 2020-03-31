package writers

import (
	"bufio"
	"io"
)
// 包了个Buffer
type BoundaryBufferedWriter struct {
	bw *bufio.Writer
}

func NewBoundaryBufferedWriter(w io.Writer, size int) *BoundaryBufferedWriter {
	return &BoundaryBufferedWriter{
		bw: bufio.NewWriterSize(w, size),
	}
}

func (b *BoundaryBufferedWriter) Write(p []byte) (int, error) {
	if len(p) > b.bw.Available() {
		err := b.bw.Flush()
		if err != nil {
			return 0, err
		}
	}
	return b.bw.Write(p)
}

func (b *BoundaryBufferedWriter) Flush() error {
	return b.bw.Flush()
}
