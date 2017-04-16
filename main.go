package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
)

var (
	verbose      bool
	ntfsFix      bool
	skipFound    bool
	outputFolder string
	inputFolder  string
)

var allDone = sync.WaitGroup{}

func init() {
	allDone.Add(4)
}

// init flags parser
func init() {
	flag.BoolVar(&verbose, "v", false, "Print logging information")
	flag.BoolVar(&ntfsFix, "ntfs", false, "Fixes the name to fit NTFS filename limits")
	flag.BoolVar(&skipFound, "s", false, "Skips songs already found")
	flag.StringVar(&outputFolder, "o", "songs/", "The destination for output songs")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Copies and fills metadata for osu! songs, organized in a beautiful way.")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] path-to-osu-songs-folder\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Available options:\n")
		flag.PrintDefaults()
	}
}

// Parse program flags
func parseInit() error {
	flag.Parse()
	inputFolder = flag.Arg(0)
	if len(inputFolder) == 0 {
		return errors.New("Please specify the input folder")
	}
	if stats, err := os.Stat(inputFolder); err != nil {
		return err
	} else if !stats.IsDir() {
		return errors.New("Input folder is not a folder")
	}
	if stats, err := os.Stat(outputFolder); err != nil {
		if err := os.MkdirAll(outputFolder, 0755|os.ModeDir); err != nil {
			return err
		}
	} else if !stats.IsDir() {
		return errors.New("Output folder is not a folder")
	}
	return nil
}

func mustOk(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	mustOk(parseInit())
	go func() {
		for {
			select {
			case b := <-Tagged:
				fmt.Printf("Success: %s\n", b.Mp3Path)
			case e := <-Failed:
				if be, ok := e.(beatmapError); ok && be.B.Mp3Path != "" {
					os.Remove(be.B.Mp3Path)
				}
				fmt.Println("Failed: ", e)
			}
		}
	}()
	mustOk(processInputFolder())
	allDone.Wait()
}
