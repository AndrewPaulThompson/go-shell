package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	m := registerFunctions()

	for {
		// Command line prefix
		fmt.Print("-> ")

		// Read the input until newline
		input, err := reader.ReadString('\n')

		// Convert CRLF to LF (for Windows)
		input = strings.Replace(input, "\r\n", "", -1)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if val, ok := m[input]; ok {
			val()
		} else {
			fmt.Printf("Command %s not found\n", input)
		}
	}
}

func registerFunctions() map[string]func() {
	// probably need to change this since the functions will need to take args etc
	m := make(map[string]func())
	m["ls"] = ls
	m["cd"] = cd

	return m
}
