package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting sleep...")

	time.Sleep(1200 * time.Second)

	fmt.Println("Sleep's over...")
}
