package requests

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/oleoneto/go-toolkit/helpers"
)

func ProcessResponses(debugFunc func(...any), count int, execChan <-chan Response, options ProcessResponseOptions) {
	var responses []Response

	for range count {
		select {
		case res := <-execChan:
			debugFunc(helpers.FuncName(), res.URL)

			if res.Error != nil {
				log.Println(res.Error)
				continue
			}

			if res.skipPersistence {
				continue
			}

			if options.BodyMarshalFunc != nil {
				res.Body, _ = options.BodyMarshalFunc(res.Body, strings.Split(res.response.Header.Get("Content-Type"), ";")[0])
			}

			responses = append(responses, res)
		}
	}

	if len(responses) < 1 {
		return
	}

	rBytes, _ := options.PersistenceMarshalFunc(Responses{responses})

	// Persist to file:
	filename := options.PersistenceNameFunc()
	debugFunc(filename)
	debugFunc(string(rBytes))

	options.PersistenceFunc(filename, rBytes, 0755)
}

type ProcessResponseOptions struct {
	ShowStatus, ShowHeaders, ShowResponseBody bool

	BodyMarshalFunc        func(body any, contentType string) (any, error)
	PersistenceMarshalFunc func(any) ([]byte, error)
	PersistenceNameFunc    func() string
	PersistenceFunc        func(name string, data []byte, mode os.FileMode) error
}

type (
	Responses struct {
		Responses []Response `yaml:"responses" json:"responses"`
	}

	Response struct {
		// Private identifier and sort key
		Id string `yaml:"id,omitempty" json:"id,omitempty"`

		url             *url.URL       `yaml:"-" json:"-"`
		response        *http.Response `yaml:"-" json:"-"`
		headers         map[string][]string
		skipPersistence bool `yaml:"-" json:"-"`

		// Public API
		Name    string `yaml:"name,omitempty" json:"name,omitempty"`
		URL     string `yaml:"url,omitempty" json:"url,omitempty"`
		Headers any    `yaml:"headers,omitempty" json:"headers,omitempty"` // TODO: Improve typing
		Status  string `yaml:"status,omitempty" json:"status,omitempty"`
		Body    any    `yaml:"body,omitempty" json:"body,omitempty"`
		Error   error  `yaml:"-,omitempty" json:"-,omitempty"` // TODO: Consider omitting field
	}
)
