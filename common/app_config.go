package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var (
	// ApplicationConfigFileName 系统配置文件名称，默认 application.yaml
	ApplicationConfigFileName = "application.yaml"
)

type (
	// 系统配置文件
	AppConfig struct {
		Server struct {
			Addr string `yaml:"addr"`
			Mode string `yaml:"mode"`
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
			Cache struct {
				CleanupInterval time.Duration `yaml:'cleanupInterval'`
			} `yaml:"cache"`
			TokenCheck struct {
				CheckAtServer bool   `yaml:"checkAtServer"`
				CheckEndpoint string `yaml:"checkEndpoint"`
			} `yaml:"tokenCheck"`
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

// SetApplicationConfigFileName 设置系统配置文件名称
func SetApplicationConfigFileName(fileName string) {
	ApplicationConfigFileName = fileName
}

// ApplicationConfig 获取系统配置
func ApplicationConfig() *AppConfig {
	if applicationConfig != nil {
		return applicationConfig
	}

	file, err := ioutil.ReadFile(ApplicationConfigFileName)
	if err != nil {
		PanicPolestarError(ERR_SYS_ERROR, "系统配置文件读取错误！"+err.Error())
	}

	applicationConfig = &AppConfig{}
	err = yaml.Unmarshal(file, applicationConfig)
	if err != nil {
		PanicPolestarError(ERR_SYS_ERROR, "系统配置文件解析错误！"+err.Error())
	}
	return applicationConfig
}
