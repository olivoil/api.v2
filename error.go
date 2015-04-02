package api

import (
	"encoding/json"
	"fmt"
	"runtime"
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
	Stack  []byte `json:"-"`
}

func NewError(status int, title string) Error {
	err := Error{Status: strconv.Itoa(status), Title: title}
	err.CaptureStackTrace()
	return err
}

// ErrorStack represents several errors
type Errors struct {
	Err []Error `json:"errors"`
}

// WrapErr automatically wraps a standard error to an api.Error
func WrapErr(err error, status int) Error {
	apiErr, ok := err.(Error)

	if ok {
		apiErr.CaptureStackTrace()
		return apiErr
	}

	apiErr = Error{
		Status: strconv.Itoa(status),
		Title:  err.Error(),
	}

	apiErr.CaptureStackTrace()
	return apiErr
}

func (e *Error) CaptureStackTrace() {
	stack := make([]byte, 1024*8)
	e.Stack = stack[:runtime.Stack(stack, true)]
}

// Add returns a stack of errors
func (e Error) Add(f Error) Errors {
	return Errors{
		Err: []Error{e, f},
	}
}

// Error returns a nice string representation of an Error
func (e Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] (%s) %s", e.Code, e.Status, e.Title)
	}

	return fmt.Sprintf("[http] (%s) %s", e.Status, e.Title)
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
	s := Errors{Err: []Error{e}}
	return s.HTTPBody()
}

// Add adds an error to a stack of error
func (s *Errors) Add(e Error) {
	s.Err = append(s.Err, e)
}

// Error returns a nice string representation of an Error
func (e Errors) Error() string {
	l := len(e.Err)

	if l == 0 {
		return ""
	}

	err := e.Err[l-1]

	if l == 1 {
		return err.Error()
	}

	if l == 2 {
		return err.Error() + ", and 1 more error"
	}

	return err.Error() + ", and " + strconv.Itoa(l-1) + " more errors"
}

// HTTPStatus returns the int value of the error status
func (e Errors) HTTPStatus() int {
	l := len(e.Err)
	err := e.Err[l-1]
	return err.HTTPStatus()
}

// HTTPBody returns the body of an http response for the error
func (e Errors) HTTPBody() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}
