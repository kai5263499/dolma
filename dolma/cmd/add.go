package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/dolma"
)

var AddCmd = &cobra.Command{
	Use: "add",
	Short: "Add files and folders to a section",
	Run: func(cmd *cobra.Command, args []string) {
		configManager.LoadConfigOnly(dolmaConfig)

		logrus.SetLevel(logrus.DebugLevel)

		logrus.Infof("loaded config %#v\n", dolmaConfig)

		wrapper := dolma.NewDolma()
		wrapper.LoadBinary(dolmaConfig.TargetBinary)
		wrapper.SaveSignature()
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
