package api

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Params", func() {
	It("parses invalid multipart requests", func() {
		r, err := http.NewRequest("POST", "/", bytes.NewReader([]byte("blabla")))
		Expect(err).ToNot(HaveOccurred())
		r.Header.Set("Content-Type", "multipart/form-data; boundary=Boundary+0xAbCdEfGbOuNdArY")
		req := &Req{
			Request:  r,
			Response: nil,
			Params:   new(Params),
		}
		err = req.ParseParams()
		Expect(strings.HasPrefix(err.Error(), "multipart:")).To(Equal(true))
	})

	It("parses params from valid multipart request", func() {
		r, err := http.NewRequest("POST", "/", bytes.NewReader([]byte("blabla")))
		Expect(err).ToNot(HaveOccurred())
		buf := bytes.NewBuffer(nil)
		mw := multipart.NewWriter(buf)
		err = mw.SetBoundary("Boundary+0xAbCdEfGbOuNdArY")
		Expect(err).ToNot(HaveOccurred())
		w, err := mw.CreateFormFile("test", "toto.txt")
		Expect(err).ToNot(HaveOccurred())
		w.Write([]byte("blabla"))
		mw.Close()
		r.Header.Set("Content-Type", "multipart/form-data; boundary=Boundary+0xAbCdEfGbOuNdArY")
		r.Body = ioutil.NopCloser(buf)
		req := &Req{
			Request:  r,
			Response: nil,
			Params:   new(Params),
		}
		err = req.ParseParams()
		Expect(err).ToNot(HaveOccurred())
	})

	It("splits values separated by commas", func() {
		r, err := http.NewRequest("POST", "/users/:id?bonjour=hello,hola", nil)
		Expect(err).ToNot(HaveOccurred())
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		req := &Req{
			Request:  r,
			Response: nil,
			Params:   new(Params),
		}
		err = req.ParseParams()
		Expect(err).ToNot(HaveOccurred())
		Expect(req.Params.Get("bonjour")).To(Equal("hello"))
		Expect(len(req.Params.GetAll("bonjour"))).To(Equal(2))
		Expect(req.Params.GetAll("bonjour")).To(Equal([]string{"hello", "hola"}))
	})
})
