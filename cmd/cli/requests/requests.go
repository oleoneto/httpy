package requests

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/oleoneto/go-toolkit/httpclient"
)

func MakeRequests(debugFunc func(...any), in Schema, client *http.Client, execChan chan Response) (count int) {
	debugFunc(helpers.FuncName())

	for _, req := range in.Requests {
		if req.Skip {
			debugFunc("Skipped request:", req.Name, req.URL)
			continue
		}

		count++ // excludes skipped requests

		// MARK: URL Parsing

		var query string = "?"
		for k, v := range req.Query {
			query += k + "=" + v + "&"
		}
		query = strings.TrimSuffix(query, "&")
		if len(query) == 1 {
			query = ""
		}

		u, err := url.Parse(req.Scheme + "://" + req.URL + req.Path + query)
		if err != nil {
			log.Fatal(err)
		}

		endpoint := u

		// MARK: Body

		body, jerr := json.Marshal(req.Body)
		if jerr != nil {
			execChan <- Response{
				Id:              req.id,
				url:             endpoint,
				response:        nil,
				skipPersistence: req.SkipPersistence,

				// Public API
				Name:  req.Name,
				URL:   endpoint.String(),
				Error: jerr,
			}

			continue
		}

		// MARK: Headers

		headers := http.Header{}
		for k, v := range req.Headers {
			headers.Set(k, v)
		}

		// MARK: Send Request
		go sendRequest(req, client, debugFunc, endpoint, headers, body, execChan)
	}

	return count
}

func sendRequest(req Request, client *http.Client, debugFunc func(...any), endpoint *url.URL, headers http.Header, body []byte, execChan chan Response) {
	if req.Timeout != nil {
		client.Timeout = *req.Timeout
	}

	debugFunc(helpers.FuncName(), "Timeout set to", client.Timeout)

	res, cerr := client.Do(&http.Request{
		Method: req.Method,
		URL:    endpoint,
		Header: headers,
		Body:   httpclient.NewBody(body),
	})

	response := Response{
		Id:              req.id,
		url:             endpoint,
		response:        res,
		skipPersistence: req.SkipPersistence,

		// Public API
		Name:  req.Name,
		URL:   endpoint.String(),
		Error: cerr,
	}

	// TODO: Implement error handling
	// if cerr != nil {}

	if cerr == nil && res.Body != nil {
		// Extract response body
		// https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go
		response.Body, response.Error = io.ReadAll(res.Body)

		defer res.Body.Close()
	}

	execChan <- response
}

type Request struct {
	// Private identifier and sort key
	id string

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
}
