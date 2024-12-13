package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type AppConfig struct {
	MYSQL_URL string
	API_PORT  string

	PrivateKey_     string
	MetaNodeVersion string

	NodeConnectionAddress string

	StorageAddress           string
	StorageConnectionAddress string

	DnsLink_ string

	FicamAddress string
	FicamABIPath string

	EmailAdmin string
	TemplateEmailOrderPath string
	
	AWSConfig  AWSConfig `json:"aws_config"`
}
type CredentialsConfig struct {
	Id        string `json:"id"`
	SecretKey string `json:"secret_key"`
	Token     string `json:"token"`
}

type AWSConfig struct {
	Region            string            `json:"region"`
	CredentialsConfig CredentialsConfig `json:"credentials_config"`
	SenderAddress     string
}

var Config *AppConfig

func LoadConfig(configFilePath string) (*AppConfig, error) {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
func (c *AppConfig) DnsLink() string {
	return c.DnsLink_
}
