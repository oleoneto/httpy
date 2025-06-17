package cli

import (
	"encoding/json"
	"log"

	"github.com/oleoneto/go-toolkit/httpclient"
	"github.com/oleoneto/httpy/pkg/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var RequestCmd = &cobra.Command{
	Use:    "http",
	Short:  "Make HTTP requests",
	PreRun: state.Flags.File.StdinHook("file"),
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Parse file

		file, err := schema.LoadSchema(state.Flags.File.Data, yaml.Unmarshal)
		if err != nil {
			log.Fatalln(err)
		}

		// 2. Make requests

		watcher := make(chan schema.ResponseWrapper)
		client := httpclient.New()

		client.Timeout = globalTimeout
		if file.Global.Timeout != nil {
			client.Timeout = *file.Global.Timeout
		}

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
	RequestCmd.Flags().VarP(&state.Flags.File, "file", "f", "")
}
