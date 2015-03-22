package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var jsonTests = []struct {
	body     string
	expected string
}{
	{"{\"foo\":42,\"bar\":\"baz\"}", "{\"foo\":42,\"bar\":\"baz\"}"},
	{"foo:42", "foo:42"},
}

var _ = Describe("*Req", func() {
	It("resolves the content-type", func() {
		for _, tt := range typeTests {
			r, _ := http.NewRequest("GET", "/", nil)
			r.Header.Set("Content-Type", tt.value)
			req := Req{Request: r}

			Expect(req.ResolveContentType()).To(Equal(tt.expected))
		}
	})

	It("decodes json body", func() {
		for _, tt := range jsonTests {
			body := bytes.NewBufferString(tt.body)
			r, _ := http.NewRequest("GET", "/", body)
			r.Header.Set("Content-Type", "application/json")
			req := Req{Request: r}

			json, err := req.JsonBody()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(json)).To(Equal(tt.expected))
		}

		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Content-Type", "application/json")
		req := Req{Request: r}

		json, err := req.JsonBody()
		if err != nil {
			Fail("Parsing an empty body must not return an error")
		}
		if string(json) != "" {
			Fail("Parsing an empty body must have returned an empty string")
		}
	})

	It("handles panics", func() {
		var test = func(rec *httptest.ResponseRecorder) *Req {
			r, _ := http.NewRequest("GET", "/", nil)
			req := WrapReq(rec, r)
			defer req.handlePanic()
			panic("test")
			return req
		}

		recorder := httptest.NewRecorder()
		test(recorder)
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
	})

	It("reponds with the right status", func() {
		rec := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		req := WrapReq(rec, r)
		if req.ResponseStatus() != 0 {
			Fail("Expected the responseStatus to not be set, was %d", req.ResponseStatus())
		}
		codes := []int{http.StatusOK, http.StatusCreated, http.StatusMultipleChoices, http.StatusBadRequest, http.StatusInternalServerError}
		for _, code := range codes {
			rec = httptest.NewRecorder()
			r, _ = http.NewRequest("GET", "/", nil)
			req = WrapReq(rec, r)
			req.Response.WriteHeader(code)
			if req.ResponseStatus() != code {
				Fail("Expected the ResponseStatus to be %d but was %d", code, req.ResponseStatus())
			}
		}
	})
})
