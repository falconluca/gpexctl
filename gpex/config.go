package gpex

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var (
	Indicator = color.CyanString("==> ")

	YellowString = color.New(color.FgYellow).SprintFunc()
	RedString    = color.New(color.FgRed).SprintFunc()
	GreenString  = color.New(color.FgGreen).SprintFunc()
	CyanString   = color.New(color.FgCyan).SprintFunc()
)

type Config struct {
	ApiKey    string `yaml:"apiKey"`
	HttpProxy string `yaml:"httpProxy"`
}

func LoadConfig(file string) (*Config, error) {
	f, _ := filepath.Abs(file)

	y, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrap(err, "read yaml file failed")
	}

	var config Config
	if err = yaml.Unmarshal(y, &config); err != nil {
		return nil, errors.Wrap(err, "parse yaml file failed")
	}
	return &config, nil
}
