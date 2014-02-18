package peernode

import (
	"fmt"
	"time"
	"log"
	"os"
	"strings"
	"net/rpc"
	"../proto"
	"../rpc"
	"../../bridge/proto"
	"../../consts"
	"../../util/songinfo"
)

const UpdatePeersWaitSecs = 10

type PeerServer struct {
	username string	// username of dedicated peer
	peerport int		// Port number of dedicated peer
	portnum int		// Port number of PeerServer
	hostport string		// Hostport of PeerServer
	assignedBridgeHP string	// Hostport of Assigned Bridge
	peerMap map[string] string	// Map of peer usernames to their hostports
	
	peersConn map[string] *rpc.Client	// list of peers connected to
	serversConn map[string] *rpc.Client	// list of Spothify servers connected to
	
	prpc *peernoderpc.PeerRPC	// RPC object used to register RPC functions	
}

func NewPeerServer(username string, portnum int, bridgeHP string) *PeerServer {
	ps := new(PeerServer)
	ps.username = username
	ps.peerport = 0
	ps.portnum = portnum
	ps.assignedBridgeHP = bridgeHP
	ps.hostport = fmt.Sprintf("localhost:%d", portnum)
	ps.peersConn = make(map[string] *rpc.Client)
	ps.serversConn = make(map[string] *rpc.Client)
	ps.peerMap = make(map[string] string)

	//Register RPC functions
	ps.prpc = peernoderpc.NewPeerRPC(ps)
	rpc.Register(ps.prpc)

	// Create userdata directory for user!
	currentPath, err := os.Getwd()
	if err!=nil {
		log.Println("Error getting working directory", err.Error())
	}
	userDir := fmt.Sprintf("%s\\..\\userdata\\%s", currentPath, ps.username)
	// Does user directory already exist?
	_, err = os.Stat(userDir)
	if os.IsNotExist(err) {
		// Does not exist, create directory
		os.Mkdir(userDir,  os.ModeDir)
	}
	
	// Connect to a Bridge
	ps.ConnectToBridge(ps.assignedBridgeHP)
	
	go ps.GetPeers()
	
	return ps
}


func (ps *PeerServer) ConnectToBridge (bridgeHP string) {
	err, conn := constants.DialToServer(bridgeHP)
	if err != nil {
		log.Fatal("Could not connect to All Storage Servers")
	}
	ps.serversConn[bridgeHP] = conn
	log.Printf("Peer Server [%s] connected to bridge server [%s]\n", ps.username, ps.assignedBridgeHP)
}


func (ps *PeerServer) GetPeers() error {
	newClientInfo := bridgenodeproto.ClientInfo{ps.username, ps.hostport}
	newArgs := bridgenodeproto.GetPeersArgs{newClientInfo}
	var newReply bridgenodeproto.GetPeersReply
	for {
		fmt.Printf("\n\t*Updating Peer List...")
		err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.GetPeerList", &newArgs, &newReply)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return err
		} else if newReply.Status != bridgenodeproto.OK {
			fmt.Println("Denied!")
			return nil
		}
		// No errors...
		ps.peerMap = newReply.PeersMap
		fmt.Printf("Updated! @ %s\n", time.Now().String())
		fmt.Println("\tPeers:", ps.peerMap, "\n")
		time.Sleep( UpdatePeersWaitSecs * time.Second )
	}
}


func (ps *PeerServer) CheckRequest ( req peernodeproto.PeerInfo ) bool {
	if req.Username == ps.username {//&& req.Portnum == ps.portnum {
		return true
	}
	fmt.Println("Message received from a different user/port!\n")
	return false
}


func (ps *PeerServer) SendMsgRequest(args *peernodeproto.SendMsgArgs, reply *peernodeproto.SendMsgReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("Message Received: ", args.Message, "from", args.PInfo.Username)
	
	if args.Recipient == peernodeproto.TO_BRIDGE {
		fmt.Println("Forwarding Message to a Bridge Server...", ps.assignedBridgeHP)
		newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
		newArgs := bridgenodeproto.SendMsgArgs{newClientInfo, args.Message}
		var newReply bridgenodeproto.SendMsgReply
		
		err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.SendMsgRequest", &newArgs, &newReply)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return err
		} else if reply.Status != bridgenodeproto.OK {
			fmt.Printf("Denied!\n")
			return nil
		}
		
	} else if args.Recipient == peernodeproto.TO_OTHERPEER {
		fmt.Println("Forwarding Message to another peer server...")
	} else {
		fmt.Println("Invalid Recipient")
		// Invalid Recipient
	}
	
	reply.Message = "Request Handled!"
	reply.Status = peernodeproto.OK
	fmt.Println("\tSendRequest done\n")
	return nil
}


func (ps *PeerServer) AddPLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("AddPlaylist [%s] Request Received from '%s'", args.TargetPlaylistName, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.ChangePLArgs{newClientInfo, args.TargetPlaylistName, ""}
	var newReply bridgenodeproto.ChangePLReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.AddPLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Printf("Denied! \n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) DeletePLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("DeletePlaylist [%s] Request Received from '%s'...", args.TargetPlaylistName, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.ChangePLArgs{newClientInfo, args.TargetPlaylistName, ""}
	var newReply bridgenodeproto.ChangePLReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.DeletePLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) SortPLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("SortPlaylist [%s] Request Received from '%s'", args.TargetPlaylistName, args.PInfo.Username)
	return nil
}


func (ps *PeerServer) RenamePLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("RenamePlaylist [%s:%s] Request Received from '%s'...", args.TargetPlaylistName, args.NewPlaylistName, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.ChangePLArgs{newClientInfo, args.TargetPlaylistName, args.NewPlaylistName}
	var newReply bridgenodeproto.ChangePLReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.RenamePLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Printf("Denied! \n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) DownloadPLRequest(args *peernodeproto.DownloadPLArgs, reply *peernodeproto.DownloadPLReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("DownloadPlaylist [%s] Request Received from '%s'", args.TargetPlaylistName, args.PInfo.Username)
	return nil
}


func (ps *PeerServer) ViewAllPLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("ViewAllPlaylists Request Received from '%s'...", args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.ChangePLArgs{newClientInfo, "", ""}
	var newReply bridgenodeproto.ChangePLReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.ViewAllPLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.DisplayStr = newReply.DisplayStr
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) AddSongRequest(args *peernodeproto.SongArgs, reply *peernodeproto.SongReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("AddSong [%s]->[%s] Request Received from '%s'...", args.SongName, args.PlaylistName, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.SongArgs{newClientInfo, args.SongName, args.PlaylistName}
	var newReply bridgenodeproto.SongReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.AddSongRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Printf("Denied! \n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) DeleteSongRequest(args *peernodeproto.SongArgs, reply *peernodeproto.SongReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("DeleteSong [%s]->[%s] Request Received from '%s'...", args.SongName, args.PlaylistName, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.SongArgs{newClientInfo, args.SongName, args.PlaylistName}
	var newReply bridgenodeproto.SongReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.DeleteSongRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Println("Denied! \n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) ViewAllSongsRequest(args *peernodeproto.SongArgs, reply *peernodeproto.SongReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Printf("ViewAllSongs in [%s] Request Received from '%s'...", args.PlaylistName, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.SongArgs{newClientInfo, "", args.PlaylistName}
	var newReply bridgenodeproto.SongReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.ViewAllSongsRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Println("Denied! \n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.DisplayStr = newReply.DisplayStr
	reply.Status = peernodeproto.OK
	return nil
}


func (ps *PeerServer) PlaySongRequest(args *peernodeproto.PlayArgs, reply *peernodeproto.PlayReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("PlaySong [%s] Request Received from '%s'", args.SInfo.Name, args.PInfo.Username)
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newSongInfo := songinfo.NewSong(args.SInfo.Name)
	newArgs := bridgenodeproto.PlayArgs{newClientInfo, *newSongInfo}
	var newReply bridgenodeproto.PlayReply
	
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.PlaySongRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		reply.Status = peernodeproto.FAILED
		return err
	} else if newReply.Status != bridgenodeproto.OK {
		fmt.Println("Denied! \n")
		reply.Status = peernodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	
	currentPath, errCD := os.Getwd()
	if errCD!=nil {
		log.Fatal("Error getting working directory", errCD.Error())
	}
	
	// Write song file into TMP directory (under ../userdata/%user%)!
	tmp_path := fmt.Sprintf("%s\\..\\userdata\\%s\\TMP", currentPath, ps.username)
	// Does TMP directory already exist?
	_, err = os.Stat(tmp_path)
	if os.IsNotExist(err) {
		// Does not exist, create directory
		fmt.Println("Create New Directory @", tmp_path)
		os.Mkdir(tmp_path,  os.ModeDir)
	}
	
	// Write song file into TMP Directory
	s := []string{tmp_path, args.SInfo.Name}
	NowPlayPath := strings.Join(s,"\\")
	fmt.Println("New Path:", NowPlayPath)
	
	// Write contents to TMP folder (NowPlayPath)
	newFO, err3 := os.Create(NowPlayPath)
	if err3 != nil {
		log.Println("Error creating New File", err3.Error())
		reply.Status = peernodeproto.FAILED
		return nil
	}
	n, err4 := newFO.Write(newReply.SongBytes)
	if err4 != nil {
		log.Println("Error writing New File", err4.Error())
		reply.Status = peernodeproto.FAILED
		return nil
	}
	if n != len(newReply.SongBytes) {
		log.Println("Error writing New File - num mismatch!", n, len(newReply.SongBytes))
		reply.Status = peernodeproto.FAILED
		return nil
	}
	
	reply.PlayPath = NowPlayPath
	reply.Status = newReply.Status
	return nil
}


func (ps *PeerServer) SearchSongRequest(args *peernodeproto.SearchArgs, reply *peernodeproto.SearchReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("SearchSong [%s] Request Received from '%s'", args.SongName, args.PInfo.Username)
	return nil
}


func (ps *PeerServer) SearchArtistRequest(args *peernodeproto.SearchArgs, reply *peernodeproto.SearchReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("SearchArtist [%s] Request Received from '%s'", args.ArtistName, args.PInfo.Username)
	return nil
}


func (ps *PeerServer) QuitRequest(args *peernodeproto.QuitArgs, reply *peernodeproto.QuitReply) error {
	if ps.CheckRequest(args.PInfo) != true {
		reply.Status = peernodeproto.INCORRECTPEER
		return nil
	}
	fmt.Println("Quit Request Received from '%s'", args.PInfo.Username)
	fmt.Println("Forwarding Request to Bridge Server...", ps.assignedBridgeHP)
	
	newClientInfo := bridgenodeproto.ClientInfo{args.PInfo.Username, ps.hostport}
	newArgs := bridgenodeproto.QuitArgs{newClientInfo}
	var newReply bridgenodeproto.QuitReply
	err := ps.serversConn[ps.assignedBridgeHP].Call("BridgeRPC.QuitRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err
	} else if reply.Status != bridgenodeproto.OK {
		fmt.Printf("Denied! \n")
		return nil
	}
		
	return nil
}