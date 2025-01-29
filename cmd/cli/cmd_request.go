package cli

import (
	"log"
	"os"
	"time"

	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/oleoneto/go-toolkit/httpclient"
	"github.com/oleoneto/mock-http/cmd/cli/requests"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var RequestCmd = &cobra.Command{
	Use:    "http",
	Short:  "Make HTTP requests",
	PreRun: state.Flags.File.StdinHook("file"),
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Parse file
		schema, err := requests.ParseSchema(debugFunc, state.Flags.File.Data, yaml.Unmarshal)
		if err != nil {
			log.Fatalln(err)
		}

		// 2. Make requests
		client := httpclient.New()
		client.Timeout = globalTimeout
		if schema.Global.Timeout != nil {
			client.Timeout = *schema.Global.Timeout
		}

		execChan := make(chan requests.Response)

		count := requests.MakeRequests(debugFunc, schema, client, execChan)

		// 3. Read Responses
		requests.ProcessResponses(debugFunc, count, execChan, requests.ProcessResponseOptions{
			BodyMarshalFunc:        bodyMarshalFunc,
			PersistenceFunc:        os.WriteFile,
			PersistenceMarshalFunc: yaml.Marshal,
			PersistenceNameFunc: func() string {
				return time.Now().Format(time.RFC3339) + ".yaml"
			},
		})
	},
}

func init() {
	RequestCmd.Flags().VarP(&state.Flags.File, "file", "f", "")
}

func bodyMarshalFunc(raw any, contentType string) (any, error) {
	loggerFunc(helpers.FuncName(), contentType)

	data, ok := raw.([]byte)
	if !ok {
		return raw, nil
	}

	switch contentType {
	case "application/json":
		return helpers.JSONPrettyPrint(string(data)), nil
	default:
		return string(data), nil
	}
}

var globalTimeout = 1 * time.Minute // TODO: Constantize default timeout
