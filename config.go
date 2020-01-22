package telewatch

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var (
	ConfigFileNotExist = fmt.Errorf("token.yml file not found")
)

type TokenConfig struct {
	Token  string `yaml:"token"`
	ChatId int64  `yaml:"chatId"`
}

func (config TokenConfig) Save() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFilePath, configData,0644)

	return err
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := path.Join(home, ".config", "telewatch")
	configFilePath := path.Join(configPath, "token.yaml")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		configFilePath = path.Join(configPath, "token.yml")
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			return "", ConfigFileNotExist
		}
	}

	return configFilePath, nil
}

func LoadConfigurationFromHomeDirectory() (TokenConfig, error) {
	ret := TokenConfig{}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return ret, err
	}
	configFile, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return ret, err
	}

	err = yaml.Unmarshal(configFile, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
