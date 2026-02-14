package main

import (
	"fmt"
	"github.com/v-72/mem-local/store"
)

func main() {
	fmt.Println("init")
	s := &store.MemLocalStore{}
	setStatus := s.Set("test", "test")
	fmt.Println("setStatus:", setStatus)
}
