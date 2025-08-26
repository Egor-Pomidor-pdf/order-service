package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env:"ENV"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
	Topic   string   `yaml:"topic"   env:"KAFKA_TOPIC"`
	GroupID string   `yaml:"group_id" env:"KAFKA_GROUP"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"     env:"DB_HOST"`
	Port     int    `yaml:"port"     env:"DB_PORT"`
	Name     string `yaml:"name"     env:"DB_NAME"`
	User     string `yaml:"user"     env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}

type ServerConfig struct {
	Address string `yaml:"address" env:"HTTP_ADDR"`
}


func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	var cfg Config
	var err error

	if configPath != "" {
		err = cleanenv.ReadConfig(configPath, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
