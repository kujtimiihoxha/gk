package fs

import (
	"github.com/spf13/afero"
	"os"
	"github.com/spf13/viper"
	"github.com/Songmu/prompter"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type FileSystem interface {
	init(dir string)
	ReadFile(path string) (string, error)
	WriteFile(path string, data string, force bool) error
	Mkdir(path string) error
	MkdirAll(path string) error
	FilePathSeparator() string
	Exists(path string) (bool, error)
	Walk(root string, fc func(path string, info os.FileInfo, err error) error) error
}

type DefaultFs struct {
	fs afero.Fs
}

func (f *DefaultFs) init(dir string) {
	var inFs afero.Fs
	if viper.Get("testing") {
		inFs = afero.NewMemMapFs()
	} else {
		inFs = afero.NewOsFs()
	}
	if dir != "" {
		f.fs = afero.NewBasePathFs(inFs, dir)
	} else {
		f.fs = inFs
	}
}
func (f *DefaultFs) ReadFile(path string) (string, error) {
	d, err := afero.ReadFile(f.fs, path)
	return string(d), err
}

func (f *DefaultFs) WriteFile(path string, data string, force bool) error {

	if b, _ := f.Exists(path); b && !(viper.GetBool("gk.force_override") || force) {
		s, _ := f.ReadFile(path)
		if s == data {
			logrus.Warnf("`%s` exists and is identical it will be ignored", path)
			return nil
		}
		b := prompter.YN(fmt.Sprintf("`%s` already exists do you want to override it ?", path), false)
		if !b {
			return nil
		}
	}
	return afero.WriteFile(f.fs, path, []byte(data), os.ModePerm)
}

func (f *DefaultFs) Mkdir(path string) error {
	return f.fs.Mkdir(path, os.ModePerm)
}

func (f *DefaultFs) MkdirAll(path string) error {
	return f.fs.MkdirAll(path, os.ModePerm)
}
func (f *DefaultFs) FilePathSeparator() string {
	return afero.FilePathSeparator
}
func (f *DefaultFs) Exists(path string) (bool, error) {
	return afero.Exists(f.fs, path)
}
func (f *DefaultFs) Walk(root string, fc func(path string, info os.FileInfo, err error) error) error {
	return afero.Walk(f.fs, root, fc)
}
func NewDefaultFs(dir string) *DefaultFs {
	dfs := &DefaultFs{}
	dfs.init(dir)
	return dfs	
}
