package main

import "github.com/google/uuid"

func GenerateSerial() string {
	return uuid.New().String()
}
