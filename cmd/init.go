package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initiates a service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logrus.Error("You must provide the service name")
			return
		}
		if viper.GetString("gk_transport") == "" {
			viper.Set("gk_transport", viper.GetString("default_transport"))
		}
		gen := generator.NewServiceInitGenerator()
		err := gen.Generate(args[0])
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("transport", "t", "", "Specify the transport you want to initiate for the service")
	viper.BindPFlag("gk_transport", initCmd.Flags().Lookup("transport"))

}
