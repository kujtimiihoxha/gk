package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/fs"
	"github.com/kujtimiihoxha/gk/templates"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gk",
	Short: "A generator for go-kit that helps you create/update boilerplate code",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().Bool("testing", false, "If testing the generator.")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "If you want to se the debug logs.")
	RootCmd.PersistentFlags().BoolP("force", "f", false, "Force overide existing files without asking.")
	RootCmd.PersistentFlags().String("folder", "", "If you want to specify the base folder of the project.")
	viper.BindPFlag("gk_testing", RootCmd.PersistentFlags().Lookup("testing"))
	viper.BindPFlag("gk_folder", RootCmd.PersistentFlags().Lookup("folder"))
	viper.BindPFlag("gk_force", RootCmd.PersistentFlags().Lookup("force"))
	viper.BindPFlag("gk_debug", RootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	initViperDefaults()
	viper.SetFs(fs.NewDefaultFs("").Fs)
	viper.SetConfigFile("gk.json")
	if viper.GetBool("gk_debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Info("No config file found initializing the project with the default config file.")
		te := template.NewEngine()
		st, err := te.Execute("gk.json", nil)
		if err != nil {
			logrus.Panic(err)
		}
		err = fs.Get().WriteFile("gk.json", st, false)
		if err != nil {
			logrus.Panic(err)
		}
		initConfig()
	}
}
func initViperDefaults() {
	viper.SetDefault("service.path", "{{toSnakeCase .ServiceName}}"+afero.FilePathSeparator+"pkg"+afero.FilePathSeparator+"service")
	viper.SetDefault("service.file_name", "service.go")
	viper.SetDefault("service.interface_name", "{{toUpperFirstCamelCase .ServiceName}}Service")
	viper.SetDefault("service.struct_name", "stub{{toCamelCase .ServiceName}}Service")
	viper.SetDefault("middleware.file_name", "middleware.go")
	viper.SetDefault("endpoints.path", "{{toSnakeCase .ServiceName}}"+afero.FilePathSeparator+"pkg"+afero.FilePathSeparator+"endpoints")
	viper.SetDefault("endpoints.file_name", "endpoints.go")
	viper.SetDefault("transport.path", "{{toSnakeCase .ServiceName}}"+afero.FilePathSeparator+"pkg"+afero.FilePathSeparator+"{{.TransportType}}")
	viper.SetDefault("transport.file_name", "handler.go")
	viper.SetDefault("default_transport", "http")
}
