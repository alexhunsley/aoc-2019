package main

import (
	"fmt"
	"time"
)

func sleepPrint() {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello wurst")
}

func mainX() {
	fmt.Println("Hello, playground")
	go sleepPrint()
	time.Sleep(2 * time.Second)
}
