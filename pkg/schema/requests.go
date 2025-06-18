package schema

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/sirupsen/logrus"
)

type Request struct {
	// A name to identify the request (not checked for uniqueness)
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// HTTP request timeout (when specified, it overrides the schema's `global.timeout`)
	//
	// Default: none
	//
	// Example: 5s
	Timeout *time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`

	// HTTP scheme (accepts: [http, https])
	//
	// Default: http
	Scheme string `yaml:"scheme,omitempty" json:"scheme,omitempty"`

	// HTTP request URL
	//
	// Default: none
	//
	// Example: example.com
	URL string `yaml:"url,omitempty" json:"url,omitempty"`

	// HTTP request path
	//
	// Default: none
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// HTTP request query
	//
	// Default: none
	//
	// Example:
	//  match=colour
	//  lang=en-GB
	Query map[string]string `yaml:"query,omitempty" json:"query,omitempty"`

	// HTTP request method
	//
	// Default: GET
	Method string `yaml:"method,omitempty" json:"method,omitempty"`

	// HTTP request headers
	//
	// Default: none
	//
	// Example:
	//  Accepts: application/json
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`

	// HTTP request body
	//
	// Default: none
	Body any `yaml:"body,omitempty" json:"body,omitempty"`

	// Request will be skipped from HTTP calls
	//
	// Default: false
	Skip bool `yaml:"skip,omitempty" json:"skip,omitempty"`

	// Request will not be persisted if session has active persistence
	//
	// Default: false
	SkipPersistence bool `yaml:"skipPersistence,omitempty" json:"skipPersistence"`

	// Assertions to be run on the response of this request
	Tests *Tests `yaml:"tests" json:"tests"`

	url *url.URL
}

func (r *Request) ParseURL() *url.URL /* TODO: Add error */ {
	if r.url != nil {
		return r.url
	}

	logrus.Debugln(helpers.FuncName())

	var query string = "?"
	for k, v := range r.Query {
		query += k + "=" + v + "&"
	}

	query = strings.TrimSuffix(query, "&")
	if len(query) == 1 {
		query = ""
	}

	// We may be in the presence of a raw url string
	rawScheme, scheme := strings.Split(r.URL, "://"), ""
	if len(rawScheme[0]) <= 5 {
		scheme = ""
	} else {
		scheme = "http://"
	}

	url, err := url.Parse(scheme + r.URL + r.Path + query)
	if err != nil {
		log.Fatal(err)
	}

	r.url = url

	logrus.Debugln(helpers.FuncName(), r.url.String())

	return r.url
}

func (r *Request) HeaderMap() http.Header {
	logrus.Debugln(helpers.FuncName())

	headers := http.Header{}
	for k, v := range r.Headers {
		headers.Set(k, v)
	}

	return headers
}

var reqOnce sync.Once
