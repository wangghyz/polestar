package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	fileApplicationYaml = "application.yaml"
)

type (
	AppConfig struct {
		Server struct {
			Addr string `yaml:"addr"`
		} `yaml:"server"`
		Mysql struct {
			Url          string `yaml:"url"'`
			MaxIdleConns int    `yaml:"maxIdleConns"`
			LogMode      bool   `yaml:"logMode"`
		} `yaml:"mysql"`
		Redis struct {
			Url string `yaml:"url"`
			DB  int    `yaml:"db"`
		} `yaml:"redis"`
		Auth struct {
			Jwt struct {
				Secret string `yaml:"secret"`
			} `yaml:"jwt"`
			TokenExp struct {
				AccessToken  int `yaml:"accessToken"`
				RefreshToken int `yaml:"refreshToken"`
			} `yaml:"tokenExp"`
			AuthUris []authUriInfo `yaml:"authUris"`
			SkipUris []authUriInfo `yaml:"skipUris"`
		}
	}

	authUriInfo struct {
		Uri         string   `yaml:"uri"`
		Methods     []string `yaml:"methods"`
		Authorities []string `yaml:"authorities"`
	}
)

var applicationConfig *AppConfig = nil

func SetAppConfigFileName(fileName string) {
	fileApplicationYaml = fileName
}

func ApplicationConfig() *AppConfig {
	if applicationConfig != nil {
		return applicationConfig
	}

	_, err := os.Stat(fileApplicationYaml)
	if err != nil {
		return nil
	}

	file, err := ioutil.ReadFile(fileApplicationYaml)
	if err != nil {
		return nil
	}

	applicationConfig = &AppConfig{}
	err = yaml.Unmarshal(file, applicationConfig)
	if err != nil {
		applicationConfig = nil
		return nil
	}
	return applicationConfig
}
