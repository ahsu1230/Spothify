Project Spothify
Author: Aaron Hsu

Purpose:
 - To apply my skills of distributed systems and working with highly-scalable systems/products.
 - To recreate the the backend infrastructure of the music application, Spotify.

Objectives:
 - Spotify is a music streaming software service that focuses on high scale, low latency music streaming.

 - P2P Networking. 
		As more users join, low latency streaming comes from peers instead of servers.
 - Storage Servers. 
		Must be consistent with user information and fault tolerant. Also stores song data.
 - Bridge Servers. 
		Act as middlemen between Peers and Spotify Storage Servers. Contain large caches for quick access of popular data/songs.
 - Offline Peer Functionality. 
		Peers can choose to be offline and still be able to play local music.
 - Music Streaming.
		Music streaming is fast. Songs are fed through storage servers OR through peers.
		Songs are not downloaded all at once, they are usually split into 15 second chunks.

Peers can request the following commands:
 - Play a requested song
 - Add/Delete/Rename Playlist
 - Add/Delete Songs from a Playlist
 - Download a Playlist (download all songs in a selected playlist)
 - Offline Mode