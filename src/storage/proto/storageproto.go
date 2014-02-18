package storagenodeproto

import(
	"../../util/songinfo"
)

// Status codes
const (
	FAILED = iota
	OK
	OBJECT_EXISTS
	OBJECT_NOT_FOUND
)

// Object Types
const (
)

type ClientInfo struct {
	Username string
	BridgeHP string
}

// RPC Function Parameters (Arguments & Replies)
/* Server-called RPCs */
type SendMsgArgs struct {
	CInfo	ClientInfo
	Message string
}
type SendMsgReply struct {
	Message string
	Status int
}

/* Add, Delete, Sort, Rename Playlist */
type ChangePLArgs struct {
	CInfo	ClientInfo
	TargetPlaylistName string
	NewPlaylistName string
}
type ChangePLReply struct {
	DisplayStr	string
	Status int
}
type DownloadPLArgs struct {
	CInfo	ClientInfo
	TargetPlaylistName string
}
type DownloadPLReply struct {
	Status int
}

/* Add, Delete Song to Playlist */
type SongArgs struct {
	CInfo	ClientInfo
	SongName	string
	PlaylistName string
}
type SongReply struct {
	DisplayStr	string
	Status int
}

/* Search Song or Artist */
type SearchArgs struct {
	CInfo	ClientInfo
	SongName	string
	ArtistName	string
}
type SearchReply struct {
	Status int
}

/* Playing a requested song */
type PlayArgs struct {
	CInfo	ClientInfo
	SInfo	songinfo.SongInfo
}
type PlayReply struct {
	SongBytes		[]byte
	Status int
}