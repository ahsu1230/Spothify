package peernode

import (
	"fmt"
	"log"
	"net/rpc"
	"strings"
	"time"
	"os"
	"os/exec"
	"../proto"
	"../../consts"
	"../../util/stringlist"
	"../../util/songinfo"
)

var MediaPlayer = "wmplayer.exe" 		// Application executable of Media Player

type PeerClient struct {
	username 		string
	portnum			int
	serverHP			string
	serverConn 		*rpc.Client
	offlineMode 		bool
	
	playlistsMap		map[string] stringlist.List		// playlistName -> list of songs
	songsMap		map[string] string			// songName -> path of song
}

func NewPeerClient(username string, portnum int, serverHP string, offline bool) *PeerClient {
	pc := new(PeerClient)
	pc.username = username
	pc.portnum = portnum
	pc.serverHP = serverHP
	pc.offlineMode = offline
	
	pc.playlistsMap = make(map[string] stringlist.List)
	pc.songsMap = make(map[string] string)
	
	fmt.Printf("New PeerClient '%s' Created ", pc.username);
	if pc.offlineMode {
		fmt.Printf("Offline!\n");
	} else {
		fmt.Printf("Online!\n");
	}
	return pc
}


func (pc *PeerClient) ConnectToServer(hostport string) {
	// Attempt to connect to peer server until successful
	count := 0
	for count < constants.RpcTries {
		conn, err := rpc.DialHTTP("tcp", hostport)
		if err != nil {
			log.Fatal("Could not connect to server:", err)
		} else {
			pc.serverConn = conn
			break
		}
		time.Sleep(time.Duration(constants.RpcWaitMillis) * time.Millisecond)
		count++
	}
	if count == constants.RpcTries {
		log.Fatal("Could not connect to server... too many tries")
		return
	}
	log.Printf("Peer Client [%s] connected to server [%s]\n", pc.username, pc.serverHP)
}


func (pc *PeerClient) DisconnectFromServer() error {
	if pc.offlineMode {
		return nil
	}
	err := pc.serverConn.Close()
	if err!=nil {
		return err
	}
	log.Printf("Player [%s] disconnected from server [%s]\n", pc.username, pc.serverHP)
	pc.serverHP = ""
	return nil
}


/* Call RPCs */
func (pc *PeerClient) PlaySong (songname string) (int, error) {
	fmt.Printf("Playing Song '%s'...\n", songname)

	// First check to see if song is in offline local cache...
	// Play Song!
	
	// Otherwise... (file does not exist on local cache)
	fmt.Println("Not in local! Look online!")
	if pc.offlineMode { // if offline, return with error
		fmt.Println("Offline - song not found!\n")
		return 0, nil
	} else { // online, send request to server
		newPInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
		newSInfo := songinfo.NewSong(songname)
		args := peernodeproto.PlayArgs{newPInfo, *newSInfo}
		var reply peernodeproto.PlayReply
		err := pc.serverConn.Call("PeerRPC.PlaySongRequest", &args, &reply)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return 0, err
		} else if reply.Status != peernodeproto.OK {
			fmt.Println("Denied!")
			return reply.Status, nil
		}
		fmt.Println("Request OK'd")
		
		// Play Song!
		currentPath, errCD := os.Getwd()
		if errCD!=nil {
			log.Fatal("Error getting working directory", errCD.Error())
		}
		s1 := strings.Split(currentPath, "\\")
		s2 := s1[0:len(s1)-1]
		s3 := strings.Join(s2, "\\")
		playPath := fmt.Sprintf("%s\\userdata\\%s\\TMP\\%s", s3, pc.username, songname)
		fmt.Println(playPath)
		//reply.PlayPath = "C:\\Users\\AaronHsu\\Documents\\GitHub\\Spothify\\src\\peernode\\userdata\\ahsu\\TMP\\SongA.mp3"
		//reply.PlayPath = "C:\\Users\\AaronHsu\\Documents\\GitHub\\Spothify\\src\\peernode\\userdata\\jfan89\\TMP\\SongB.mp3"
		//fmt.Println(reply.PlayPath)
		//fmt.Println(playPath == reply.PlayPath)
		// Does the file exist?
		_, errOS := os.Stat(playPath)
		if !os.IsNotExist(errOS) { // If it does exist, play it!
			cmd := exec.Command("wmplayer.exe", "/open", playPath)
			runerror := cmd.Run()
			if runerror != nil {
				log.Println("Error playing music file!", runerror.Error())
				return 0, runerror
			}
		} else {
			log.Println("File does not exist? ")
		}
	}

	return 0, nil
}


func (pc *PeerClient) AddPlaylist (newPL string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Adding New Playlist '%s'\n", newPL)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.ChangePLArgs{newPeerInfo, newPL, ""}
	var reply peernodeproto.ChangePLReply
	err := pc.serverConn.Call("PeerRPC.AddPLRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Request OK'd")
	return 0, nil
}


func (pc *PeerClient) DeletePlaylist (pl string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Deleting Existing Playlist '%s'\n", pl)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.ChangePLArgs{newPeerInfo, pl, ""}
	var reply peernodeproto.ChangePLReply
	err := pc.serverConn.Call("PeerRPC.DeletePLRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Request OK'd")
	return 0, nil
}


func (pc *PeerClient) RenamePlaylist (oldPL, newPL string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Renaming Playlist '%s' to '%s'\n", oldPL, newPL)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.ChangePLArgs{newPeerInfo, oldPL, newPL}
	var reply peernodeproto.ChangePLReply
	err := pc.serverConn.Call("PeerRPC.RenamePLRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Request OK'd")
	return 0, nil
}


func (pc *PeerClient) SortPlaylist (pl string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Sorting Playlist '%s'\n", pl)
	newMsg := fmt.Sprintf("Sorting Playlist:[%s]", pl)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SendMsgArgs{newPeerInfo, peernodeproto.TO_BRIDGE, newMsg}
	var reply peernodeproto.SendMsgReply
	err := pc.serverConn.Call("PeerRPC.SendMsgRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Message Replied: ", reply.Message)
	return 0, nil
}


func (pc *PeerClient) DownloadPlaylist (pl string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Downloading Playlist '%s'\n", pl)
	newMsg := fmt.Sprintf("Downloading Playlist:[%s]", pl)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SendMsgArgs{newPeerInfo, peernodeproto.TO_BRIDGE, newMsg}
	var reply peernodeproto.SendMsgReply
	err := pc.serverConn.Call("PeerRPC.SendMsgRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Message Replied: ", reply.Message)
	return 0, nil
}


func (pc *PeerClient) ViewAllPlaylists () (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Viewing All Playlists\n")
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.ChangePLArgs{newPeerInfo, "", ""}
	var reply peernodeproto.ChangePLReply
	err := pc.serverConn.Call("PeerRPC.ViewAllPLRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("All Playlists: {", reply.DisplayStr, "}")
	return 0, nil
}


func (pc *PeerClient) SearchSong (songName string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Searching Song '%s'\n", songName)
	newMsg := fmt.Sprintf("Search for Song:[%s]", songName)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SendMsgArgs{newPeerInfo, peernodeproto.TO_BRIDGE, newMsg}
	var reply peernodeproto.SendMsgReply
	err := pc.serverConn.Call("PeerRPC.SendMsgRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Message Replied: ", reply.Message)
	return 0, nil
}


func (pc *PeerClient) SearchArtist (artistName string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Searching Artist '%s'\n", artistName)
	newMsg := fmt.Sprintf("Search for Artist:[%s]", artistName)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SendMsgArgs{newPeerInfo, peernodeproto.TO_BRIDGE, newMsg}
	var reply peernodeproto.SendMsgReply
	err := pc.serverConn.Call("PeerRPC.SendMsgRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Message Replied: ", reply.Message)
	return 0, nil
}


func (pc *PeerClient) AddSong (songName, pl string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Adding Song '%s' to Playlist '%s'\n", songName, pl)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SongArgs{newPeerInfo, songName, pl}
	var reply peernodeproto.SongReply
	err := pc.serverConn.Call("PeerRPC.AddSongRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Request Ok'd")
	return 0, nil
}


func (pc *PeerClient) DeleteSong (songName, pl string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Deleting Song '%s' from Playlist '%s'\n", songName, pl)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SongArgs{newPeerInfo, songName, pl}
	var reply peernodeproto.SongReply
	err := pc.serverConn.Call("PeerRPC.DeleteSongRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Request Ok'd")
	return 0, nil
}


func (pc *PeerClient) ViewAllSongs (pl string) (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("View All Songs in Playlist '%s'\n", pl)
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.SongArgs{newPeerInfo, "", pl}
	var reply peernodeproto.SongReply
	err := pc.serverConn.Call("PeerRPC.ViewAllSongsRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	fmt.Println("Songs in Playlist: ", reply.DisplayStr)
	return 0, nil
}


func (pc *PeerClient) Quit () (int, error) {
	if pc.offlineMode {
		return peernodeproto.OFFLINEERROR, nil
	}

	fmt.Printf("Requesting Disconnect\n")
	newPeerInfo := peernodeproto.PeerInfo{pc.username, pc.portnum}
	args := peernodeproto.QuitArgs{newPeerInfo}
	var reply peernodeproto.QuitReply
	err := pc.serverConn.Call("PeerRPC.QuitRequest", &args, &reply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if reply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return reply.Status, nil
	}
	
	pc.DisconnectFromServer()
	fmt.Printf("Quit! %s logged out!\n", pc.username)
	return 0, nil
}


/*
// Functions include: Add/Remove/Rename Playlist, Add/Remove Song, Download Playlist, Search Song, Play Song, Quit
func (pc *PeerClient) CallServerRPC(requestType string, args interface{}, reply interface{}) (int, error) {
	RPCname := ""

	var RPCargs interface{}	// don't really work because can't access status
	var RPCreply interface{}	// in which case it doesn't save much space from using totally different function
	switch(requestType) {
	case "Message":
		RPCname = "PeerRPC.SendMsgRequest"
		RPCargs := args.(peernodeproto.SendMsgArgs)
		RPCreply := args.(peernodeproto.SendMsgReply)
		break
	case "AddPlaylist":
		break
	case "RemovePlaylist":
		break
	case "RenamePlaylist":
		break
	case "DownloadPlaylist":
		break
	case "AddSong":
		break
	case "RemoveSong":
		break
	case "SearchSong":
		break
	case "PlaySong":
		break
	case "Quit":
		break
	default: // Unsupported Request
		fmt.Println("Unsupported Request")
		return -1, nil
	}
	
	err := pc.serverConn.Call(RPCname, &RPCargs, &RPCreply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return 0, err
	} else if RPCreply.Status != peernodeproto.OK {
		fmt.Println("Denied!")
		return RPCreply.Status, nil
	}
	return 0, nil
}
//*/