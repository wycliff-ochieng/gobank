package main

import "fmt"

func main() {
	fmt.Println("Hawayu")
	server := NewAPIServer(":3000")
	server.Run()
}
