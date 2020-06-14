package main

import (
	"fmt"
	"os"
)

func cd(args []string) {
	// Return if more than 1 argument is passed
	if len(args) > 1 {
		fmt.Printf("Too many arguments: have %d, expected 1\n", len(args))
		return
	}

	// Else change to the specified directory
	err := os.Chdir(args[0])

	if err != nil {
		fmt.Println(err)
		return
	}
}
