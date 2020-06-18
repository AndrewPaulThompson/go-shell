package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type lsOpts struct {
	arguments       []string
	includeHidden   bool
	longListing     bool
	humanReadable   bool
	sortByTime      bool
	sortBySize      bool
	sortByExtension bool
	doNotSort       bool
	reverseSort     bool
	sort            string
}

func ls(args []string) {
	opts := lsOpts{}
	err := opts.splitArgs("ls", args)
	if err != nil {
		return
	}

	files, err := opts.getFileList()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = opts.sortFiles(files)
	if err != nil {
		fmt.Println(err)
		return
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

func (opts lsOpts) sortFiles(files []os.FileInfo) error {
	// If we have a sort by value
	if opts.sort != "" {
		switch opts.sort {
		case "none":
			opts.doNotSort = true
		case "extension":
			opts.sortByExtension = true
		case "size":
			opts.sortBySize = true
		case "time":
			opts.sortByTime = true
		default:
			return errors.New("Invalid sort type " + opts.sort)
		}
	}

	switch true {
	case opts.doNotSort:
		return nil
	case opts.sortByTime:
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})
	case opts.sortBySize:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size() > files[j].Size()
		})
	case opts.sortByExtension:
		sort.Slice(files, func(i, j int) bool {
			first := filepath.Ext(files[i].Name())
			second := filepath.Ext(files[j].Name())

			if strings.Compare(first, second) < 0 {
				return true
			}
			return false
		})
	default:
		sort.Slice(files, func(i, j int) bool {
			if strings.Compare(strings.ToLower(files[i].Name()), strings.ToLower(files[j].Name())) < 0 {
				return true
			}
			return false
		})
	}

	// Reverse the results if needed
	if opts.reverseSort != false {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	return nil
}

func (opts lsOpts) getFileList() ([]os.FileInfo, error) {
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

func (opts *lsOpts) splitArgs(command string, args []string) error {
	// Create new flagset
	fs := flag.NewFlagSet(command, flag.ContinueOnError)

	// Register flags
	fs.BoolVar(&opts.longListing, "l", false, "Long listing")
	fs.BoolVar(&opts.includeHidden, "a", false, "Show all files")
	fs.BoolVar(&opts.includeHidden, "all", false, "Show all files")
	fs.BoolVar(&opts.humanReadable, "h", false, "Human readable size")
	fs.BoolVar(&opts.humanReadable, "human-readable", false, "Human readable size")
	fs.BoolVar(&opts.sortByTime, "t", false, "Sort files by time")
	fs.BoolVar(&opts.sortBySize, "S", false, "Sort files by file size")
	fs.BoolVar(&opts.sortByExtension, "X", false, "Sort alphabetically by entry extension")
	fs.BoolVar(&opts.doNotSort, "U", false, "Do not sort; list entries in directory order")
	fs.BoolVar(&opts.doNotSort, "f", false, "Do not sort; list entries in directory order")
	fs.BoolVar(&opts.reverseSort, "r", false, "Reverse order while sorting")
	fs.BoolVar(&opts.reverseSort, "reverse", false, "Reverse order while sorting")
	fs.StringVar(&opts.sort, "sort", "", "Sort by WORD instead of name: none -U, extension -X, size -S, time -t")

	//  Parse arguments into flags
	err := fs.Parse(args)
	opts.arguments = fs.Args()

	return err
}
