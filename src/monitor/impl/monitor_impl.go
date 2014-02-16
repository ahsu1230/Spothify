package monitor

import (
	"fmt"
	"net/rpc"
	"../proto"
	"../rpc"
	"sync"
	"log"
	//"../../consts"
	//"../../util/stringlist"
)

type MonitorServer struct {
	portnum int		// Port number of Monitor Server
	hostport string		// Hostport of Monitor Server
	numstorage int		// Expected Number of Storage Nodes
	storageCnt int		// Number of storage nodes registered so far
	bridgeCnt uint32	// Number of bridge nodes registered so far
	
	registerMutex *sync.RWMutex
	storageConn map[string] *rpc.Client		// list of storages connected to	
	bridgeConn map[string] *rpc.Client 		// list of bridges connected to
	servers map[string] monitorproto.Node	// list of servers
	//serverIDs map[uint32] string			// maps serverIDs to respective hostports
	
	mrpc *monitorrpc.MonitorRPC	// RPC object used to register RPC functions
}


func NewMonitorServer(portnum int, numstorage int) *MonitorServer {
	ms := new(MonitorServer)
	ms.portnum = portnum
	ms.numstorage = numstorage
	ms.storageCnt = 0
	ms.hostport = fmt.Sprintf("localhost:%d", portnum)
	
	ms.registerMutex = new(sync.RWMutex)
	ms.storageConn = make(map[string] *rpc.Client)
	ms.bridgeConn = make(map[string] *rpc.Client)
	ms.servers = make(map[string] monitorproto.Node)
	//ms.serverIDs = make(map[uint32] string)
	
	//Register RPC functions
	ms.mrpc = monitorrpc.NewMonitorRPC(ms)
	rpc.Register(ms.mrpc)
	return ms
}


func (ms *MonitorServer) RegisterServer(args *monitorproto.RegisterArgs, reply *monitorproto.RegisterReply) error {
	// add new storage server to map
	reply.Ready = false

	ms.registerMutex.Lock()
	targetHostport := args.NodeInfo.Hostport

	// Check if node has already notified its ready
	_, exists := ms.servers[targetHostport]
	if exists {
		// Checks if all nodes registered yet -> do nothing if not
		ms.CheckAllNodesRegistered(args, reply)
		ms.registerMutex.Unlock()
		return nil
	}

	// Mark the node as registered and add its info to the node list
	ms.servers[targetHostport] = args.NodeInfo
	
	switch args.NodeInfo.Type {
	case monitorproto.STORAGE:
		/*
		_, exists := ms.serverIDs[args.NodeInfo.ID]
		if !exists {
			ms.serverIDs[args.NodeInfo.ID] = args.NodeInfo.Hostport
		} else {
			log.Printf("Storage server with ID %d already registered!", args.NodeInfo.ID)
		}
		//*/
		
		ms.storageCnt++
		if ms.storageCnt >= ms.numstorage {
			log.Println("All Storage Servers have joined with Monitor!")
		}
	case monitorproto.BRIDGE:
		// save stuff in monitor?
		// can't do anything until all storages registered
		// assign ID
		reply.BridgeID = ms.bridgeCnt
		ms.bridgeCnt++
	}

	log.Printf("[%s] Server @ %s has joined and is connected\n", monitorproto.DisplayServerType(args.NodeInfo.Type), args.NodeInfo.Hostport)
	
	// Checks if all nodes registered yet -> do nothing if not
	ms.CheckAllNodesRegistered(args, reply)
	ms.registerMutex.Unlock()
	reply.Status = monitorproto.OK
	return nil
}


/* Checks if all nodes have registered and if they have, 
   set reply RPC fields */
func (ms *MonitorServer) CheckAllNodesRegistered(args *monitorproto.RegisterArgs, reply *monitorproto.RegisterReply) bool {
	reply.Ready = false
	// If all storage nodes are registered, set ready status to true
	if ms.storageCnt >= ms.numstorage {
		reply.Ready = true
		targetNode := args.NodeInfo
		switch targetNode.Type {
		case monitorproto.STORAGE:

		case monitorproto.BRIDGE:
			// reply.list of storages
			newStorageMap := make(map[uint32] string)
			for currentHP := range(ms.servers) {
				targetServer := ms.servers[currentHP]
				if targetServer.Type == monitorproto.STORAGE {
					newStorageMap[targetServer.ID] = targetServer.Hostport
				}
			}
			reply.StorageMap = newStorageMap
		}
	}
	return reply.Ready
}