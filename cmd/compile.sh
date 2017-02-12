#!/bin/bash
GOOS=windows GOARCH=amd64 go build -o osu-music-copy_win64.exe ..
GOOS=windows GOARCH=386 go build -o osu-music-copy_win32.exe ..
GOOS=linux GOARCH=amd64 go build -o osu-music-copy_linux64 ..
GOOS=linux GOARCH=386 go build -o osu-music-copy_linux32 ..
GOOS=darwin GOARCH=amd64 go build -o osu-music-copy_darwin64 ..
