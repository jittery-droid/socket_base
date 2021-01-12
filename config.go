package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "ian",
		Password: "password",
		Name:     "chat_db",
	}
}

type Config struct {
	Port      int            `json:"port"`
	Env       string         `json:"env"`
	Pepper    string         `json:"pepper"`
	Database  PostgresConfig `json:"database"`
	JWTSecret string         `json:"jwt_secret"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port:      5000,
		Env:       "dev",
		Pepper:    "secret-random-string",
		JWTSecret: "silly-string",
		Database:  DefaultPostgresConfig(),
	}
}

func LoadConfig(configReq bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if configReq {
			panic(err)
		}
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .config")
	return c
}
