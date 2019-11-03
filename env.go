package cronic

import (
	"os"
)

// Like os.Getenv, but accepts a default value if the variable is not set
func Getenv(name, defaultValue string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	}

	return defaultValue
}
