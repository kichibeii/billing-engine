package config

import "github.com/billing-engine/internal/service"

type Config struct {
	App      App
	Database DatabaseConfig
}

type App struct {
	Host string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Driver               string
	Name                 string
	User                 string
	Password             string
	Host                 string
	Port                 string
	AdditionalParameters string
}

type AppConfig struct {
	Config  *Config
	Service service.ServiceInterface
}
