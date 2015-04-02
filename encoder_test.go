package api

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Sample struct {
	Visible string `json:"visible"`
	Hidden  string `json:"hidden" out:"false"`
}

var _ = Describe("Encoder", func() {
	It("works", func() {
		src := &Sample{Visible: "visible", Hidden: "this field won't be exported"}
		dst := &Sample{}

		enc := &JsonEncoder{}
		result, err := enc.Encode(src)
		Expect(err).ToNot(HaveOccurred())

		err = json.Unmarshal(result, dst)
		Expect(err).ToNot(HaveOccurred())

		Expect(dst.Hidden).To(Equal(""))
	})
})
