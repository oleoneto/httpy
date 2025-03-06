package cli

import (
	"os"
	"reflect"
	"strconv"
	"time"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	"github.com/oleoneto/mock-http/pkg"
	"github.com/oleoneto/mock-http/pkg/dbsql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:              "mockhttp",
	Short:            "Mock HTTP is a CLI tool for making HTTP requests",
	PersistentPreRun: DatabaseConnect,
	PostRun:          func(cmd *cobra.Command, args []string) {},
	Run:              func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func Execute(config pkg.CLIConfig) error {
	plugins = config.Plugins
	sqlSchema = config.SQLSchema

	if config.DefaultTimeout != nil {
		globalTimeout = *config.DefaultTimeout
	}

	return RootCmd.Execute()
}

func init() {
	logrus.SetLevel(func() logrus.Level {
		var base logrus.Level = logrus.InfoLevel

		s := os.Getenv("LOG_LEVEL")
		level, err := strconv.Atoi(s)
		if err == nil {
			return logrus.Level(level)
		}

		return base
	}())

	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cobra.OnInitialize(func() { /* config code */ })

	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(RequestCmd)
	RootCmd.AddCommand(MockServerCmd)

	// MARK: Set up global flags
	RootCmd.PersistentFlags().BoolVar(&state.Flags.VerboseLogging, "verbose", state.Flags.VerboseLogging, "enable detailed logging")
	RootCmd.PersistentFlags().VarP(state.Flags.OutputFormat, "output", "o", "output format")

	RootCmd.PersistentFlags().VarP(dbAdapter, "db-adapter", "a", "database adapter")
	RootCmd.PersistentFlags().BoolVar(&state.Flags.TimeExecutions, "time", state.Flags.TimeExecutions, "time executions")
}

var (
	dbAdapter = &toolkit.FlagEnum{
		Allowed: []string{"sqlite3"},
		Default: "sqlite3",
	}

	outputFormat = OutputFormat{
		&toolkit.FlagEnum{
			Allowed: []string{"json", "yaml", "silent"},
			Default: "yaml",
		},
	}

	state = toolkit.CommandState{
		Flags: toolkit.CommandFlags{
			OutputTemplate: "",
			OutputFormat:   outputFormat.FlagEnum,
		},
	}

	globalTimeout = 1 * time.Minute

	plugins = make(map[string]reflect.Value)

	sqlSchema []byte

	database dbsql.SqlBackend
)
