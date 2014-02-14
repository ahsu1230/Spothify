package bridgenoderpc

import (
	"../proto"
)

type BridgeInterface interface {
	GetPeerList(*bridgenodeproto.GetPeersArgs, *bridgenodeproto.GetPeersReply) error
	SendMsgRequest(*bridgenodeproto.SendMsgArgs, *bridgenodeproto.SendMsgReply) error
	QuitRequest(*bridgenodeproto.QuitArgs, *bridgenodeproto.QuitReply) error
	
	AddPLRequest(*bridgenodeproto.ChangePLArgs, *bridgenodeproto.ChangePLReply) error
	DeletePLRequest(*bridgenodeproto.ChangePLArgs, *bridgenodeproto.ChangePLReply) error
	SortPLRequest(*bridgenodeproto.ChangePLArgs, *bridgenodeproto.ChangePLReply) error
	RenamePLRequest(*bridgenodeproto.ChangePLArgs, *bridgenodeproto.ChangePLReply) error
	DownloadPLRequest(*bridgenodeproto.DownloadPLArgs, *bridgenodeproto.DownloadPLReply) error
	ViewAllPLRequest(*bridgenodeproto.ChangePLArgs, *bridgenodeproto.ChangePLReply) error
	
	AddSongRequest(*bridgenodeproto.SongArgs, *bridgenodeproto.SongReply) error
	DeleteSongRequest(*bridgenodeproto.SongArgs, *bridgenodeproto.SongReply) error
	ViewAllSongsRequest(*bridgenodeproto.SongArgs, *bridgenodeproto.SongReply) error
	
	SearchSongRequest(*bridgenodeproto.SearchArgs, *bridgenodeproto.SearchReply) error
	SearchArtistRequest(*bridgenodeproto.SearchArgs, *bridgenodeproto.SearchReply) error
	
	PlaySongRequest(*bridgenodeproto.PlayArgs, *bridgenodeproto.PlayReply) error
}

type BridgeRPC struct {
	bs BridgeInterface
}

func NewBridgeRPC(bs BridgeInterface) *BridgeRPC {
	return &BridgeRPC{bs}
}

func (brpc *BridgeRPC) GetPeerList(args *bridgenodeproto.GetPeersArgs, reply *bridgenodeproto.GetPeersReply) error {
	return brpc.bs.GetPeerList(args, reply)
}

func (brpc *BridgeRPC) SendMsgRequest(args *bridgenodeproto.SendMsgArgs, reply *bridgenodeproto.SendMsgReply) error {
	return brpc.bs.SendMsgRequest(args, reply)
}

func (brpc *BridgeRPC) AddPLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	return brpc.bs.AddPLRequest(args, reply)
}
func (brpc *BridgeRPC) DeletePLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	return brpc.bs.DeletePLRequest(args, reply)
}
func (brpc *BridgeRPC) SortPLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	return brpc.bs.SortPLRequest(args, reply)
}
func (brpc *BridgeRPC) RenamePLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	return brpc.bs.RenamePLRequest(args, reply)
}
func (brpc *BridgeRPC) DownloadPLRequest(args *bridgenodeproto.DownloadPLArgs, reply *bridgenodeproto.DownloadPLReply) error {
	return brpc.bs.DownloadPLRequest(args, reply)
}
func (brpc *BridgeRPC) ViewAllPLRequest(args *bridgenodeproto.ChangePLArgs, reply *bridgenodeproto.ChangePLReply) error {
	return brpc.bs.ViewAllPLRequest(args, reply)
}

func (brpc *BridgeRPC) AddSongRequest(args *bridgenodeproto.SongArgs, reply *bridgenodeproto.SongReply) error {
	return brpc.bs.AddSongRequest(args, reply)
}
func (brpc *BridgeRPC) DeleteSongRequest(args *bridgenodeproto.SongArgs, reply *bridgenodeproto.SongReply) error {
	return brpc.bs.DeleteSongRequest(args, reply)
}
func (brpc *BridgeRPC) ViewAllSongsRequest(args *bridgenodeproto.SongArgs, reply *bridgenodeproto.SongReply) error {
	return brpc.bs.ViewAllSongsRequest(args, reply)
}


func (brpc *BridgeRPC) SearchSongRequest(args *bridgenodeproto.SearchArgs, reply *bridgenodeproto.SearchReply) error {
	return brpc.bs.SearchSongRequest(args, reply)
}
func (brpc *BridgeRPC) SearchArtistRequest(args *bridgenodeproto.SearchArgs, reply *bridgenodeproto.SearchReply) error {
	return brpc.bs.SearchArtistRequest(args, reply)
}

func (brpc *BridgeRPC) PlaySongRequest(args *bridgenodeproto.PlayArgs, reply *bridgenodeproto.PlayReply) error {
	return brpc.bs.PlaySongRequest(args, reply)
}

func (brpc *BridgeRPC) QuitRequest(args *bridgenodeproto.QuitArgs, reply *bridgenodeproto.QuitReply) error {
	return brpc.bs.QuitRequest(args, reply)
}

