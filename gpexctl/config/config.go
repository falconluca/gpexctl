package config

import (
	"github.com/fatih/color"
	log "github.com/golang/glog"
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

	Conf *Config
)

type Config struct {
	ApiKey    string `yaml:"apiKey"`
	HTTPProxy string `yaml:"HTTPProxy"`
}

func InitConfig() {
	var err error
	if Conf, err = LoadConfig("../config.yaml"); err != nil {
		log.Errorf("%#+v", err)
	}
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
