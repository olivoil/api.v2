package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Middleware", func() {
	Context("MiddlewareFunc", func() {
		It("find a Name", func() {
			middlewareOK := func(ctx context.Context, r *Req) (context.Context, error) {
				return ctx, nil
			}

			f := MiddlewareFunc(middlewareOK)

			Expect(f.Name()).To(ContainSubstring("github.com/olivoil"))
		})
	})
})
