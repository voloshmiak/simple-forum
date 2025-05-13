package config

import "github.com/joho/godotenv"

func LoadEnv() error {
	// load environment variables
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}
