package cmd

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

		logrus.Infof("loaded config %#v\n", dolmaConfig)

		flatBuilder := flatbuffers.NewBuilder(0)
		flatBuilder.Reset()
		ipPosition := flatBuilder.CreateString("test")
		generated.EndpointStart(flatBuilder)
		generated.EndpointAddIp(flatBuilder, ipPosition)
		endpointPosition := generated.EndpointEnd(flatBuilder)
		flatBuilder.Finish(endpointPosition)
		logrus.Infof("flatbuffer: %#v", string(flatBuilder.Bytes[flatBuilder.Head():]))

		wrapper := dolma.NewDolma()
		wrapper.LoadBinary(dolmaConfig.TargetBinary)
		wrapper.SaveSignature()
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
