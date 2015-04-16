package api

import (
	"strings"

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

			Expect(strings.Contains(f.Name(), "github.com/olivoil/api2.func")).To(Equal(true))
		})
	})
})
