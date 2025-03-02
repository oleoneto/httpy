package cli

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"time"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	"github.com/oleoneto/mock-http/pkg/dbsql"
	"github.com/oleoneto/mock-http/pkg/schema"
	"github.com/spf13/cobra"
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

func DatabaseConnect(cmd *cobra.Command, args []string) {
	var err error
	database, err = dbsql.ConnectDatabase(dbsql.DBConnectOptions{
		VerboseLogging: state.Flags.VerboseLogging,
		Filename:       "mockhttp.sqlite3",
	})

	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	row := database.QueryRowContext(
		ctx,
		`SELECT s.name FROM sqlite_schema s WHERE s.type = 'table' AND s.name = 'responses' LIMIT 1`,
	)

	if reflect.DeepEqual(row.Scan(nil), sql.ErrNoRows) {
		if _, err = database.ExecContext(ctx, string(sqlSchema)); err != nil {
			panic(err)
		}
	}
}

type OutputFormat struct{ *toolkit.FlagEnum }

// Supported formats:
// - silent
// - json
// - yaml
func (f *OutputFormat) ProcessResponseOptions() schema.ProcessingOptions {
	switch f.FlagEnum.String() {
	case "yaml":
		return schema.ProcessingOptions{
			SQLPersistenceFunc: database.ExecContext,
			BodyMarshalFunc:    schema.BodyMarshalFunc,
			Plugins:            plugins,
			// FilePersistenceFunc:        os.WriteFile,
			// FilePersistenceMarshalFunc: yaml.Marshal,
			// FilePersistenceNamingFunc: func() string { return time.Now(time.RFC3339) + ".yaml" }
		}
	default:
		return schema.ProcessingOptions{
			SQLPersistenceFunc: database.ExecContext,
			BodyMarshalFunc:    schema.BodyMarshalFunc,
			Plugins:            plugins,
			// FilePersistenceFunc:        os.WriteFile,
			// FilePersistenceMarshalFunc: json.Marshal,
			// FilePersistenceNamingFunc: func() string { return time.Now(time.RFC3339) + ".json" }
		}
	}
}
