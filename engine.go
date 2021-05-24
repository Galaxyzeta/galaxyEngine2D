package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	buffer := make([]byte, 1)
	go func() {
		for {
			fmt.Print("1")
			_, err := os.Stdin.Read(buffer)
			fmt.Print(string(buffer))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	time.Sleep(time.Second * 5)
}
