package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/v-72/mem-local/store"
)

func generateSecureString(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func main() {
	fmt.Println("init")
	s := &store.MemLocalStore{}
	err := s.Init()

	if err != nil {
		fmt.Println("Error initializing local store:", err)
		return
	}

	key, _ := generateSecureString(8)
	value, _ := generateSecureString(64)
	fmt.Printf("Generated key: %s\n", key)
	fmt.Printf("Generated value: %s\n", value)

	setStatus := s.Set(key, value)
	fmt.Println("\nsetStatus:", setStatus)
	val, _ := s.Get(key)
	fmt.Println(val)
	fmt.Println("value from cache", val)
}
