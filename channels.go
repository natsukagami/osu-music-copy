package main

// The following channels are shared
var (
	OsuFiles = make(chan string)
	Beatmaps = make(chan BeatmapFile)
	Tagged   = make(chan BeatmapFile)
	Failed   = make(chan error)
)
