package api

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	Context("Marshalling", func() {
		It("will be marshalled correctly with a wrapped error", func() {
			apiErr := WrapErr(errors.New("boom!"), 500)
			result := []byte(apiErr.HTTPBody())
			expected := []byte(`{"errors":[{"status":"500","title":"boom!"}]}`)
			Expect(result).To(MatchJSON(expected))
		})

		It("will be marshalled correctly with one error", func() {
			apiErr := &Error{Title: "Bad Request", Status: "400"}
			result := []byte(apiErr.HTTPBody())
			expected := []byte(`{"errors":[{"status":"400","title":"Bad Request"}]}`)
			Expect(result).To(MatchJSON(expected))
		})

		It("will be marshalled correctly with several errors", func() {
			errorOne := &Error{Title: "Bad Request", Status: "400"}

			errorTwo := &Error{
				ID:     "001",
				Href:   "http://bla/blub",
				Status: "500",
				Code:   "001",
				Title:  "Title must not be empty",
				Detail: "Never occures in real life",
				Path:   "#titleField",
			}

			apiErr := errorOne.Add(errorTwo)
			result := []byte(apiErr.HTTPBody())
			expected := []byte(`{"errors":[
				{"status":"400","title":"Bad Request"},
				{"id":"001","href":"http://bla/blub","status":"500","code":"001","title":"Title must not be empty","detail":"Never occures in real life","path":"#titleField"}
			]}`)
			Expect(result).To(MatchJSON(expected))
		})
	})
})
