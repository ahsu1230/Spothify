package storagenode

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"../proto"
	"../rpc"
	"../../consts"
	"../../monitor/proto"
	"../../util/stringlist"
)

type SongInfo struct {
	name	string
	artist	string
}

type StorageServer struct {
	id uint32			// Storage ID
	portnum int			// Port number of Storage Server
	hostport string		// Hostport of Storage Server
	monitorHP string	// Hostport of Monitor Server
	srpc *storagenoderpc.StorageRPC	// RPC object used to register RPC functions
	connected map[string] *rpc.Client	// list of servers connected to
	
	// Storing Info
	songDB	map[SongInfo] string		// map song info to file path
	userDB	map[string] (map[string] *stringlist.List) // map usernames to a (map of playlist names -> list of songs)	
	userLocks	map[string] *sync.RWMutex
	
	DB_users	*stringlist.List
	DB_lock	*sync.RWMutex
}

func NewStorageServer(portnum int, monitorPort int, id uint32) *StorageServer {

	ss := new(StorageServer)
	ss.portnum = portnum
	ss.hostport = fmt.Sprintf("localhost:%d", portnum)
	ss.monitorHP = fmt.Sprintf("localhost:%d", monitorPort)
	ss.connected = make(map[string] *rpc.Client)
	ss.id = id
	
	ss.songDB = make(map[SongInfo] string)
	ss.userDB = make(map[string] (map[string] *stringlist.List)) // when new playlist, make interior map
	ss.userLocks = make(map[string] *sync.RWMutex)

	ss.DB_users = stringlist.NewList()	
	ss.DB_lock = new (sync.RWMutex)
	//ss.DB_lock.Unlock()
	
	
	//Register RPC functions
	ss.srpc = storagenoderpc.NewStorageRPC(ss)
	rpc.Register(ss.srpc)
	
	// Register With Monitor
	ss.RegisterWithMonitor()
	return ss
}


func (ss *StorageServer) RegisterWithMonitor() {
	fmt.Println("Registering with Monitor...")
	err, conn := constants.DialToServer(ss.monitorHP)
	if err != nil {
		log.Fatal("Problem Detected Dialing to Server", ss.monitorHP)
	}
	ss.connected[ss.monitorHP] = conn
	
	registerInfo := monitorproto.Node{monitorproto.STORAGE, ss.hostport, ss.id}
	newArgs := monitorproto.RegisterArgs{registerInfo}
	var newReply monitorproto.RegisterReply
	
	err = conn.Call("MonitorRPC.RegisterServer", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		log.Fatal(err.Error())
	} else if newReply.Status != monitorproto.OK {
		fmt.Println("Denied!")
		log.Fatal("Registration Denied!")
	}
	fmt.Println("Registered!")
}


func (ss *StorageServer) NewUsernameRequest(username string) bool {
	ss.DB_lock.Lock()
	exists := ss.DB_users.Contains(username)
	if !exists {
		fmt.Println("New User!", username)
		// make a new lock & new map for new user
		ss.DB_users.Insert(username)
		ss.userLocks[username] = new (sync.RWMutex)
		ss.userDB[username] = make(map[string] *stringlist.List)
		fmt.Println("Done adding new user")
	}
	ss.DB_lock.Unlock()
	return exists
}


func (ss *StorageServer) SendMsgRequest(args *storagenodeproto.SendMsgArgs, reply *storagenodeproto.SendMsgReply) error {
	fmt.Printf("Message: [%s:%s] from Bridge [%s] \n", args.CInfo.Username, args.Message, args.CInfo.BridgeHP)
	reply.Message = "Request Handled!"
	reply.Status = storagenodeproto.OK
	fmt.Println("\tSendRequest done\n")
	return nil
}


func (ss *StorageServer) AddPLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	fmt.Printf("AddPlaylist [%s] Request Received from '%s' \n", args.TargetPlaylistName, args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	// is it user's first request?
	// if so, create new map of playlists -> list of songs, create lock for user, add username to storage's list of users
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	
	// does playlist already exist?
	_, exists := ss.userDB[targetUser][args.TargetPlaylistName]
	if !exists { // if so, add playlist to user map with empty list of songs
		ss.userDB[targetUser][args.TargetPlaylistName] = stringlist.NewList()
		reply.Status = storagenodeproto.OK
	} else { // otherwise, return error status - playlist already exists!
		fmt.Println("\tError: Object Exists!")
		reply.Status = storagenodeproto.OBJECT_EXISTS
	}
	
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) DeletePLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	fmt.Printf("DeletePlaylist [%s] Request Received from '%s' \n", args.TargetPlaylistName, args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	_, exists := ss.userDB[targetUser][args.TargetPlaylistName]
	if !exists {
		fmt.Println("\tError: Object Not Found!")
		reply.Status = storagenodeproto.OBJECT_NOT_FOUND
	} else {
		delete( ss.userDB[targetUser], args.TargetPlaylistName)
		reply.Status = storagenodeproto.OK
	}
	
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) SortPLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	fmt.Printf("SortPlaylist [%s] Request Received from '%s' \n", args.TargetPlaylistName, args.CInfo.Username)
	return nil
}


func (ss *StorageServer) RenamePLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	fmt.Printf("RenamePlaylist [%s:%s] Request Received from '%s' \n", args.TargetPlaylistName, args.NewPlaylistName, args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	_, exists := ss.userDB[targetUser][args.TargetPlaylistName]
	if !exists {
		fmt.Println("\tError: Object Not Found!")
		reply.Status = storagenodeproto.OBJECT_NOT_FOUND
	} else {
		ss.userDB[targetUser][args.NewPlaylistName] = ss.userDB[targetUser][args.TargetPlaylistName]
		delete( ss.userDB[targetUser], args.TargetPlaylistName)
		reply.Status = storagenodeproto.OK
	}
	
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) DownloadPLRequest(args *storagenodeproto.DownloadPLArgs, reply *storagenodeproto.DownloadPLReply) error {
	fmt.Printf("DownloadPlaylist [%s] Request Received from '%s' \n", args.TargetPlaylistName, args.CInfo.Username)
	return nil
}


func (ss *StorageServer) ViewAllPLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	fmt.Printf("ViewAllPlaylist Request Received from '%s' \n", args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	display := ""
	for pl, _ := range ss.userDB[targetUser] {
		display += pl
		display += ", "
	}
	reply.DisplayStr = display
	reply.Status = storagenodeproto.OK
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) AddSongRequest(args *storagenodeproto.SongArgs, reply *storagenodeproto.SongReply) error {
	fmt.Printf("AddSong [%s]->[%s] Request Received from '%s' \n", args.SongName, args.PlaylistName, args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	_, exists := ss.userDB[targetUser][args.PlaylistName]
	if !exists {
		fmt.Println("\tError: Object Not Found!")
		reply.Status = storagenodeproto.OBJECT_NOT_FOUND
		ss.userLocks[targetUser].Unlock()
		return nil
	}
	
	if ss.userDB[targetUser][args.PlaylistName].Contains(args.SongName) {
		fmt.Println("\tError: Object Exists!")
		reply.Status = storagenodeproto.OBJECT_EXISTS
	} else {
		ss.userDB[targetUser][args.PlaylistName].Insert(args.SongName)
		reply.Status = storagenodeproto.OK
	}
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) DeleteSongRequest(args *storagenodeproto.SongArgs, reply *storagenodeproto.SongReply) error {
	fmt.Printf("DeleteSong [%s]->[%s] Request Received from '%s' \n", args.SongName, args.PlaylistName, args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	_, exists := ss.userDB[targetUser][args.PlaylistName]
	if !exists {
		fmt.Println("\tError: Object Not Found!")
		reply.Status = storagenodeproto.OBJECT_NOT_FOUND
		ss.userLocks[targetUser].Unlock()
		return nil
	}
	
	if ss.userDB[targetUser][args.PlaylistName].Contains(args.SongName) {
		err := ss.userDB[targetUser][args.PlaylistName].Remove(args.SongName)
		if err != nil {
			fmt.Println("\tError: ", err.Error())
			reply.Status = storagenodeproto.FAILED
		} else {
			reply.Status = storagenodeproto.OK
		}
	} else {
		fmt.Println("\tError: Object Not Found!")
		reply.Status = storagenodeproto.OBJECT_NOT_FOUND
	}
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) ViewAllSongsRequest(args *storagenodeproto.SongArgs, reply *storagenodeproto.SongReply) error {
	fmt.Printf("ViewAllSongs Request for [%s] Received from '%s' \n", args.PlaylistName, args.CInfo.Username)
	reply.Status = storagenodeproto.FAILED
	targetUser := args.CInfo.Username
	ss.NewUsernameRequest(targetUser)
	
	ss.userLocks[targetUser].Lock()
	_, exists := ss.userDB[targetUser][args.PlaylistName]
	if !exists {
		fmt.Println("\tError: Object Not Found!")
		reply.Status = storagenodeproto.OBJECT_NOT_FOUND
		ss.userLocks[targetUser].Unlock()
		return nil
	}
	songArray := ss.userDB[targetUser][args.PlaylistName].ToArray()
	
	display := args.PlaylistName + ": {"
	for i := 0; i < len(songArray); i++ {
		display += songArray[i]
		display += ", "
	}
	display += "}\n"
	reply.DisplayStr = display
	reply.Status = storagenodeproto.OK
	ss.userLocks[targetUser].Unlock()
	return nil
}


func (ss *StorageServer) PlaySongRequest(args *storagenodeproto.PlayArgs, reply *storagenodeproto.PlayReply) error {
	fmt.Printf("PlaySong [%s] Request Received from '%s' \n", args.SongName, args.CInfo.Username)
	return nil
}


func (ss *StorageServer) SearchSongRequest(args *storagenodeproto.SearchArgs, reply *storagenodeproto.SearchReply) error {
	fmt.Printf("SearchSong [%s] Request Received from '%s' \n", args.SongName, args.CInfo.Username)
	return nil
}


func (ss *StorageServer) SearchArtistRequest(args *storagenodeproto.SearchArgs, reply *storagenodeproto.SearchReply) error {
	fmt.Printf("SearchArtist [%s] Request Received from '%s' \n", args.ArtistName, args.CInfo.Username)
	return nil
}


