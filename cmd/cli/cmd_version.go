package cli

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/oleoneto/go-toolkit/cli/formatters"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of the CLI",
	// PersistentPreRun:  BeforeHook(state),
	// PersistentPostRun: AfterHook(state),
	Run: func(cmd *cobra.Command, args []string) {
		state.Writer.Print(version)
	},
}

type Version struct {
	Major string `json:"major" yaml:"major" toml:"major" csv:"major"`
	Minor string `json:"minor" yaml:"minor" toml:"minor" csv:"minor"`
	Patch string `json:"patch" yaml:"patch" toml:"patch" csv:"patch"`
}

func (v *Version) String() string { return fmt.Sprintf(`cli %v.%v.%v`, v.Major, v.Minor, v.Patch) }

func (v *Version) TableWriter() table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(nil) // Delegate printing to gout tool

	t.SetTitle("version")
	t.AppendHeader(table.Row{"major", "minor", "patch"})
	t.AppendRow(table.Row{v.Major, v.Minor, v.Patch})

	return t
}

var (
	version = &Version{Major: "0", Minor: "1", Patch: "0-alpha"}

	_ formatters.TableFormattable = (*Version)(nil)
)
