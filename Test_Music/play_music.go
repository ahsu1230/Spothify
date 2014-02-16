package main

import (	"fmt"
		"flag"
		"log"
		"os"
		"strings"
		"os/exec"
)

// To Run:
// go run play_music.go -song="SongA.mp3"

var songName *string = flag.String("song", "", "Song File Name")
var songFolder = "songfiles"
var MediaPlayer = "wmplayer.exe" 		// Application executable of Media Player

func main() {
	fmt.Println("Hello, World.\n")
	flag.Parse()
	if (*songName == "") {
		log.Fatal("Need a File Name!")
	}
	
	currentPath, err := os.Getwd()
	if err!=nil {
		log.Fatal("Error getting working directory", err.Error())
	}
	s:= []string{currentPath, songFolder, *songName}
	songPath := strings.Join(s,"\\")
	
	fmt.Printf("Playing song... %s\n", songPath)
	
	// Does TMP directory already exist?
	_, err = os.Stat("TMP")
	if os.IsNotExist(err) {
		// Does not exist, create directory
		os.Mkdir("TMP",  os.ModeDir)
	}
	
	// Read song file from songPath into some array...
	fo, err1 := os.Open(songPath)
	if err1 != nil {
		log.Fatal("Error opening file", err1.Error())
	}
	
	// Write song file into TMP Directory
	s = []string{currentPath, "TMP", *songName}
	NowPlayPath:= strings.Join(s,"\\")
	//
	
	// Play song in TMP Directory
	cmd := exec.Command("wmplayer.exe", "/open", songPath)
	//cmd := exec.Command("wmplayer.exe", "/open", NowPlayPath)
	
	runerror := cmd.Run()
	if runerror != nil {
		log.Fatal(runerror)
	}
}

/*
Command:
wmplayer.exe /open C:\SongA.mp3
wmplayer.exe /open C:\SongB.mp3
wmplayer.exe /play /close C:\SongC.mp3
*/