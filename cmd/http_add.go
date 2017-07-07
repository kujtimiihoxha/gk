package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/generator"
	"github.com/spf13/cobra"
)

// http_addCmd represents the http_add command
var http_addCmd = &cobra.Command{
	Use:   "http",
	Short: "Add http transport to service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide the service name")
			return
		}
		g := generator.NewAddHttpGenerator()
		err := g.Generate(args[0])
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

func init() {
	addCmd.AddCommand(http_addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// http_addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// http_addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
