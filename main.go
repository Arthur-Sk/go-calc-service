package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting sleep...")

	time.Sleep(12000 * time.Second)

	fmt.Println("Sleep's over...")
}
