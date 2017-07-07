package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/generator"
	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Initiates grpc transport after creating the protobuf",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide the service name")
			return
		}
		g := generator.NewGRPCInitGenerator()
		err := g.Generate(args[0])
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

func init() {
	initCmd.AddCommand(grpcCmd)
}
