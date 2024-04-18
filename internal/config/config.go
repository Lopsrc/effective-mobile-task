package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Env     string `yaml:"env" env-default:"local"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`	
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
		} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var instance *Config
var once sync.Once

func GetConfig(pathConfig string) *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
			panic(err)
		}
	})
	return instance
}