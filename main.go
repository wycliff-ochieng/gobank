package main

import (
	"fmt"

	"github.com/wycliff-ochieng/handler"
)

func main() {
	fmt.Println("Hawayu")
	server := handler.NewAPIServer(":3000")
	server.Run()
}
