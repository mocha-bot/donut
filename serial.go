package main

import "github.com/gofrs/uuid/v5"

func GenerateSerial() string {
	return uuid.Must(uuid.NewV4()).String()
}
