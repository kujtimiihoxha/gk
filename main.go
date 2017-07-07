package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/kujtimiihoxha/gk/cmd"
	"github.com/kujtimiihoxha/gk/utils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func main() {
	viper.AutomaticEnv()
	gosrc := utils.GetGOPATH() + afero.FilePathSeparator + "src" + afero.FilePathSeparator
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
		return
	}
	if !strings.HasPrefix(pwd, gosrc) {
		logrus.Error("The project must be in the $GOPATH/src folder for the generator to work.")
		return
	}
	cmd.Execute()
}
