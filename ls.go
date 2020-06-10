package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type opts struct {
	includeHidden bool
	longListing   bool
}

func ls(args []string) {
	opts := handleArgs(args)

	// Get the current directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Get files in the given directory
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// If this is a hidden file/directory AND we shouldn't show it, skip
		if string(file.Name()[0]) == "." && opts.includeHidden == false {
			continue
		}

		// Use long listing format
		if opts.longListing {
			fmt.Printf("%s %d %s %s\n", file.Mode(), file.Size(), file.ModTime().Format("Jan 2 15:04"), file.Name())
			continue
		}

		// Default to just printing the name
		fmt.Println(file.Name())
	}
}

func handleArgs(args []string) opts {
	opts := opts{}

	// For each additional argument
	for _, i := range args {
		// If the first character is "-", and the second character is NOT "-"
		// Flags have been combined and we should split them
		if string(i[0]) == "-" && string(i[1]) != "-" {
			flags := strings.Split(i[1:], "")

			for _, j := range flags {
				switch j {
				case "a":
					opts.includeHidden = true
				case "l":
					opts.longListing = true
				}
			}
		}
	}

	return opts
}
