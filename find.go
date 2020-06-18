package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sync"
)

type findOptions struct {
	startingDir string
	filename    string
}

var wg sync.WaitGroup

func find(args []string) {
	opts := findOptions{}
	err := opts.splitArgs("find", args)
	if err != nil {
		fmt.Println(err)
	}

	// Recursively check directories for filename
	opts.checkDir(opts.startingDir)
}

func (opts findOptions) checkDir(dir string) {
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range fileList {
		if file.IsDir() {
			opts.checkDir(dir + "/" + file.Name())
		} else {
			if file.Name() == opts.filename {
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
