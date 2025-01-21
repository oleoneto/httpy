package main

import (
	"embed"

	_ "github.com/joho/godotenv/autoload"
	"github.com/oleoneto/mock-http/cmd/cli"
)

func main() {
	cli.Execute(data, buildHash)
}

var (
	//go:embed data
	data embed.FS

	buildHash string = ""
)
