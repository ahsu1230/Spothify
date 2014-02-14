package songinfo

type SongInfo {
	name string
	artist string
	album string
}

func NewSong(newName, newArtist, newAlbum string) *SongInfo {
	newSong := new(SongInfo)
	newSong.name = newName
	newSong.artist = newArtist
	newSong.album = newAlbum
	return newSong
}

func equals(s1, s2 SongInfo) bool{
	return (s1.name == s2.name) && (s1.artist == s2.artist) && (s1.album == s2.artist)
}