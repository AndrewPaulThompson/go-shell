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

	// Get the Home Directory and change to it so we start there
	homeDir, _ := os.UserHomeDir()
	// err := os.Chdir(homeDir)

	// if err != nil {
	// 	log.Panic(err)
	// }

	for {
		// If the home directory is in the working directory path, replace it with ~
		workingDir, _ := os.Getwd()
		dir := strings.Replace(workingDir, homeDir, "~", 1)

		// Replace \ for /
		dir = strings.Replace(dir, "\\", "/", -1)

		// Command line prefix
		fmt.Print(dir + "$ ")

		// Read the input until newline
		input, err := reader.ReadString('\n')

		// Convert CRLF to LF (for Windows)
		input = strings.Replace(input, "\r\n", "", -1)
		args := strings.Split(input, " ")

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if val, ok := m[args[0]]; ok {
			val(args[1:])
		} else {
			fmt.Printf("Command %s not found\n", args[0])
		}
	}
}

func registerFunctions() map[string]func([]string) {
	// probably need to change this since the functions will need to take args etc
	m := make(map[string]func([]string))
	m["ls"] = ls
	m["cd"] = cd
	m["pwd"] = pwd
	m["find"] = find

	return m
}
