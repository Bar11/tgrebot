// Package log
//
// @author: xwc1125
// @date: 2021/3/29
package log

import (
	"io"
)

var Writer writer

type writer struct {
	io.WriteCloser
	Handler
}

func (w writer) WriterHandler(writer io.WriteCloser, fmtr Format) Handler {
	return WriterHandler(writer, fmtr)
}

func (w *writer) Close() error {
	return w.WriteCloser.Close()
}

func RotatingDayFileHandler(w *TimeWriter, fmtr Format) Handler {
	return writer{w, StreamHandler(w, fmtr)}
}

func WriterHandler(w io.WriteCloser, fmtr Format) Handler {
	return writer{w, StreamHandler(w, fmtr)}
}
