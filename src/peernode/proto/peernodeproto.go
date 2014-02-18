package peernodeproto

import(
	"../../util/songinfo"
)

// Status codes
const (
	OK = iota
	INCORRECTPEER
	OFFLINEERROR
	FAILED
)

// Direction codes
const (
	TO_OTHERPEER = iota
	TO_BRIDGE
	TO_STORAGE
	TO_MONITOR
)

// Object Types
const (
)

type PeerInfo struct {
	Username string
	Portnum int
}

type RecipientInfo struct {
	Recipient 	int
	hostport 	string
	username 	string	// For peers only
}

// RPC Function Parameters (Arguments & Replies)
/* Server-called RPCs */
type SendMsgArgs struct {
	PInfo		PeerInfo
	Recipient 	int
	//RInfo		RecipientInfo
	Message 	string
}
type SendMsgReply struct {
	Message string
	Status int
}

/* Add, Delete, Sort, Rename Playlist */
type ChangePLArgs struct {
	PInfo	PeerInfo
	TargetPlaylistName string
	NewPlaylistName string
}
type ChangePLReply struct {
	DisplayStr string
	Status int
}
type DownloadPLArgs struct {
	PInfo	PeerInfo
	TargetPlaylistName string
}
type DownloadPLReply struct {
	Status int
}

/* Add, Delete Song to Playlist */
type SongArgs struct {
	PInfo	PeerInfo
	SongName	string
	PlaylistName string
}
type SongReply struct {
	DisplayStr string
	Status int
}

/* Search Song or Artist */
type SearchArgs struct {
	PInfo	PeerInfo
	SongName	string
	ArtistName	string
}
type SearchReply struct {
	Status int
}

/* Playing a requested song */
type PlayArgs struct {
	PInfo	PeerInfo
	SInfo	songinfo.SongInfo
}
type PlayReply struct {
	PlayPath	string
	Status int
}

/* Quit */
type QuitArgs struct {
	PInfo	PeerInfo
}
type QuitReply struct {
	Status int
}