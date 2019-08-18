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

func (s sphinxRow) short() string {
	return fmt.Sprintf("%s %s", s.Name, s.preAt().String())
}

func (s sphinxRow) preAt() time.Time {
	return time.Unix(s.PreAt, 0)
}