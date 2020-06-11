package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type opts struct {
	arguments     []string
	includeHidden bool
	longListing   bool
	humanReadable bool
}

func ls(args []string) {
	opts, err := handleArgs("ls", args)
	if err != nil {
		return
	}

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

		var size string
		if opts.humanReadable {
			size = prettifySize(file.Size())
		} else {
			size = strconv.FormatInt(file.Size(), 10)
		}

		// Use long listing format
		if opts.longListing {
			fmt.Printf("%s %s %s %s\n", file.Mode(), size, file.ModTime().Format("Jan 2 15:04"), file.Name())
			continue
		}

		// Default to just printing the name
		fmt.Println(file.Name())
	}
}

func handleArgs(command string, args []string) (opts, error) {
	// Create new opts
	opts := opts{}

	// Create new flagset
	fs := flag.NewFlagSet(command, flag.ContinueOnError)

	// Register flags
	fs.BoolVar(&opts.longListing, "l", false, "Long listing")
	fs.BoolVar(&opts.includeHidden, "a", false, "Show all files")
	fs.BoolVar(&opts.humanReadable, "h", false, "Human readable size")

	//  Parse arguments into flags
	err := fs.Parse(args)
	opts.arguments = fs.Args()

	return opts, err
}

func prettifySize(b int64) string {
	// Unit size
	const unit = 1024

	// If the input is less than 1KB, show bytes suffix
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	// Convert unit size to int64, initialize counter variable
	div, exp := int64(unit), 0

	// Keep dividing the input value by the unit size until it equals less than the unit size
	for n := b / unit; n >= unit; n /= unit {
		// Keep track of the total units we've divided by
		div *= unit

		// Increment the counter
		exp++
	}

	// Divide the input value by the final count of units
	// Use the counter for number of times divided to get the unit suffix
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
