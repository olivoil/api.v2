package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Error implements the jsonapi.org spec for errors
type Error struct {
	ID     string `json:"id,omitempty"`
	Href   string `json:"href,omitempty"`
	Status string `json:"status,omitempty,required"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty,required"`
	Detail string `json:"detail,omitempty"`
	Path   string `json:"path,omitempty"`
}

// ErrorStack represents several errors
type ErrorStack struct {
	Errors []Error `json:"errors"`
}

// WrapErr automatically wraps a standard error to an api.Error
func WrapErr(err error, status int) Error {
	apiErr, ok := err.(Error)
	s := strconv.Itoa(status)

	if ok {
		apiErr.Status = s
		return apiErr
	}

	return Error{
		Status: s,
		Title:  err.Error(),
	}
}

// Add returns a stack of errors
func (e Error) Add(f Error) ErrorStack {
	return ErrorStack{
		Errors: []Error{e, f},
	}
}

// Error returns a nice string representation of an Error
func (e Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s error (%s) %s", e.Code, e.Status, e.Title)
	}

	return fmt.Sprintf("error (%s) %s", e.Status, e.Title)
}

// HTTPStatus returns the int value of the error status
func (e Error) HTTPStatus() int {
	v, err := strconv.Atoi(e.Status)
	if err != nil {
		return 0
	}
	return v
}

// HTTPBody returns the body of an http response for the error
func (e Error) HTTPBody() string {
	s := ErrorStack{Errors: []Error{e}}
	return s.HTTPBody()
}

// Add adds an error to a stack of error
func (s *ErrorStack) Add(e Error) {
	s.Errors = append(s.Errors, e)
}

// Error returns a nice string representation of an Error
func (e ErrorStack) Error() string {
	l := len(e.Errors)

	if l == 0 {
		return ""
	}

	err := e.Errors[l-1]

	if l == 1 {
		return err.Error()
	}

	return err.Error() + ", and " + strconv.Itoa(l) + " more errors"
}

// HTTPStatus returns the int value of the error status
func (e ErrorStack) HTTPStatus() int {
	l := len(e.Errors)
	err := e.Errors[l-1]
	return err.HTTPStatus()
}

// HTTPBody returns the body of an http response for the error
func (e ErrorStack) HTTPBody() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}
