package main

import (
	"log"
	"fmt"
	"flag"
	"os"
	//"time"
	"../impl"
	"../proto"
)


// For parsing the command line
type cmd_info struct {
	cmdline string
	funcname string
	nargs int // number of required args
}
const (
	CMD_PUT = iota
	CMD_GET
)

var portnum *int = flag.Int("port", 0, "server port # to connect to")
var offline *bool = flag.Bool("offline", false, "run in offline mode?")

const DEFAULT_PORTNUM = 9009

func main() {

	fmt.Printf("starting...\n");
	
	flag.Parse()
	if (flag.NArg() < 2) {
		fmt.Printf("arguments: [username] -port=#? -offline [command] [arg1][arg2][arg3]")
		log.Fatal("Insufficient arguments to client")
	}
	if (*portnum == 0) {
		*portnum = 9009
		log.Printf("DEFAULTING PORTNUM to %d...\n", DEFAULT_PORTNUM)
	}
	fmt.Println("num flags", flag.NArg())
	
	username := flag.Arg(0)
	fmt.Printf("Attempt UserNode sign-in [%s]\n", username)

	serverPort := fmt.Sprintf("%d", *portnum)
	serverAddress := fmt.Sprintf("localhost:%d", *portnum)
	fmt.Printf("Server address: '%s'\n", serverAddress)
	fmt.Printf("Server portnum: '%s'\n", serverPort)
	pc := peernode.NewPeerClient(username, *portnum, serverAddress, *offline)
	
	cmd := flag.Arg(1)
	
	// Command List
	cmdlist := []cmd_info {
		{ "ap", "ServerRPC.AddPlaylist", 2 },		// playlistname
		{ "dp", "ServerRPC.DeletePlaylist", 2 },		// playlistname
		{ "rp", "ServerRPC.RenamePlaylist", 3 },		// playlistname, newplaylistname
		{ "vp", "ServerRPC.ViewAllPlaylists", 1 },		// (nothing)
		{ "sp", "ServerRPC.SortPlaylist", 2 },		// playlistname (by name for now)
		{ "dlp", "ServerRPC.DownloadPlaylist", 2 },	// playlistname
		
		{ "ss", "ServerRPC.SearchSong", 2 },		// songname
		{ "sa", "ServerRPC.SearchArtist", 2 },		// artistname
		
		{ "ps", "ServerRPC.PlaySong", 2 },			// songname
		{ "as", "ServerRPC.AddSong", 3 },			// playlistname, songname
		{ "ds", "ServerRPC.DeleteSong", 3 },		// playlistname, songname
		{ "vs", "ServerRPC.ViewAllSongs", 2 },		// playlistname
		{ "q", "RegionRPC.Quit", 1},
	}
	cmdmap := make(map[string]cmd_info)
	for _, j := range(cmdlist) {
			cmdmap[j.cmdline] = j
	}

	ci, found := cmdmap[cmd]
	if (!found) {
		log.Fatal("Unknown command ", cmd)
	}
	if (flag.NArg() < (ci.nargs+1)) {
		log.Fatal("Insufficient arguments for ", cmd)
	}
	
	if !(*offline) {
		pc.ConnectToServer(serverAddress)
	}
	switch(cmd) {
		case "ap":  // Add Playlist
			status, err := pc.AddPlaylist(flag.Arg(2))
			PrintStatus("AddPlaylist", status, err)
			break;
		case "dp":  // Delete Playlist
			status, err := pc.DeletePlaylist(flag.Arg(2))
			PrintStatus("DeletePlaylist", status, err)
			break;
		case "rp":  // Rename Playlist
			status, err := pc.RenamePlaylist(flag.Arg(2), flag.Arg(3))
			PrintStatus("RenamePlaylist", status, err)
			break;
		case "sp":  // Sort Playlist
			status, err := pc.SortPlaylist(flag.Arg(2))
			PrintStatus("SortPlaylist", status, err)
			break;
		case "dlp":  // Download Playlist
			status, err := pc.DownloadPlaylist(flag.Arg(2))
			PrintStatus("DownloadPlaylist", status, err)
			break;
		case "vp":  // View Playlists
			status, err := pc.ViewAllPlaylists()
			PrintStatus("ViewPlaylist", status, err)
			break;
		case "ss":  // Search Song
			status, err := pc.SearchSong(flag.Arg(2))
			PrintStatus("SearchSong", status, err)
			break;
		case "sa":  // Search Artist
			status, err := pc.SearchArtist(flag.Arg(2))
			PrintStatus("SearchArtist", status, err)
			break;
		case "ps":  // Play Song
			status, err := pc.PlaySong(flag.Arg(2))
			PrintStatus("PlaySong", status, err)
			break;
		case "as":  // Add Song
			status, err := pc.AddSong(flag.Arg(2), flag.Arg(3))
			PrintStatus("AddSong", status, err)
			break;
		case "ds":  // Delete Song
			status, err := pc.DeleteSong(flag.Arg(2), flag.Arg(3))
			PrintStatus("DeleteSong", status, err)
			break;
		case "vs":  // View Songs
			status, err := pc.ViewAllSongs(flag.Arg(2))
			PrintStatus("ViewSongs", status, err)
			break;
		case "q":  // Quit
			status, err := pc.Quit()
			PrintStatus("Quit", status, err)
			os.Exit(0);
			break;
		default: // Unknown Command
			fmt.Println("Unknown Command\n")
	}
	pc.DisconnectFromServer()
}

func PrintStatus(cmdname string, status int, err error) {
	if err!=nil {
		fmt.Printf("%s\n", err.Error())
	} else if status == peernodeproto.OFFLINEERROR {
		fmt.Printf("%s denied - offline mode\n", cmdname)
	} else if status != peernodeproto.OK {
		fmt.Printf("%s failed... %d\n", cmdname, status)
	} else {
		fmt.Printf("%s success!\n", cmdname)
	}
	/*
	if (status == nodeproto.OK) {
		fmt.Printf("%s succeeded\n", cmdname)
	} else {
		fmt.Printf("%s failed: %s\n", cmdname, PlayerStatusToString(status))
	}
	*/
}
