package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Default struct {
		SSHHost       string `mapstructure:"ssh_host"`
		SSHPort       int    `mapstructure:"ssh_port"`
		SSHUser       string `mapstructure:"ssh_user"`
		SSHPrivateKey string `mapstructure:"ssh_private_key"`
	} `mapstructure:"default"`
	Database struct {
		ROTraffic struct {
			Server   string `mapstructure:"server"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
			SSLMode  string `mapstructure:"sslmode"`
		} `mapstructure:"ro-traffic"`
	} `mapstructure:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		configPath = filepath.Join(homeDir, ".ssh", "dnc_db_info")
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return &config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.ROTraffic.Server,
		c.Database.ROTraffic.Port,
		c.Database.ROTraffic.Username,
		c.Database.ROTraffic.Password,
		c.Database.ROTraffic.Database,
		c.Database.ROTraffic.SSLMode,
	)
}
