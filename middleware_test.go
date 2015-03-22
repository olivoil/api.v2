package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Middleware", func() {
	Context("MiddlewareFunc", func() {
		It("find a Name", func() {
			middlewareOK := func(r *Req) (error, int) {
				return nil, 200
			}

			f := MiddlewareFunc(middlewareOK)

			Expect(f.Name()).To(Equal("github.com/olivoil/api.func·021"))
		})
	})
})
