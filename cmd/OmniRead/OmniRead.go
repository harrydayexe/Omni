package main

import (
	"fmt"
	"github.com/harrydayexe/Omni/internal/core"
)

func main() {
	user := core.User{
		Name: "John",
		Age:  25,
	}

	fmt.Println("Hello, Read!")
	fmt.Println(user.Name)
}
