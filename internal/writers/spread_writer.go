package writers

import (
	"io"
	"time"
)
// 分批写，每次write都是写到内存的队列，最后flush的时候 用送间隔长度/队列长度，每次这个间隔之后写出去一匹
type SpreadWriter struct {
	w        io.Writer
	interval time.Duration
	buf      [][]byte
	exitCh   chan int
}

func NewSpreadWriter(w io.Writer, interval time.Duration, exitCh chan int) *SpreadWriter {
	return &SpreadWriter{
		w:        w,
		interval: interval,
		buf:      make([][]byte, 0),
		exitCh:   exitCh,
	}
}

func (s *SpreadWriter) Write(p []byte) (int, error) {
	b := make([]byte, len(p))
	copy(b, p)
	s.buf = append(s.buf, b)
	return len(p), nil
}

func (s *SpreadWriter) Flush() {
	sleep := s.interval / time.Duration(len(s.buf))
	ticker := time.NewTicker(sleep)
	for _, b := range s.buf {
		s.w.Write(b)
		select {
		case <-ticker.C:
		case <-s.exitCh: // skip sleeps finish writes
		}
	}
	ticker.Stop()
	s.buf = s.buf[:0]
}
