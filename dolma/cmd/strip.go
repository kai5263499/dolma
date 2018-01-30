package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/dolma"
)

var StripCmd = &cobra.Command{
	Use:   "strip",
	Short: "Strip sections and signature",
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
		if err = wrapper.StripSections(); err != nil {
			logrus.Errorf("Error stripping sections + signature %s", err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(StripCmd)
}
