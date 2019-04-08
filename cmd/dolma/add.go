package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/dolma"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files and folders to a section",
	Run: func(cmd *cobra.Command, args []string) {
		configManager.LoadConfigOnly(dolmaConfig)

		logrus.SetLevel(logrus.DebugLevel)

		logrus.Debugf("loaded config %#v\n", dolmaConfig)

		wrapper := dolma.NewDolma()
		wrapper.LoadBinary(dolmaConfig.TargetBinary)

		wrapper.AddContent(dolmaConfig.SectionPrefix, dolmaConfig.SectionContent)

		// wrapper.SaveSignature()
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
