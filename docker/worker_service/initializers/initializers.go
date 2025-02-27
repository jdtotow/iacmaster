package initializers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	pwd, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(pwd, ".env"))
	if err != nil {
		log.Fatal("Could not load .env")
	}
}
