package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type opts struct {
	arguments     []string
	includeHidden bool
	longListing   bool
	humanReadable bool
	sortByTime    bool
	sortBySize    bool
	doNotSort     bool
}

func ls(args []string) {
	opts, err := splitArgs("ls", args)
	if err != nil {
		return
	}

	files, err := getFileList(opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	sortFiles(files, opts)

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

func sortFiles(files []os.FileInfo, opts opts) {
	switch true {
	case opts.doNotSort:
		return
	case opts.sortByTime:
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})
	case opts.sortBySize:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size() > files[j].Size()
		})
	default:
		sort.Slice(files, func(i, j int) bool {
			if strings.Compare(strings.ToLower(files[i].Name()), strings.ToLower(files[j].Name())) < 0 {
				return true
			}
			return false
		})
	}
}

func getFileList(opts opts) ([]os.FileInfo, error) {
	// If we have some arguments, we want to list those instead of the working directory
	var files []os.FileInfo
	if len(opts.arguments) > 0 {
		for _, arg := range opts.arguments {
			fileList, err := ioutil.ReadDir(arg)

			// If we error here, we've tried to read a file, not a directory
			if err != nil {
				// So try to read the arg as a file instead
				file, err := os.Stat(arg)
				if err != nil {
					// If we error here, we can't read the argument
					// Return early to stop looping
					return files, err
				}

				// Append this single file to the file slice
				files = append(files, file)
			}

			// Append the files from ReadDir to the file slice
			files = append(files, fileList...)
		}
	} else {
		// Get the current directory
		wd, err := os.Getwd()
		if err != nil {
			return files, err
		}

		// Get files in the working directory
		files, err = ioutil.ReadDir(wd)
		if err != nil {
			return files, err
		}
	}

	return files, nil
}

func splitArgs(command string, args []string) (opts, error) {
	// Create new opts
	opts := opts{}

	// Create new flagset
	fs := flag.NewFlagSet(command, flag.ContinueOnError)

	// Register flags
	fs.BoolVar(&opts.longListing, "l", false, "Long listing")
	fs.BoolVar(&opts.includeHidden, "a", false, "Show all files")
	fs.BoolVar(&opts.humanReadable, "h", false, "Human readable size")
	fs.BoolVar(&opts.sortByTime, "t", false, "Sort files by time")
	fs.BoolVar(&opts.sortBySize, "S", false, "Sort files by file size")
	fs.BoolVar(&opts.doNotSort, "U", false, "Do not sort; list entries in directory order")
	fs.BoolVar(&opts.doNotSort, "f", false, "Do not sort; list entries in directory order")

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
