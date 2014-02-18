package bridgenode

import (
	"fmt"
	"net/rpc"
	"../proto"
	"../rpc"
	"log"
	"time"
	"../../storage/proto"
	"../../monitor/proto"
	"../../consts"
	//"../../util/stringlist"
	"../../util/uint32list"
	"../../util/hasher"
	"../../util/songinfo"
)

type SongInfo struct {
	name	string
	artist	string
}

type BridgeServer struct {
	id uint32			// Bridge ID
	portnum int			// Portnum of Bridge Server
	hostport string		// Hostport of Bridge Server
	monitorHP string		// Hostport of Monitor Server
		
	serverAddress map[uint32] string	// list of server node IDs -> their hostports
	connected map[string] *rpc.Client	// list of servers connected to
	
	peers map[string] string		// map peer username to its HP
	servers map[uint32] string		// maps storage IDs to HP
	serverIDs *uint32list.List		// dynamic list of server IDs from lowest to highest 
							// 	allows insert/delete of server IDs
	// cached requests
	
	brpc *bridgenoderpc.BridgeRPC	// RPC object used to register RPC functions
}


func NewBridgeServer(portnum int, monitorPort int) *BridgeServer {
	bs := new(BridgeServer)
	bs.portnum = portnum
	bs.hostport = fmt.Sprintf("localhost:%d", portnum)
	bs.monitorHP = fmt.Sprintf("localhost:%d", monitorPort)
	bs.serverAddress = make(map[uint32] string)
	bs.connected = make(map[string] *rpc.Client)
	bs.peers = make(map[string] string)
	bs.servers = make(map[uint32] string)
	bs.serverIDs = uint32list.NewList()

	//Register RPC functions
	bs.brpc = bridgenoderpc.NewBridgeRPC(bs)
	rpc.Register(bs.brpc)
	
	// Register Bridge Server with Monitor - reply with list of storage servers.
	bs.RegisterWithMonitor()
	// Dial to all storage servers
	bs.DialToStorages()
	
	return bs
}


func (bs *BridgeServer) RegisterWithMonitor() {
	fmt.Println("Registering with Monitor...")
	err, conn := constants.DialToServer(bs.monitorHP)
	if err != nil {
		log.Fatal("Problem Detected Dialing to Server", bs.monitorHP)
	}
	bs.connected[bs.monitorHP] = conn
	
	registerInfo := monitorproto.Node{monitorproto.BRIDGE, bs.hostport, bs.id}
	newArgs := monitorproto.RegisterArgs{registerInfo}
	var newReply monitorproto.RegisterReply
	newReply.Ready = false
	
	for !newReply.Ready {
		err = conn.Call("MonitorRPC.RegisterServer", &newArgs, &newReply)
		if err != nil {
			fmt.Println("\tError: ", err.Error())
			log.Fatal(err.Error())
		} else if newReply.Status != monitorproto.OK {
			fmt.Println("Denied!")
			log.Fatal("Registration Denied!")
		}
		time.Sleep (constants.RpcWaitMillis * time.Millisecond)
	}
	fmt.Println("Returned!")
	// Should have returned in reply, map of storage server IDs to hostports
	for id, hp := range(newReply.StorageMap) {
		bs.serverAddress[id] = hp
	}
	bs.id = newReply.BridgeID
	bs.servers = newReply.StorageMap // map[uint32 id] -> string HP
	
	fmt.Println("Try to Sort...")
	for targetID := range(bs.servers) {
		// insertion sort (easy & only one time use & probably not many items needed to sort)
		fmt.Println("\t Inserting...", targetID)
		bs.serverIDs.InsertInSort(targetID)
	}
	fmt.Println("Sorted List of Servers...", bs.serverIDs.ToArray())
	fmt.Println("Registered and Ready with all Storage Servers and BridgeID", bs.id)
}


func (bs *BridgeServer) DialToStorages() {
	for _, hp := range(bs.serverAddress) {
		err, conn := constants.DialToServer(hp)
		if err != nil {
			log.Fatal("Could not connect to All Storage Servers")
		}
		bs.connected[hp] = conn
	}
}


func (bs *BridgeServer) GetPeerList(args *bridgenodeproto.GetPeersArgs, reply *bridgenodeproto.GetPeersReply) error {
	// Peer Server calls this RPC to register themselves with bridge and obtain a list of peers
	// Add peer to bridge's peer list with 
	_, exists := bs.peers[args.CInfo.Username]
	if !exists {
		fmt.Printf("\tAdd New User %s!\n", args.CInfo.Username)
		bs.peers[args.CInfo.Username] = args.CInfo.Hostport // need to verify it's ok?
	}
	//fmt.Printf("Update %s Peer List @ [%s]\n", args.CInfo.Username, args.CInfo.Hostport)
	// Return with bridge's peer list
	reply.PeersMap = bs.peers
	reply.Status = bridgenodeproto.OK
	return nil
}


// Given a key, we hash it according to hasher package to get a hash value
// We use hash value to figure out which server is associated
func (bs *BridgeServer) GetStorageHP(key string) string {
	// Remember, bs.serverIDs are sorted...
	targetID := hasher.Storehash(key)
	fmt.Printf("\ttargetID: %d\n", targetID)
	idArray := bs.serverIDs.ToArray()
	c := 0
	for (c < len(idArray) && targetID > idArray[c]) {
		c++
	}
	// this case happens when the targetID surpasses the "largest" server ID
	// in that case, consistent hashing requires us to loop around, so it's the first serverID
	if c == len(idArray) {
		return bs.servers[idArray[0]]
	}
	// otherwise return appropriate server id
	return bs.servers[idArray[c]]
}

// ------------------------- API Functions -------------------------

func (bs *BridgeServer) SendMsgRequest(args *bridgenodeproto.SendMsgArgs, reply *bridgenodeproto.SendMsgReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	
	fmt.Println("Message Received: ", args.Message, "from", args.CInfo.Username)
	
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.SendMsgArgs{newClientInfo, args.Message}
	var newReply storagenodeproto.SendMsgReply
	fmt.Println("Forwarding Message to Storage...")
	reply.Status = bridgenodeproto.FAILED
	err := bs.connected[constants.SAMPLE_STORAGEHP].Call("StorageRPC.SendMsgRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if reply.Status != storagenodeproto.OK {
		fmt.Println("Denied!")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	
	reply.Message = "Request Handled!"
	fmt.Println("\tSendRequest done\n")
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) AddPLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("AddPlaylist [%s] Request Received from '%s'...", args.TargetPlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.ChangePLArgs{newClientInfo, args.TargetPlaylistName, ""}
	var newReply storagenodeproto.ChangePLReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.AddPLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied! \n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) DeletePLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("DeletePlaylist [%s] Request Received from '%s'...", args.TargetPlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.ChangePLArgs{newClientInfo, args.TargetPlaylistName, ""}
	var newReply storagenodeproto.ChangePLReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.DeletePLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied! \n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) SortPLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Println("SortPlaylist [%s] Request Received from '%s'", args.TargetPlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) RenamePLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("RenamePlaylist [%s:%s] Request Received from '%s'...", args.TargetPlaylistName, args.NewPlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.ChangePLArgs{newClientInfo, args.TargetPlaylistName, args.NewPlaylistName}
	var newReply storagenodeproto.ChangePLReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.RenamePLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied! \n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) ViewAllPLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("ViewAllPlaylists Request Received from '%s'...", args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.ChangePLArgs{newClientInfo, "", ""}
	var newReply storagenodeproto.ChangePLReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.ViewAllPLRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.DisplayStr = newReply.DisplayStr
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) DownloadPLRequest(args *bridgenodeproto.DownloadPLArgs, reply *bridgenodeproto.DownloadPLReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Println("DownloadPlaylist [%s] Request Received from '%s'", args.TargetPlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) AddSongRequest(args *bridgenodeproto.SongArgs, reply *bridgenodeproto.SongReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("AddSong [%s]->[%s] Request Received from '%s'...", args.SongName, args.PlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.SongArgs{newClientInfo, args.SongName, args.PlaylistName}
	var newReply storagenodeproto.SongReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.AddSongRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) DeleteSongRequest(args *bridgenodeproto.SongArgs, reply *bridgenodeproto.SongReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("DeleteSong [%s]->[%s] Request Received from '%s'...", args.SongName, args.PlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.SongArgs{newClientInfo, args.SongName, args.PlaylistName}
	var newReply storagenodeproto.SongReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.DeleteSongRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) ViewAllSongsRequest(args *bridgenodeproto.SongArgs, reply *bridgenodeproto.SongReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Printf("ViewAllSongs Request [%s] Received from '%s'\n", args.PlaylistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newArgs := storagenodeproto.SongArgs{newClientInfo, "", args.PlaylistName}
	var newReply storagenodeproto.SongReply
	targetHP := bs.GetStorageHP(args.CInfo.Username)
	err := bs.connected[targetHP].Call("StorageRPC.ViewAllSongsRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.DisplayStr = newReply.DisplayStr
	reply.Status = bridgenodeproto.OK
	return nil
}



func (bs *BridgeServer) PlaySongRequest(args *bridgenodeproto.PlayArgs, reply *bridgenodeproto.PlayReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Println("PlaySong [%s] Request Received from '%s'", args.SInfo.Name, args.CInfo.Username)

	reply.Status = bridgenodeproto.FAILED
	newClientInfo := storagenodeproto.ClientInfo{args.CInfo.Username, bs.hostport}
	newSongInfo := songinfo.NewSong(args.SInfo.Name)
	newArgs := storagenodeproto.PlayArgs{newClientInfo, *newSongInfo}
	var newReply storagenodeproto.PlayReply
	
	// Hash by song name NOT username (like most operations)
	targetHP := bs.GetStorageHP(args.SInfo.Name)
	err := bs.connected[targetHP].Call("StorageRPC.PlaySongRequest", &newArgs, &newReply)
	if err != nil {
		fmt.Println("\tError: ", err.Error())
		return err
	} else if newReply.Status != storagenodeproto.OK {
		fmt.Printf("Denied!\n")
		reply.Status = bridgenodeproto.FAILED
		return nil
	}
	fmt.Printf("Success! \n")
	reply.SongBytes = newReply.SongBytes
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) SearchSongRequest(args *bridgenodeproto.SearchArgs, reply *bridgenodeproto.SearchReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Println("SearchSong [%s] Request Received from '%s'", args.SongName, args.CInfo.Username)
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) SearchArtistRequest(args *bridgenodeproto.SearchArgs, reply *bridgenodeproto.SearchReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Println("SearchArtist [%s] Request Received from '%s'", args.ArtistName, args.CInfo.Username)
	reply.Status = bridgenodeproto.OK
	return nil
}


func (bs *BridgeServer) QuitRequest(args *bridgenodeproto.QuitArgs, reply *bridgenodeproto.QuitReply) error {
	_, peerexists := bs.peers[args.CInfo.Username]
	if !peerexists {
		reply.Status = bridgenodeproto.UNREGISTERED
		return nil
	}
	fmt.Println("Quit Request Received from '%s'", args.CInfo.Username)
	
	// Delete username from map of peers
	delete (bs.peers, args.CInfo.Username)
	
	reply.Status = bridgenodeproto.OK
	return nil
}