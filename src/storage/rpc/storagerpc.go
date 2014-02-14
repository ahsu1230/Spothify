package storagenoderpc

import (
	"../proto"
)

type StorageInterface interface {
	SendMsgRequest(*storagenodeproto.SendMsgArgs, *storagenodeproto.SendMsgReply) error
	
	AddPLRequest(*storagenodeproto.ChangePLArgs, *storagenodeproto.ChangePLReply) error
	DeletePLRequest(*storagenodeproto.ChangePLArgs, *storagenodeproto.ChangePLReply) error
	SortPLRequest(*storagenodeproto.ChangePLArgs, *storagenodeproto.ChangePLReply) error
	RenamePLRequest(*storagenodeproto.ChangePLArgs, *storagenodeproto.ChangePLReply) error
	DownloadPLRequest(*storagenodeproto.DownloadPLArgs, *storagenodeproto.DownloadPLReply) error
	ViewAllPLRequest(*storagenodeproto.ChangePLArgs, *storagenodeproto.ChangePLReply) error
	
	AddSongRequest(*storagenodeproto.SongArgs, *storagenodeproto.SongReply) error
	DeleteSongRequest(*storagenodeproto.SongArgs, *storagenodeproto.SongReply) error
	ViewAllSongsRequest(*storagenodeproto.SongArgs, *storagenodeproto.SongReply) error
	
	SearchSongRequest(*storagenodeproto.SearchArgs, *storagenodeproto.SearchReply) error
	SearchArtistRequest(*storagenodeproto.SearchArgs, *storagenodeproto.SearchReply) error
	
	PlaySongRequest(*storagenodeproto.PlayArgs, *storagenodeproto.PlayReply) error
}

type StorageRPC struct {
	ss StorageInterface
}

func NewStorageRPC(ss StorageInterface) *StorageRPC {
	return &StorageRPC{ss}
}

func (srpc *StorageRPC) SendMsgRequest(args *storagenodeproto.SendMsgArgs, reply *storagenodeproto.SendMsgReply) error {
	return srpc.ss.SendMsgRequest(args, reply)
}

func (srpc *StorageRPC) AddPLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	return srpc.ss.AddPLRequest(args, reply)
}
func (srpc *StorageRPC) DeletePLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	return srpc.ss.DeletePLRequest(args, reply)
}
func (srpc *StorageRPC) SortPLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	return srpc.ss.SortPLRequest(args, reply)
}
func (srpc *StorageRPC) RenamePLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	return srpc.ss.RenamePLRequest(args, reply)
}
func (srpc *StorageRPC) DownloadPLRequest(args *storagenodeproto.DownloadPLArgs, reply *storagenodeproto.DownloadPLReply) error {
	return srpc.ss.DownloadPLRequest(args, reply)
}
func (srpc *StorageRPC) ViewAllPLRequest(args *storagenodeproto.ChangePLArgs, reply *storagenodeproto.ChangePLReply) error {
	return srpc.ss.ViewAllPLRequest(args, reply)
}

func (srpc *StorageRPC) AddSongRequest(args *storagenodeproto.SongArgs, reply *storagenodeproto.SongReply) error {
	return srpc.ss.AddSongRequest(args, reply)
}
func (srpc *StorageRPC) DeleteSongRequest(args *storagenodeproto.SongArgs, reply *storagenodeproto.SongReply) error {
	return srpc.ss.DeleteSongRequest(args, reply)
}
func (srpc *StorageRPC) ViewAllSongsRequest(args *storagenodeproto.SongArgs, reply *storagenodeproto.SongReply) error {
	return srpc.ss.ViewAllSongsRequest(args, reply)
}

func (srpc *StorageRPC) SearchSongRequest(args *storagenodeproto.SearchArgs, reply *storagenodeproto.SearchReply) error {
	return srpc.ss.SearchSongRequest(args, reply)
}
func (srpc *StorageRPC) SearchArtistRequest(args *storagenodeproto.SearchArgs, reply *storagenodeproto.SearchReply) error {
	return srpc.ss.SearchArtistRequest(args, reply)
}

func (srpc *StorageRPC) PlaySongRequest(args *storagenodeproto.PlayArgs, reply *storagenodeproto.PlayReply) error {
	return srpc.ss.PlaySongRequest(args, reply)
}

