package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/dolma"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check and print the signature",
	Run: func(cmd *cobra.Command, args []string) {
		configManager.LoadConfigOnly(dolmaConfig)

		logrus.SetLevel(logrus.DebugLevel)

		logrus.Infof("loaded config %#v\n", dolmaConfig)

		var err error
		wrapper := dolma.NewDolma()
		if err = wrapper.LoadBinary(dolmaConfig.TargetBinary); err != nil {
			logrus.Errorf("Error loading binary %s", err)
			return
		}

		logrus.Infof("Signature %#v", wrapper.Signature)
	},
}

func init() {
	RootCmd.AddCommand(CheckCmd)
}
