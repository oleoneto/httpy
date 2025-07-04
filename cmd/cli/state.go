package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	"github.com/oleoneto/httpy/pkg/dbsql"
	"github.com/oleoneto/httpy/pkg/schema"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func DatabaseConnect(config, filename string) {
	filepath := fmt.Sprintf("%s/%s", *&config, filename)

	var err error
	database, err = dbsql.ConnectDatabase(dbsql.DBConnectOptions{
		VerboseLogging: state.Flags.VerboseLogging,
		Filename:       filepath,
	})

	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	row := database.QueryRowContext(
		ctx,
		`SELECT s.name FROM sqlite_schema s WHERE s.type = 'table' AND s.name = 'responses' LIMIT 1`,
	)

	if errors.Is(row.Scan(nil), sql.ErrNoRows) {
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
		}
	default:
		return schema.ProcessingOptions{
			SQLPersistenceFunc: database.ExecContext,
			BodyMarshalFunc:    schema.BodyMarshalFunc,
			Plugins:            plugins,
		}
	}
}

func newConfig(dburl string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	configPath := fmt.Sprintf("%s/.%s", home, RootCmd.Name()) // i.e /Users/alice/.httpy

	dir, derr := os.ReadDir(configPath)
	if derr == nil || len(dir) < 1 {
		// Directory does not exist. Create one.
		if merr := os.MkdirAll(configPath, 0755); merr != nil {
			log.Fatalf("Error: unable to create config at: %s. %s", configPath, merr)
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.SetDefault("database.engine", dbAdapter.String())
	viper.SetDefault("database.url", dburl)

	viper.SetDefault("server.port", serverPort)

	viper.SetDefault("stdout.format", state.Flags.OutputFormat)
	viper.SetDefault("stdout.verbose", state.Flags.VerboseLogging)

	err = viper.WriteConfig() // TODO: consider using `SafeWriteConfig()`
	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		return configPath, nil
	}

	return configPath, err
}
