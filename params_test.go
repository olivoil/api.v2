package api

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func TestParseParamsInvalidMultipart(t *testing.T) {
	r, err := http.NewRequest("POST", "/", bytes.NewReader([]byte("blabla")))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "multipart/form-data; boundary=Boundary+0xAbCdEfGbOuNdArY")
	req := &Req{Request: r,
		Response: nil,
		Params:   new(Params),
	}
	err = req.ParseParams()
	if !strings.HasPrefix(err.Error(), "multipart:") {
		t.Errorf("expected multipart error, got %s", err)
	}
}

func TestParseParamsValidMultipart(t *testing.T) {
	r, err := http.NewRequest("POST", "/", bytes.NewReader([]byte("blabla")))
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(buf)
	if err := mw.SetBoundary("Boundary+0xAbCdEfGbOuNdArY"); err != nil {
		t.Fatal(err)
	}
	w, err := mw.CreateFormFile("test", "toto.txt")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte("blabla"))
	mw.Close()
	r.Header.Set("Content-Type", "multipart/form-data; boundary=Boundary+0xAbCdEfGbOuNdArY")
	r.Body = ioutil.NopCloser(buf)
	req := &Req{Request: r,
		Response: nil,
		Params:   new(Params),
	}
	err = req.ParseParams()
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}
