package main

import (
	"fmt"
	"sync"

	parser "github.com/natsukagami/go-osu-parser"
)

const concurrentExtractors = 8

func extractor(input <-chan string, success chan<- BeatmapFile, fail chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range input {
		log("%s being parsed\n", file)
		beatmap, err := parser.ParseFile(file)
		if err != nil {
			fail <- fmt.Errorf("Cannot parse %s: %v", file, err)
		} else {
			success <- BeatmapFile{
				&beatmap,
				file,
				"",
			}
		}
	}
}

func init() {
	wg := sync.WaitGroup{}
	wg.Add(concurrentExtractors)
	for i := 0; i < concurrentExtractors; i++ {
		go extractor(OsuFiles, Beatmaps, Failed, &wg)
	}
	go func() {
		wg.Wait()
		close(Beatmaps)
		allDone.Done()
	}()
}
