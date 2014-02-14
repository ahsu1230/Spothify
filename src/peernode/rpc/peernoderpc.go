package peernoderpc

import (
	"../proto"
)

type PeerInterface interface {
	SendMsgRequest(*peernodeproto.SendMsgArgs, *peernodeproto.SendMsgReply) error
	QuitRequest(*peernodeproto.QuitArgs, *peernodeproto.QuitReply) error
	
	AddPLRequest(*peernodeproto.ChangePLArgs, *peernodeproto.ChangePLReply) error
	DeletePLRequest(*peernodeproto.ChangePLArgs, *peernodeproto.ChangePLReply) error
	SortPLRequest(*peernodeproto.ChangePLArgs, *peernodeproto.ChangePLReply) error
	RenamePLRequest(*peernodeproto.ChangePLArgs, *peernodeproto.ChangePLReply) error
	DownloadPLRequest(*peernodeproto.DownloadPLArgs, *peernodeproto.DownloadPLReply) error
	ViewAllPLRequest(*peernodeproto.ChangePLArgs, *peernodeproto.ChangePLReply) error
	
	AddSongRequest(*peernodeproto.SongArgs, *peernodeproto.SongReply) error
	DeleteSongRequest(*peernodeproto.SongArgs, *peernodeproto.SongReply) error
	ViewAllSongsRequest(*peernodeproto.SongArgs, *peernodeproto.SongReply) error
	
	SearchSongRequest(*peernodeproto.SearchArgs, *peernodeproto.SearchReply) error
	SearchArtistRequest(*peernodeproto.SearchArgs, *peernodeproto.SearchReply) error
	
	PlaySongRequest(*peernodeproto.PlayArgs, *peernodeproto.PlayReply) error
}

type PeerRPC struct {
	ps PeerInterface
}

func NewPeerRPC(ps PeerInterface) *PeerRPC {
	return &PeerRPC{ps}
}

func (prpc *PeerRPC) SendMsgRequest(args *peernodeproto.SendMsgArgs, reply *peernodeproto.SendMsgReply) error {
	return prpc.ps.SendMsgRequest(args, reply)
}
func (prpc *PeerRPC) AddPLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	return prpc.ps.AddPLRequest(args, reply)
}
func (prpc *PeerRPC) DeletePLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	return prpc.ps.DeletePLRequest(args, reply)
}
func (prpc *PeerRPC) SortPLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	return prpc.ps.SortPLRequest(args, reply)
}
func (prpc *PeerRPC) RenamePLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	return prpc.ps.RenamePLRequest(args, reply)
}
func (prpc *PeerRPC) DownloadPLRequest(args *peernodeproto.DownloadPLArgs, reply *peernodeproto.DownloadPLReply) error {
	return prpc.ps.DownloadPLRequest(args, reply)
}
func (prpc *PeerRPC) ViewAllPLRequest(args *peernodeproto.ChangePLArgs, reply *peernodeproto.ChangePLReply) error {
	return prpc.ps.ViewAllPLRequest(args, reply)
}

func (prpc *PeerRPC) AddSongRequest(args *peernodeproto.SongArgs, reply *peernodeproto.SongReply) error {
	return prpc.ps.AddSongRequest(args, reply)
}
func (prpc *PeerRPC) DeleteSongRequest(args *peernodeproto.SongArgs, reply *peernodeproto.SongReply) error {
	return prpc.ps.DeleteSongRequest(args, reply)
}
func (prpc *PeerRPC) ViewAllSongsRequest(args *peernodeproto.SongArgs, reply *peernodeproto.SongReply) error {
	return prpc.ps.ViewAllSongsRequest(args, reply)
}

func (prpc *PeerRPC) SearchSongRequest(args *peernodeproto.SearchArgs, reply *peernodeproto.SearchReply) error {
	return prpc.ps.SearchSongRequest(args, reply)
}
func (prpc *PeerRPC) SearchArtistRequest(args *peernodeproto.SearchArgs, reply *peernodeproto.SearchReply) error {
	return prpc.ps.SearchArtistRequest(args, reply)
}

func (prpc *PeerRPC) PlaySongRequest(args *peernodeproto.PlayArgs, reply *peernodeproto.PlayReply) error {
	return prpc.ps.PlaySongRequest(args, reply)
}

func (prpc *PeerRPC) QuitRequest(args *peernodeproto.QuitArgs, reply *peernodeproto.QuitReply) error {
	return prpc.ps.QuitRequest(args, reply)
}