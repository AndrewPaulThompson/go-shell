package main

import (
	"fmt"
	"os"
)

func pwd(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(wd)
}
