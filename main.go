package main

import (
	"os"
	"reflect"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/oleoneto/httpy/cmd/cli"
	"github.com/oleoneto/httpy/pkg"
	"github.com/oleoneto/httpy/pkg/extensions"
	"github.com/sirupsen/logrus"
)

func main() {
	LoadExtensions()

	cli.Execute(pkg.CLIConfig{
		Plugins:        plugins,
		DefaultTimeout: helpers.PointerTo(1 * time.Minute),
	})
}

var data []byte
var plugins map[string]reflect.Value

var supportedExtensions = []string{
	"RequestTransformerFunc",
	"ResponseTransformerFunc",
	"ResponsePassesValidationFunc",
}

func LoadExtensions() {
	if filepath, ok := os.LookupEnv("PLUGINS_FILEPATH"); ok {
		var err error
		if data, err = os.ReadFile(filepath); err != nil {
			panic(err)
		}

		plugins = extensions.Load(string(data), supportedExtensions)
		for k, v := range plugins {
			logrus.Warnln("Loaded extension:", k, v)
		}
	}
}
