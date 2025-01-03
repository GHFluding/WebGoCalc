package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"` //env-required:"true"
	Storage    `yaml:"storage"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
type Storage struct {
	Path     string `yaml:"path" env-required:"true"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"qwerty"`
	DBName   string `yaml:"dbname" env-default:"storage"`
	AsString string
}

// load config
func MustLoad() *Config {
	// if err := godotenv.Load("../../.env"); err != nil {
	// 	log.Print("No .env file found")
	// }
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("config_path is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config_path is not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalln("cannot read config")
	}
	cfg.Storage.AsString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName)
	return &cfg
}
