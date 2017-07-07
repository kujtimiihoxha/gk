package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/generator"
	"github.com/spf13/cobra"
)

// grpc_addCmd represents the grpc_add command
var grpc_addCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Add grpc transport to service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide the service name")
			return
		}
		g := generator.NewAddGRPCGenerator()
		err := g.Generate(args[0])
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

func init() {
	addCmd.AddCommand(grpc_addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpc_addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpc_addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
