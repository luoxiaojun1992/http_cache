package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func Env(key, defaultValue string) string {
	val := os.Getenv(key)

	if len(val) > 0 {
		return val
	}

	return defaultValue
}

func EnvInt(key string, defaultValue int) int {
	val := Env(key, "")
	if len(val) > 0 {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
	}

	return defaultValue
}
