package main

import (
	"fmt"

	parser "github.com/natsukagami/go-osu-parser"
)

type beatmapError struct {
	B BeatmapFile
	E error
}

func (b beatmapError) Error() string {
	return fmt.Sprintf("Error when processing %s: %s", b.B.OsuPath, b.E.Error())
}

// BeatmapFile represents a parsed beatmap file
type BeatmapFile struct {
	*parser.Beatmap
	OsuPath string
	Mp3Path string
}
