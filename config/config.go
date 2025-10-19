package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Server   *Server             `yaml:"service"`
	Services map[string]*Service `yaml:"services"`
	Etcd     *Etcd               `yaml:"etcd"`
	Kafka    *Kafka              `yaml:"kafka"`
}

type Server struct {
	Port      string `yaml:"port"`
	Version   string `yaml:"version"`
	JwtSecret string `yaml:"jwtSecret"`
	Metrics   string `yaml:"metrics"`
}

type Service struct {
	Name    string   `yaml:"name"`
	Addr    []string `yaml:"addr"`
	Metrics []string `yaml:"metrics"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Kafka struct {
	Address []string `yaml:"address"`
}

var Conf Config

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}
