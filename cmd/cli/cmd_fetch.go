package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/oleoneto/go-toolkit/httpclient"
	"github.com/oleoneto/httpy/pkg/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var FetchCmd = &cobra.Command{
	Use:       "fetch",
	Short:     "Make HTTP requests",
	ValidArgs: []string{"url"},
	Args:      cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		fetchMode = func(cmd *cobra.Command) FetchMode {
			if cmd.Flag("file").Changed {
				return CollectionMode
			}
			return RawMode
		}(cmd)

		if fetchMode == RawMode {
			if len(args) != 1 {
				log.Fatalln("missing required argument URL")
				os.Exit(1)
			}

			httpRequest.URL = args[0]

			return
		}

		state.Flags.File.StdinHook("file")
	},
	Run: func(cmd *cobra.Command, args []string) {
		client := httpclient.New()
		client.Timeout = httpyFlags.timeout

		var file schema.Schema
		var err error

		switch fetchMode {
		case CollectionMode:
			// MARK: Parse contents from file or stdin

			file, err = schema.LoadSchema(state.Flags.File.Data, yaml.Unmarshal)
			if err != nil {
				log.Fatalln(err)
			}

			if file.Global.Timeout != nil {
				client.Timeout = *file.Global.Timeout
			}
		case RawMode:
			if httpRequest.Name == "" {
				httpRequest.Name = httpRequest.ParseURL().Host
			}
			file = schema.Schema{Requests: []schema.Request{httpRequest}}
		}

		// MARK: Make requests

		watcher := make(chan schema.ResponseWrapper)
		count := schema.Execute(file, client, json.Marshal /* TODO: check this */, watcher)

		// Send Requests + Read Responses
		schema.Process(
			count,
			outputFormat.ProcessResponseOptions(),
			watcher,
		)
	},
}

func init() {
	FetchCmd.Flags().VarP(&state.Flags.File, "file", "f", "FILE or stdin")

	FetchCmd.Flags().StringVar(&httpRequest.Method, "method", httpRequest.Method, "http method")
	FetchCmd.Flags().StringToStringVar(&httpRequest.Headers, "headers", nil, "http headers")
	FetchCmd.Flags().BytesBase64Var(&body, "body", body, "http request body")
	FetchCmd.Flags().StringVar(&httpRequest.Name, "name", httpRequest.Name, "a name for this request (useful when persisting responses)")
}

type FetchMode string

const (
	RawMode        FetchMode = "raw"
	CollectionMode FetchMode = "collection"
)

var (
	body        []byte
	httpRequest schema.Request = schema.Request{Method: "GET", Scheme: "http"}
	fetchMode   FetchMode
)
