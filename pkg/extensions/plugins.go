package extensions

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var extensionsPackageName = "httpy_extensions"

func Load(src string, names []string) map[string]reflect.Value {
	interpreter := interp.New(interp.Options{Unrestricted: false})
	interpreter.Use(stdlib.Symbols)

	_, err := interpreter.Eval(src)
	if err != nil {
		panic(err)
	}

	methods := make(map[string]reflect.Value)
	for _, name := range names {
		f, err := interpreter.Eval(extensionsPackageName + "." + name)
		if err != nil {
			continue
		}

		methods[name] = f
	}

	return methods
}

func Unwrap[T any](v reflect.Value) (f T, err error) {
	logrus.Debugln("Attempting to cast", v)
	if f, ok := v.Interface().(T); ok {
		return f, nil
	}

	n := reflect.New(reflect.TypeOf(f))

	// n.Call(in []reflect.Value)

	return n.Interface().(T), fmt.Errorf("unable to safely unwrap value")
}
