package main

import (
	"fmt"

	"github.com/v-72/mem-local/store"
)

func main() {
	fmt.Println("init")
	s := &store.MemLocalStore{}
	err := s.Init()

	if err != nil {
		fmt.Println("Error initializing local store:", err)
		return
	}
	setStatus := s.Set("TestKey", "TestValue")
	fmt.Println("\nsetStatus:", setStatus)
	val, _ := s.Get("TestKey")
	fmt.Println(val)
	fmt.Println("value from cache", val)
}
