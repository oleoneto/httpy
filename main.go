package main

import (
	"os"
	"reflect"

	_ "github.com/joho/godotenv/autoload"
	"github.com/oleoneto/mock-http/cmd/cli"
	"github.com/oleoneto/mock-http/pkg"
	"github.com/oleoneto/mock-http/pkg/extensions"
	"github.com/sirupsen/logrus"
)

func main() {
	LoadExtensions()

	cli.Execute(pkg.CLIConfig{Plugins: plugins})
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
