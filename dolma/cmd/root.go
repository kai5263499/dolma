package cmd

import (
	"github.com/CrowdStrike/fortio"
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
)

var (
	configManager *fortio.Manager
	dolmaConfig   *Config
)

var RootCmd = &cobra.Command{
	Use:   "dolma",
	Short: "Dolma is utility for adding static resources to a go binary",
}

func init() {
	configManager = fortio.NewConfigManagerWithRootCmd(RootCmd, &fortio.CmdLineConfigLoader{})

	dolmaConfig = &Config{}

	err := configManager.CreateCommandLineFlags(dolmaConfig, RootCmd)
	if err != nil {
		logrus.Errorf("Unable to load config - %v", err)
	}
}
