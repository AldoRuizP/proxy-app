package utils

import (
	env "github.com/joho/godotenv"
)

// LoadEnv loads envs
func LoadEnv() {
	env.Load()
}
