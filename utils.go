package main

import (
	"log"
	"os"
)

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == "" {
			log.Fatal("Missing mandatory env variable : " + key)
		}
		return defaultValue
	}
	return value
}