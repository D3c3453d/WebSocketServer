package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type ServerConfig struct {
	Bind                string `mapstructure:"server_bind,omitempty"`
	Port                string `mapstructure:"server_port,omitempty"`
	Schema              string `mapstructure:"server_schema,omitempty"`
	ReadTimeout         int    `mapstructure:"server_read_timeout,omitempty"`
	WriteTimeout        int    `mapstructure:"server_write_timeout,omitempty"`
	TokenExpireDuration int    `mapstructure:"server_token_expire_duration,omitempty"`
	Salt                string `mapstructure:"server_salt,omitempty"`
	TokenKey            string `mapstructure:"server_token_key,omitempty"`
	LogLevel            string `mapstructure:"server_log_level" `
	LogShowPath         bool   `mapstructure:"server_log_show_path"`
	ExternalDomainName  string `mapstructure:"server_external_domain_name"`
	ExternalPort        string `mapstructure:"server_external_port"`
}

type DbConfig struct {
	Name              string `mapstructure:"db_name,omitempty"`
	User              string `mapstructure:"db_user,omitempty"`
	Password          string `mapstructure:"db_password,omitempty"`
	Type              string `mapstructure:"db_type,omitempty"`
	Host              string `mapstructure:"db_host,omitempty"`
	Port              string `mapstructure:"db_port,omitempty"`
	SSLMode           string `mapstructure:"db_ssl_mode"`
	MigrationsDirPath string `mapstructure:"db_migration_path"`
}

func InitConfig() {

	wddir, _ := os.Getwd()
	appDir := strings.Split(wddir, "/app")[0] + "/app"

	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../app")
	viper.AddConfigPath(appDir)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetEnvPrefix("ARCHIVE")

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Errorf("Error read config file: default:  %v", err)
		os.Exit(1)
	}

	viper.AutomaticEnv()
	logrus.Info("Config prepared")

	logrus.Info("Loaded config database: %v", GetDbConfig())
	logrus.Info("Loaded config server: %v", GetServerConfig())

}

func GetServerConfig() *ServerConfig {
	serverConfig := &ServerConfig{}

	if err := viper.Unmarshal(&serverConfig); err != nil {
		logrus.Panicf("Can't load server config. %v", err)
		os.Exit(1)
	}

	return serverConfig
}

func GetDbConfig() *DbConfig {
	dbConfig := &DbConfig{}

	if err := viper.Unmarshal(&dbConfig); err != nil {
		logrus.Panicf("Can't load db_adapter config. %v", err)
		os.Exit(1)
	}

	return dbConfig
}
