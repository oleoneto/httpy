package schema

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/oleoneto/go-toolkit/httpclient"
	"github.com/sirupsen/logrus"
)

// LoadSchema
func LoadSchema(input []byte, parserFunc func([]byte, any) error) (data Schema, err error) {
	logrus.Debugln(helpers.FuncName())
	err = parserFunc(input, &data)
	return
}

// Execute performs requests
func Execute(s Schema, client *http.Client, bodyDeserializerFunc func(any) ([]byte, error), receiver chan<- ResponseWrapper /* receive-only channel */) (total int) {
	logrus.Debugln(helpers.FuncName())

	for _, request := range s.Requests {
		// Nothing to do here...
		if request.Skip {
			continue
		}

		total++ // TODO: Review statement

		// Load payload, if one exists
		var err error
		var body []byte
		if request.Body != nil && bodyDeserializerFunc != nil {
			if request.Body, err = bodyDeserializerFunc(request.Body); err != nil {
				receiver <- ResponseWrapper{Request: request, Response: nil, Error: err}
				continue
			}

			fmt.Println(body)
		}

		// Perform request
		go func() {
			response, err := Send(client, request)
			receiver <- ResponseWrapper{Request: request, Response: response, Error: err}
		}()
	}

	return
}

// Send makes HTTP request using the provided client.
func Send(client *http.Client, request Request) (*http.Response, error) {
	logrus.Debugln(helpers.FuncName())

	if request.Timeout != nil {
		client.Timeout = *request.Timeout
	}

	r, err := client.Do(&http.Request{
		Method: request.Method,
		URL:    request.ParseURL(),
		Header: request.HeaderMap(),
		Body: func() io.ReadCloser {
			if request.Body != nil {
				return httpclient.NewBody(request.Body.([]byte))
			}
			return nil
		}(),
	})

	logrus.Debugln(helpers.FuncName(), "request sent...")
	return r, err
}

// Consume processes HTTP responses.
func Consume(count int, options ProcessingOptions, sender <-chan ResponseWrapper /* send-only channel */) {
	logrus.Debugln(helpers.FuncName())

	var persistLater []Response

	for range count {
		select {
		case responseWrapper := <-sender:
			if responseWrapper.Error != nil {
				logrus.Errorln(responseWrapper.Error)
				continue
			}

			logrus.Debugln("Response for:", responseWrapper.Request.Name)

			if p, ok := options.Plugins["ResponseTransformerFunc"]; ok {
				// NOTE: Proceed with caution
				if transformer, ok := p.Interface().(ResponseTransformerFunc); ok {
					logrus.Warnln("Running code for: ResponseTransformerFunc")
					responseWrapper.Response = transformer(responseWrapper.Response)
				}
			}

			response := Response{
				Error: responseWrapper.Error,
				Name:  responseWrapper.Request.Name,
				URL:   responseWrapper.Request.ParseURL().String(),
			}

			// MARK: Extract response body
			if responseWrapper.Response.Body != nil {
				// Fails if res.Body == nil OR cerr != nil
				response.Status = responseWrapper.Response.StatusCode

				// Fails if res.Body == nil OR cerr != nil
				response.Headers = responseWrapper.Response.Header

				// Extract response body
				// https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go
				body, err := io.ReadAll(responseWrapper.Response.Body)
				if err != nil {
					logrus.Warnln(err)
					continue
				}

				response.Body = body
				responseWrapper.Response.Body.Close()

				if options.BodyMarshalFunc != nil {
					b, err := options.BodyMarshalFunc(response.Body, response.Headers["Content-Type"][0])
					if err == nil {
						response.Body = b
					}
				}
			}

			if responseWrapper.ShouldPersist() {
				persistLater = append(persistLater, response)
			}

			// MARK: Run externally-defined tests

			if responseWrapper.Request.Tests == nil || responseWrapper.Request.Tests.Expectations.HTTP == nil {
				fmt.Println("No tests to run...")
				continue
			}

			if p, ok := options.Plugins["ResponsePassesValidationFunc"]; ok {
				// NOTE: Proceed with caution
				tester, ok := p.Interface().(func(*http.Request, *http.Response, struct {
					HTTP *struct {
						Status  *int
						Headers map[string][]string
						Body    any
					}
					Duration *time.Duration
				}) bool)
				if ok {
					logrus.Warnln("Running code for: ResponsePassesValidationFunc")
					logrus.Warnln("Running tests for:", responseWrapper.Request.Name)

					tester(responseWrapper.Response.Request, responseWrapper.Response, struct {
						HTTP *struct {
							Status  *int
							Headers map[string][]string
							Body    any
						}
						Duration *time.Duration
					}{
						HTTP: (*struct {
							Status  *int
							Headers map[string][]string
							Body    any
						})(responseWrapper.Request.Tests.Expectations.HTTP),
						// Duration: responseWrapper.Request.Tests.Expectations.Duration,
					})
				}
			}
		}
	}

	// TODO: Save responses in SQLite

	if len(persistLater) < 1 {
		return
	}

	rBytes, err := options.PersistenceMarshalFunc(Responses{persistLater})
	if err != nil {
		logrus.Errorln(err)
		return
	}

	// Persist response:
	if options.PersistenceFunc != nil {
		options.PersistenceFunc(rBytes)
	}
}

func BodyMarshalFunc(raw any, contentType string) (any, error) {
	logrus.Debugln(helpers.FuncName(), contentType)

	data, ok := raw.([]byte)
	if !ok {
		return raw, nil
	}

	switch contentType {
	case "application/json":
		return helpers.JSONPrettyPrint(string(data)), nil
	default:
		return raw, nil
	}
}

type Schema struct {
	Global struct {
		Timeout *time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	} `yaml:"global,omitempty" json:"global,omitempty"`
	Requests []Request `yaml:"requests" json:"requests"`
}

type RequestTransformerFunc func(*http.Request) *http.Request
type ResponseTransformerFunc func(*http.Response) *http.Response
type ResponsePassesValidationFunc func(*http.Request, *http.Response) bool
