package api

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// This struct wraps a ResponseWriter to keep track of the status code, for logging purpose.
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (srw *statusResponseWriter) WriteHeader(status int) {
	srw.status = status
	srw.ResponseWriter.WriteHeader(status)
}

func (srw *statusResponseWriter) Write(b []byte) (int, error) {
	if srw.status == 0 {
		srw.status = http.StatusOK
	}
	return srw.ResponseWriter.Write(b)
}

// Implementation of the various interfaces that we may have hidden because of the wrapped ResponseWriter.
// See https://groups.google.com/d/topic/golang-nuts/zq_i3Hf7Nbs/discussion for details.
func (srw *statusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := srw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("ResponseWriter does not implement http.Hijacker")
}

func (srw *statusResponseWriter) Flush() {
	if f, ok := srw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (srw *statusResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := srw.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}
