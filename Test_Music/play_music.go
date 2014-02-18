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
	
	// Open file
	fo, err1 := os.Open(songPath)
	if err1 != nil {
		log.Fatal("Error opening file", err1.Error())
	}
	fmt.Println("Opened File...")
	
	// Read song file from songPath into some array...
	mp3Info, err2 := fo.Stat()
	fmt.Println("Obtained File Stats!")
	if err2 != nil {
		log.Fatal("Error getting FileInfo", err2.Error())
	}
	songArray := make([]byte, mp3Info.Size())
	m, errRead := fo.Read(songArray)
	if errRead != nil {
		log.Fatal("Error reading file!", errRead.Error())
	}
	fmt.Printf("%d bytes read!", m)
	
	/* // using File.Sys() - doesn't work... not too sure what .Sys() does
	songSys := mp3Info.Sys()
	if songSys == nil {
		log.Fatal("Song Sys is nil!")
	}
	fmt.Println("Converting Song to Byte Array...")
	songArray :=songSys.([]byte)
	//*/
	
	fmt.Println("Have Song Array!")
	
	// Does TMP directory already exist?
	_, err = os.Stat("TMP")
	if os.IsNotExist(err) {
		// Does not exist, create directory
		os.Mkdir("TMP",  os.ModeDir)
	}
	
	// Write song file into TMP Directory
	s = []string{currentPath, "TMP", *songName}
	NowPlayPath := strings.Join(s,"\\")
	fmt.Println("New Path:", NowPlayPath)
	
	// Write contents to TMP folder (NowPlayPath)
	newFO, err3 := os.Create(NowPlayPath)
	if err3 != nil {
		log.Fatal("Error creating New File", err3.Error())
	}
	n, err4 := newFO.Write(songArray)
	if err4 != nil {
		log.Fatal("Error writing New File", err4.Error())
	}
	
	fmt.Printf("FileSizes: %d  Copied: %d\n", len(songArray), n)
	
	
	// Play song in TMP Directory
	//cmd := exec.Command("wmplayer.exe", "/open", songPath)
	cmd := exec.Command("wmplayer.exe", "/open", NowPlayPath)
	
	runerror := cmd.Run()
	if runerror != nil {
		log.Fatal(runerror)
	}
}

/*
Exec Command:
wmplayer.exe /open C:\SongA.mp3
wmplayer.exe /open C:\SongB.mp3
wmplayer.exe /play /close C:\SongC.mp3
*/