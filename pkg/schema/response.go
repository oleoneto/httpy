package schema

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"reflect"
)

type Responses struct {
	Responses []Response `yaml:"responses" json:"responses"`
}

type Response struct {
	Error error `yaml:"-" json:"-"`

	// -----------------------------------------------------------
	// MARK: Serializable
	// -----------------------------------------------------------

	// A name to identify the request (not checked for uniqueness)
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// Full HTTP request URL
	//
	// Example: https://example.com/articles?sort=title
	URL string `yaml:"url,omitempty" json:"url,omitempty"`

	// HTTP response status code
	//
	// Example: 200
	Status int `yaml:"status,omitempty" json:"status,omitempty"`

	// TODO: Improve typing
	//
	// HTTP response headers
	Headers map[string][]string `yaml:"headers,omitempty" json:"headers,omitempty"`

	// HTTP response body
	Body any `yaml:"body,omitempty" json:"body,omitempty"`
}

type ResponseWrapper struct {
	Request  Request
	Response *http.Response
	Error    error
}

type ProcessingOptions struct {
	ShowStatus       bool
	ShowHeaders      bool
	ShowResponseBody bool

	Plugins map[string]reflect.Value

	BodyMarshalFunc func(body any, contentType string) (any, error)

	FilePersistenceMarshalFunc func(any) ([]byte, error)
	FilePersistenceFunc        func(name string, data []byte, perm os.FileMode) error
	FilePersistenceNamingFunc  func() string

	SQLPersistenceFunc func(context.Context, string, ...any) (sql.Result, error)
}

func (r *ResponseWrapper) ShouldPersist() bool { return !r.Request.SkipPersistence }
