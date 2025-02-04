package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	"github.com/oleoneto/mock-http/pkg/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func BeforeHook(state toolkit.CommandState) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		state.SetFormatter(cmd, args)

		if !state.Flags.TimeExecutions {
			return
		}

		state.ExecutionStartTime = time.Now()
	}
}

func AfterHook(state toolkit.CommandState) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if !state.Flags.TimeExecutions {
			return
		}

		fmt.Fprintln(
			os.Stderr,
			append([]any{"Elapsed time:", time.Since(state.ExecutionStartTime)}, state.ExecutionExitLog...)...,
		)
	}
}

/*
func (c *State) ConnectDatabase(cmd *cobra.Command, args []string) {
	dbpath := viper.Get("database.path")

	db, err := dbsql.ConnectDatabase(dbsql.DBConnectOptions{
		Adapter:        dbsql.SQLAdapter(cmd.Flag("adapter").Value.String()),
		DSN:            *c.Flags.DatabaseURL,
		Filename:       dbpath.(string),
		VerboseLogging: c.Flags.VerboseLogging,
	})
	if err != nil {
		panic(err)
	}

	c.Database = db
}
*/

type OutputFormat struct{ *toolkit.FlagEnum }

// Supported formats:
// - silent
// - json
// - yaml
func (f *OutputFormat) ProcessResponseOptions() schema.ProcessingOptions {
	switch f.FlagEnum.String() {
	case "yaml":
		return schema.ProcessingOptions{
			// PersistenceFunc:        os.WriteFile,
			PersistenceMarshalFunc: yaml.Marshal,
			BodyMarshalFunc:        schema.BodyMarshalFunc,
			Plugins:                plugins,
		}
	default:
		return schema.ProcessingOptions{
			// PersistenceFunc:        os.WriteFile,
			PersistenceMarshalFunc: json.Marshal,
			BodyMarshalFunc:        schema.BodyMarshalFunc,
			Plugins:                plugins,
		}
	}
}
