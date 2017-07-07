package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/generator"
	"github.com/spf13/cobra"
)

var thriftCmd = &cobra.Command{
	Use:   "thrift",
	Short: "Initiates thrift transport after creating the thrift service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide the service name")
			return
		}
		g := generator.NewThriftInitGenerator()
		err := g.Generate(args[0])
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

func init() {
	initCmd.AddCommand(thriftCmd)
}
