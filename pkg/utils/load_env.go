package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadENV() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error load .env file : %s", err.Error())
	}

	return nil
}
