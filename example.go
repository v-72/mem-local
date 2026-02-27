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
	ttl := 1000
	setStatus := s.Set("TestKey", "TestValue", &ttl)
	fmt.Println("\nsetStatus:", setStatus)
	val, _ := s.Get("TestKey")
	fmt.Println(val)
	fmt.Println("value from cache", val)

	deleteStatus := s.Delete("TestKey")
	fmt.Println("\ndeleteStatus:", deleteStatus)
	val, ok := s.Get("TestKey")
	if !ok {
		fmt.Println("Key not found after deletion")
	} else {
		fmt.Println("Value after deletion:", val)
	}
}
