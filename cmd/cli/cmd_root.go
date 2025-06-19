package cli

import (
	"embed"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	"github.com/oleoneto/go-toolkit/helpers"
	"github.com/oleoneto/httpy/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed schemas
var schemas embed.FS

var entrypoint = "httpy"

var RootCmd = &cobra.Command{
	Use:   entrypoint,
	Short: "HTTPy, a CLI tool for programmatically managing collections of HTTP requests",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func Execute(config pkg.CLIConfig) error {
	httpyFlags.plugins = config.Plugins

	if config.DefaultTimeout != nil {
		httpyFlags.timeout = *config.DefaultTimeout
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
		httpyFlags.configDir, _ = newConfig(*httpyFlags.dbFilePath)
	})

	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(FetchCmd)
	RootCmd.AddCommand(MockServerCmd)

	// MARK: Set up global flags
	RootCmd.PersistentFlags().BoolVar(&state.Flags.VerboseLogging, "verbose", state.Flags.VerboseLogging, "enable detailed logging")
	RootCmd.PersistentFlags().VarP(state.Flags.OutputFormat, "output", "o", "output format")

	// RootCmd.PersistentFlags().VarP(dbAdapter, "db-adapter", "a", "database adapter")
	RootCmd.PersistentFlags().StringVar(&httpyFlags.configDir, "config-dir", httpyFlags.configDir, "config directory")
	RootCmd.PersistentFlags().StringVar(httpyFlags.dbFilePath, "db-url", *httpyFlags.dbFilePath, "database url")
	RootCmd.PersistentFlags().BoolVar(&state.Flags.TimeExecutions, "time", state.Flags.TimeExecutions, "time executions")
	RootCmd.PersistentFlags().BoolVar(&httpyFlags.ephemeral, "ephemeral", httpyFlags.ephemeral, "halt recording of any output")
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

	httpyFlags = HTTPyFlags{
		timeout:    1 * time.Minute,
		plugins:    make(map[string]reflect.Value),
		dbFilePath: helpers.PointerTo(fmt.Sprintf("%s.sqlite3", entrypoint)),
		configDir:  fmt.Sprintf("$HOME/.%s", entrypoint),
	}
)
