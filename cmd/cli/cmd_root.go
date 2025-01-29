package cli

import (
	"embed"

	toolkit "github.com/oleoneto/go-toolkit/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI, a cli tool.",
	// PersistentPreRun: BeforeHook(state), // TODO: Replace command
	// PersistentPostRun: AfterHook(state), // TODO: Replace command
	PostRun: func(cmd *cobra.Command, args []string) {
		if buildHash != "" {
			log.Debug("build", buildHash)
		}
	},
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func Execute(vfs embed.FS, buildHash string) error {
	// virtualFS = vfs
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(RequestCmd)

	// MARK: Set up global flags
	RootCmd.PersistentFlags().BoolVar(&state.Flags.VerboseLogging, "verbose", state.Flags.VerboseLogging, "enable detailed logging")
	RootCmd.PersistentFlags().VarP(state.Flags.OutputFormat, "output", "o", "output format")
	RootCmd.PersistentFlags().StringVarP(&state.Flags.OutputTemplate, "output-template", "y", state.Flags.OutputTemplate, "template (used when output format is 'gotemplate')")

	RootCmd.Flags().VarP(&state.Flags.File, "file", "f", "")

	// Migrator configuration
	RootCmd.PersistentFlags().VarP(dbAdapter, "db-adapter", "a", "database adapter")
	RootCmd.PersistentFlags().BoolVar(&state.Flags.TimeExecutions, "time", state.Flags.TimeExecutions, "time executions")
}

var (
	// wherein dictionary files are stored
	virtualFS embed.FS

	buildHash string

	dbAdapter = &toolkit.FlagEnum{
		Allowed: []string{"postgresql", "sqlite3"},
		Default: "sqlite3",
	}

	state = toolkit.CommandState{
		Flags: toolkit.CommandFlags{
			OutputTemplate: "",
			OutputFormat: &toolkit.FlagEnum{
				Allowed: []string{"plain", "json", "yaml", "table", "gotemplate", "silent"},
				Default: "plain",
			},
		},
	}

	loggerFunc func(...any) = log.Infoln
	debugFunc  func(...any) = log.Debugln
)
