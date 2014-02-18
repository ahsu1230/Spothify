package songinfo

type SongInfo struct {
	Name string
	//Artist string
	//Album string
}


//func NewSong(newName, newArtist, newAlbum string) *SongInfo {
func NewSong(newName string) *SongInfo {
	newSong := new(SongInfo)
	newSong.Name = newName
	//newSong.Artist = newArtist
	//newSong.Album = newAlbum
	return newSong
}

func equals(s1, s2 SongInfo) bool{
	//return (s1.Name == s2.Name) && (s1.Artist == s2.Artist) && (s1.Album == s2.Album)
	return (s1.Name == s2.Name) 
}