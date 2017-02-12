package main

import (
	"io/ioutil"
	"path/filepath"
	"sync"
)

var dirChan = make(chan string)

const concurrentDirProcessors = 2

// processDir scans a folder and sends the first .osu file found.
func processDir(input <-chan string, output chan<- string, fail chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for dir := range input {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fail <- err
			continue
		}
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".osu" {
				p := filepath.Join(dir, file.Name())
				log("%s found, entering queue\n", p)
				output <- p
				break
			}
		}
	}
}

func init() {
	wg := sync.WaitGroup{}
	wg.Add(concurrentDirProcessors)
	for i := 0; i < concurrentDirProcessors; i++ {
		go processDir(dirChan, OsuFiles, Failed, &wg)
	}
	go func() {
		wg.Wait()
		close(OsuFiles)
		allDone.Done()
	}()
}

func processInputFolder() error {
	files, err := ioutil.ReadDir(inputFolder)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go func(f string) {
				dirChan <- filepath.Join(inputFolder, f)
				wg.Done()
			}(file.Name())
		}
	}
	wg.Wait()
	close(dirChan)
	return nil
}
