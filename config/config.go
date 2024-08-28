package config

import (
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/spf13/viper"
)

const DefaultYmlFile = "default.config.yml"

type (
	Config struct {
		APIPort int           `mapstructure:"apiPort" validate:"required"`
		Kafka   KafkaConfig   `mapstructure:"kafka" validate:"required"`
		MongoDB MongoDBConfig `mapstructure:"mongoDB" validate:"required"`
	}

	KafkaConfig struct {
		Brokers                    string `mapstructure:"brokers" validate:"required"`
		Topic                      string `mapstructure:"topic" validate:"required"`
		Group                      string `mapstructure:"group" validate:"required"`
		WithTls                    bool   `mapstructure:"withTls"`
		DialerTimeoutMilliseconds  int    `mapstructure:"dialerTimeoutMilliseconds"`
		MaxWaitTimeoutMilliseconds int    `mapstructure:"maxWaitTimeoutMilliseconds"`
		NumConsumers               int    `mapstructure:"numConsumers"`
	}

	MongoDBConfig struct {
		Host     string `mapstructure:"host" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		Database string `mapstructure:"database" validate:"required"`
		Port     int    `mapstructure:"port"`
		Options  string `mapstructure:"options"`
	}
)

func getFilePath(fileName string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return basepath + "/" + fileName
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	viper.SetConfigFile(getFilePath(DefaultYmlFile))

	if err := viper.MergeInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (m *MongoDBConfig) MongoURI() string {
	return "mongodb://" + m.User + ":" + m.Password + "@" + m.Host + ":" + strconv.Itoa(m.Port) + "/" + m.Database + "?" + m.Options
}
