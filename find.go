package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

type findOptions struct {
	startingDir string
	filename    string
}

func find(args []string) {
	// Parse flags, arguments
	opts := findOptions{}
	err := opts.splitArgs("find", args)
	if err != nil {
		fmt.Println(err)
	}

	// Recursively check directories for filename
	opts.checkDir(opts.startingDir)
}

func (opts findOptions) checkDir(dir string) {
	// Read the current directory
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	// For each file/directory in this directory
	for _, file := range fileList {
		if file.IsDir() {
			// If the file is a directory, search it
			opts.checkDir(dir + "/" + file.Name())
		} else {
			// Else if the filename is what we're searching for
			if file.Name() == opts.filename {
				// Print the filepath from the starting directory
				fmt.Println(dir + "/" + file.Name())
			}
		}
	}

	return
}

func (opts *findOptions) splitArgs(command string, args []string) error {
	// For find, the first argument is always the starting directory
	opts.startingDir = args[0]

	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.StringVar(&opts.filename, "name", "", "Search pattern")

	//  Parse arguments into flags
	return fs.Parse(args[1:])
}
