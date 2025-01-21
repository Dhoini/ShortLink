package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateJWTSecret() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func main() {
	secret, err := generateJWTSecret()
	if err != nil {
		fmt.Println("Error generating secret:", err)
		return
	}
	fmt.Println("Generated JWT Secret:", secret)
}
