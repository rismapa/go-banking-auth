package domain

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	logger "github.com/rismapa/go-banking-lib/config"
)

type Config struct {
	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`

	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructur:"host"`
		API  string `mapstructure:"apikey"`
	} `mapstructure:"server"`

	DB struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Database string `mapstructure:"name"`
	} `mapstructure:"database"`
}

/*
 * Implemtasi database dengan config dari .env
 */
func (c *Config) GetDatabaseENVConfig() string {
	err := godotenv.Load(".env")
	if err != nil {
		logger.GetLog().Fatal().Err(err).Msg("Error loading .env file")
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
}
