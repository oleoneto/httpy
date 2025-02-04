package pkg

import (
	"reflect"
	"time"
)

type CLIConfig struct {
	Plugins        map[string]reflect.Value
	DefaultTimeout *time.Duration
}
