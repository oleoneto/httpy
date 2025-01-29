package cli

import (
	"fmt"
	"log"

	"github.com/mitchellh/go-homedir"
	"github.com/oleoneto/go-toolkit/files"
	"github.com/spf13/viper"
)

func initConfig() {
	if state.Flags.CLIConfig != "" {
		viper.SetConfigFile(state.Flags.CLIConfig)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		path := fmt.Sprintf("%v/%s", home, state.Flags.ConfigDir.Name)

		viper.AddConfigPath(path)
		viper.SetConfigName("config")
	}

	// NOTE: File does not exist... create one!
	if err := viper.ReadInConfig(); err != nil {
		home, herr := homedir.Dir()
		if herr != nil {
			log.Fatal(err)
		}

		if f := state.Flags.ConfigDir.Create(files.FileGenerator{}, home); len(f) == 0 {
			log.Fatal("Cannot read config. Hint: You may need to run `init` to create the config file")
		}
	}
}
