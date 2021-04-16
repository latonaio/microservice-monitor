package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v2"
)

type Env struct {
	Interval         int    `envconfig:"INTERVAL" default:"10"`
	WindowSize       int    `envconfig:"WINDOW_SIZE" default:"5"`
	ConfigDir        string `envconfig:"CONFIG_DIR" default:"./config"`
	SlackIncomingURL string `envconfig:"SLACK_URL"`
}

func GetEnv() (Env, error) {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		return Env{}, err
	}

	return env, nil
}

type device struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
}

type threshold struct {
	CPU    float64 `yaml:"cpu,omitempty"`
	Memory float64 `yaml:"memory,omitempty"`
}

type alert struct {
	Name      string    `yaml:"name"`
	Threshold threshold `yaml:"threshold"`
}

type AlertSetting struct {
	Device device  `yaml:"device"`
	Alert  []alert `yaml:"alert,omitempty"`
}

func LoadAlertSetting(env Env) (AlertSetting, error) {
	buf, err := ioutil.ReadFile(filepath.Join(env.ConfigDir, "alert-setting.yml"))
	if err != nil {
		return AlertSetting{}, err
	}

	var as AlertSetting
	err = yaml.Unmarshal(buf, &as)
	return as, err
}
