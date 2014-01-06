package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var typeTests = []struct {
	value    string
	expected string
}{
	{"", "text/html"},
	{"Text/HTML", "text/html"},
	{"text/*, text/html, text/html;level=1, */*", "text/*, text/html, text/html"},
	{"text/html;level=2", "text/html"},
	{"text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c", "text/plain"},
	{" application/json ", "application/json"},
}

func TestResolveContentType(t *testing.T) {
	for _, tt := range typeTests {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Content-Type", tt.value)
		req := Req{Request: r}

		if req.ResolveContentType() != tt.expected {
			t.Fatalf("%s must be %s", req.ResolveContentType(), tt.expected)
		}
	}
}

var jsonTests = []struct {
	body     string
	expected string
}{
	{"{\"foo\":42,\"bar\":\"baz\"}", "{\"foo\":42,\"bar\":\"baz\"}"},
	{"foo:42", "foo:42"},
}

func TestJsonBody(t *testing.T) {
	for _, tt := range jsonTests {
		body := bytes.NewBufferString(tt.body)
		r, _ := http.NewRequest("GET", "/", body)
		r.Header.Set("Content-Type", "application/json")
		req := Req{Request: r}

		json, err := req.JsonBody()
		if err != nil {
			t.Fatalf("%s failed to parse: %s", tt.body, err)
		}
		if string(json) != tt.expected {
			t.Fatalf("%s must be %s", string(json), tt.expected)
		}
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", "application/json")
	req := Req{Request: r}

	json, err := req.JsonBody()
	if err != nil {
		t.Fatalf("Parsing an empty body must not return an error")
	}
	if string(json) != "" {
		t.Fatalf("Parsing an empty body must have returned an empty string")
	}
}

func TestHandlePanic(t *testing.T) {
	var test = func(rec *httptest.ResponseRecorder) *Req {
		r, _ := http.NewRequest("GET", "/", nil)
		req := WrapReq(rec, r)
		defer req.handlePanic()
		panic("test")
		return req
	}

	recorder := httptest.NewRecorder()
	r := test(recorder)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Expected the panic to be handled with a 500 got a %d, %#v", recorder.Code, r)
	}
}

func TestResponseStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	req := WrapReq(rec, r)
	if req.ResponseStatus() != 0 {
		t.Fatalf("Expected the responseStatus to not be set, was %d", req.ResponseStatus())
	}
	codes := []int{http.StatusOK, http.StatusCreated, http.StatusMultipleChoices, http.StatusBadRequest, http.StatusInternalServerError}
	for _, code := range codes {
		rec = httptest.NewRecorder()
		r, _ = http.NewRequest("GET", "/", nil)
		req = WrapReq(rec, r)
		req.Response.WriteHeader(code)
		if req.ResponseStatus() != code {
			t.Fatalf("Expected the ResponseStatus to be %d but was %d", code, req.ResponseStatus())
		}
	}
}
