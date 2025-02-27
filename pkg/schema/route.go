package schema

type Route struct {
	// A name to identify the route (not checked for uniqueness)
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// HTTP request path
	//
	// Default: none
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// HTTP route method
	//
	// Default: GET
	Method string `yaml:"method,omitempty" json:"method,omitempty"`

	// HTTP response body
	//
	// Default: none
	Body any `yaml:"body,omitempty" json:"body,omitempty"`

	// HTTP response status code
	//
	// Default: none
	StatusCode int `yaml:"status,omitempty" json:"status,omitempty"`
}
