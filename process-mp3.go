package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/bogem/id3v2"
)

const (
	concurrentCopiers = 8
	concurrentTaggers = 20
)

var ntfsInvalidRegex = regexp.MustCompile(`[\<\>\:\"\/\\\|\?\*]`)

// A map with concurrent processing
type mutexMap struct {
	m map[string]bool
	l sync.Mutex
}

func (m *mutexMap) Add(f string) bool {
	m.l.Lock()
	defer m.l.Unlock()
	_, ok := m.m[f]
	if !ok {
		m.m[f] = true
	}
	return !ok
}

var mp = mutexMap{m: make(map[string]bool)}

func doCopy(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	err = os.MkdirAll(filepath.Dir(dst), 0755|os.ModeDir)
	if err != nil {
		return
	}
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return
}

// Copies the mp3 file to another folder
func copyMp3(input <-chan BeatmapFile, success chan<- BeatmapFile, fail chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for b := range input {
		Artist := strings.Replace(b.Artist, "/", "_", -1)
		Title := strings.Replace(b.Title, "/", "_", -1)
		Source := strings.Replace(b.Source, "/", "_", -1)
		var (
			oldPath = filepath.Join(filepath.Dir(b.OsuPath), b.AudioFilename)
			newPath = filepath.Join(outputFolder, Artist, Title+".mp3")
		)
		if b.Source != "" {
			newPath = filepath.Join(outputFolder, Artist, Source+" - "+Title+".mp3")
		}
		if ntfsFix {
			f := func(s string) string { return ntfsInvalidRegex.ReplaceAllString(s, "_") }
			if b.Source != "" {
				newPath = filepath.Join(outputFolder, f(b.Artist), f(b.Source+" - "+b.Title)+".mp3")
			} else {
				newPath = filepath.Join(outputFolder, f(b.Artist), f(b.Title)+".mp3")
			}
		}
		var err error
		if newPath, err = filepath.Abs(newPath); err != nil {
			fail <- beatmapError{b, err}
			continue
		}
		log("Copying %s to %s...\n", oldPath, newPath)
		if !mp.Add(newPath) {
			// Path already exists
			continue
		}
		if _, err := os.Stat(newPath); err != nil && skipFound {
			// Skips file that exists
			continue
		}
		if err := doCopy(oldPath, newPath); err != nil {
			fmt.Println(err)
			fail <- beatmapError{b, err}
		} else {
			b.Mp3Path = newPath
			success <- b
		}
	}
}

var copyDone = make(chan BeatmapFile)

func init() {
	wg := sync.WaitGroup{}
	wg.Add(concurrentCopiers)
	for i := 0; i < concurrentCopiers; i++ {
		go copyMp3(Beatmaps, copyDone, Failed, &wg)
	}
	go func() {
		wg.Wait()
		close(copyDone)
		allDone.Done()
	}()
}

func openMp3(p string) (t *id3v2.Tag, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	t, err = id3v2.Open(p)
	return
}

func cleanMp3(p string) error {
	mp3, err := openMp3(p)
	if err != nil {
		return err
	}

	defer mp3.Close()
	mp3.DeleteAllFrames()
	if err := mp3.Save(); err != nil {
		return err
	}
	return nil
}

// Tags the mp3's metadata
func tagMp3(input <-chan BeatmapFile, success chan<- BeatmapFile, fail chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for b := range input {
		if err := cleanMp3(b.Mp3Path); err != nil {
			fail <- beatmapError{b, err}
			continue
		}
		log("Tagging %s\n", b.OsuPath)
		mp3, err := openMp3(b.Mp3Path)
		if err != nil {
			fail <- beatmapError{b, err}
			continue
		}
		mp3.SetArtist(b.Artist)
		mp3.SetTitle(b.Title)
		mp3.SetGenre("osu!")
		if b.Source != "" {
			mp3.AddFrame("TPE2", id3v2.TextFrame{Encoding: id3v2.ENUTF8, Text: b.Artist})
			mp3.SetAlbum(b.Source)
		} else {
			mp3.AddFrame("TPE2", id3v2.TextFrame{Encoding: id3v2.ENUTF8, Text: "Various Artists"})
			mp3.SetAlbum("osu!")
		}
		if b.BgFilename != "" {
			artwork, err := ioutil.ReadFile(filepath.Join(filepath.Dir(b.OsuPath), b.BgFilename))
			if err == nil {
				picType := "image/png"
				if filepath.Ext(b.BgFilename) == ".jpg" {
					picType = "image/jpeg"
				}
				fmt.Println("Found ", b.BgFilename, picType)
				mp3.AddAttachedPicture(id3v2.PictureFrame{
					Encoding:    id3v2.ENUTF8,
					MimeType:    picType,
					PictureType: id3v2.PTFrontCover,
					Description: "Front cover",
					Picture:     artwork,
				})
			}
		}
		if err := mp3.Save(); err != nil {
			fail <- beatmapError{b, err}
			continue
		}
		mp3.Close()
		success <- b
	}
}

func init() {
	wg := sync.WaitGroup{}
	wg.Add(concurrentTaggers)
	for i := 0; i < concurrentTaggers; i++ {
		go tagMp3(copyDone, Tagged, Failed, &wg)
	}
	go func() {
		wg.Wait()
		close(Tagged)
		allDone.Done()
	}()
}
