package monitorproto

import(
	//"../../util/stringlist"
)

// Status codes
const (
	OK = iota
	STORAGES_NOT_READY
)

// Object Types
const (
	UNKNOWN = iota
	STORAGE
	BRIDGE
	STORAGE_BACKUP
	BRIDGE_BACKUP
)

type Node struct {
	Type int
	Hostport string
	ID	uint32
}

func DisplayServerType(nodetype int) string {
	display := ""
	switch nodetype {
	case STORAGE:
		display = "Storage"
	case BRIDGE:
		display = "Bridge"
	case STORAGE_BACKUP:
		display = "Storage Backup"
	case BRIDGE_BACKUP:
		display = "Bridge Backup"
	default:
		display = "Unknown"
	}
	return display
}

// RPC Function Parameters (Arguments & Replies)
/* Server-called RPCs */
type RegisterArgs struct {
	NodeInfo Node
}
type RegisterReply struct {
	Ready bool
	StorageMap	map[uint32] string		// For Bridges, list of storage HPs and respective Id's
	BridgeID uint32
	Status int
}

