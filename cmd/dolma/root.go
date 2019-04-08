package main

import (
	"github.com/CrowdStrike/fortio"
	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/dolma/domain"
	"github.com/spf13/cobra"
)

var (
	configManager *fortio.Manager
	dolmaConfig   *domain.Config
)

var RootCmd = &cobra.Command{
	Use:   "dolma",
	Short: "Dolma is utility for adding static resources to a go binary",
}

func init() {
	configManager = fortio.NewConfigManagerWithRootCmd(RootCmd, &fortio.CmdLineConfigLoader{})

	dolmaConfig = &domain.Config{}

	err := configManager.CreateCommandLineFlags(dolmaConfig, RootCmd)
	if err != nil {
		logrus.Errorf("Unable to load config err=%v", err)
	}
}
