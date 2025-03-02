package schema

import "time"

type Tests struct {
	Expectations Expectations `yaml:"expectations" json:"expectations"`
}

type Expectations struct {
	HTTP *struct {
		Status  *int                `yaml:"status"`
		Headers map[string][]string `yaml:"headers" json:"headers"`
		Body    any                 `yaml:"body" json:"body"`
	} `yaml:"http" json:"http"`

	Duration *time.Duration `yaml:"duration" json:"duration"`
}
