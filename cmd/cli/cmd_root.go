package cli

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/oleoneto/httpy/pkg"
	"github.com/oleoneto/httpy/pkg/dbsql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var entrypoint = "httpy"

var RootCmd = &cobra.Command{
	Use:              entrypoint,
	Short:            "HTTPy, a CLI tool for programmatically managing collections of HTTP requests",
	PersistentPreRun: func(cmd *cobra.Command, args []string) { DatabaseConnect(configDir, *dbFilePath) },
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
	// Logger setup
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

	// Config
	cobra.OnInitialize(func() {
		configDir, _ = newConfig(*dbFilePath)
	})

	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(FetchCmd)
	RootCmd.AddCommand(MockServerCmd)

	// MARK: Set up global flags
	RootCmd.PersistentFlags().BoolVar(&state.Flags.VerboseLogging, "verbose", state.Flags.VerboseLogging, "enable detailed logging")
	RootCmd.PersistentFlags().VarP(state.Flags.OutputFormat, "output", "o", "output format")

	// RootCmd.PersistentFlags().VarP(dbAdapter, "db-adapter", "a", "database adapter")
	RootCmd.PersistentFlags().StringVar(&configDir, "config-dir", configDir, "config directory")
	RootCmd.PersistentFlags().StringVar(dbFilePath, "db-url", *dbFilePath, "database url")
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

	globalTimeout         = 1 * time.Minute
	plugins               = make(map[string]reflect.Value)
	dbFilePath    *string = helpers.PointerTo(fmt.Sprintf("%s.sqlite3", entrypoint))
	sqlSchema     []byte
	database      dbsql.SqlBackend

	configDir string = fmt.Sprintf("$HOME/.%s", entrypoint)
)
