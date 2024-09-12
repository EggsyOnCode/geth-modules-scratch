package main

import (
	"fmt"
	"time"
)

func main() {
	// Assign the string to an interface
	var i interface{} = "umair"

	// Now perform a type assertion to extract the string value
	if actual, ok := i.(string); ok {
		fmt.Println("actual: ", actual) // Will print: actual: umair
	} else {
		fmt.Println("Type assertion failed")
	}
}

func giveTime(ch chan<- int64) {
	time.Sleep(1 * time.Second)
	ch <- time.Now().UnixMicro()
}
