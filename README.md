# osu-music-copy
[![Build Status](https://travis-ci.org/natsukagami/osu-music-copy.svg?branch=master)](https://travis-ci.org/natsukagami/osu-music-copy)

A program that allows you to copy and organize your osu! music library with ease.

## Running
### From pre-compiled Executable
One can easily download a binary executable on the [Release](https://github.com/natsukagami/osu-music-copy/releases) section.
Once downloaded, the program can be run directly from a command prompt.

### Compiling from Source
To compile the program you will need the Go compiler. It can be obtained from the [Golang website](https://golang.org). 
The recommended version is 1.7.x or later.

After installing Go compiler, run the following command from the command prompt:
```
go get github.com/natsukagami/osu-music-copy
go build github.com/natsukagami/osu-music-copy
```

After that you should see an executable named `osu-music-copy`. 

## Command-line Arguments
```
Copies and fills metadata for osu! songs, organized in a beautiful way.
Usage: ./osu-music-copy [options] path-to-osu-songs-folder
Available options:
	-h, -help 
    	Show this help
  -ntfs
    	Fixes the name to fit NTFS filename limits
  -o string
    	The destination for output songs (default "songs/")
  -s	Skips songs already found
  -v	Print logging information
```