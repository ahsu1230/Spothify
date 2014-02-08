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
		log.Fatal("Incorrect File Name")
	}
	
	currentPath, err := os.Getwd()
	if err!=nil {
		log.Fatal("Error getting working directory", err.Error())
	}
	s:= []string{currentPath, songFolder, *songName}
	songPath := strings.Join(s,"\\")
	
	fmt.Printf("Playing song... %s\n", songPath)
	
	cmd := exec.Command("wmplayer.exe", "/open", songPath)
	
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