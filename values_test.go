package api

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Values", func() {
	It("has getters and setters", func() {
		v := Values{"foo": {"bar", "baz"}, "bar": {"baz"}}
		Expect(v.Get("foo")).To(Equal("bar"))
		Expect(v.GetAll("bar")).To(Equal([]string{"baz"}))

		v.Append("baz", "zyx")
		Expect(v.Get("baz")).To(Equal("zyx"))

		v.Set("baz", "xyz")
		Expect(v.Get("baz")).To(Equal("xyz"))
	})

	It("converts url.Values", func() {
		v := Values(url.Values{"foo": {"bar", "baz"}, "bar": {"baz"}})
		Expect(v.Get("foo")).To(Equal("bar"))
		Expect(v.GetAll("bar")).To(Equal([]string{"baz"}))

		v.Append("baz", "zyx")
		Expect(v.Get("baz")).To(Equal("zyx"))

		v.Set("baz", "xyz")
		Expect(v.Get("baz")).To(Equal("xyz"))
	})

	It("converts a map", func() {
		v := Values(map[string][]string{"foo": {"bar", "baz"}, "bar": {"baz"}})
		Expect(v.Get("foo")).To(Equal("bar"))
		Expect(v.GetAll("bar")).To(Equal([]string{"baz"}))

		v.Append("baz", "zyx")
		Expect(v.Get("baz")).To(Equal("zyx"))

		v.Set("baz", "xyz")
		Expect(v.Get("baz")).To(Equal("xyz"))
	})
})
