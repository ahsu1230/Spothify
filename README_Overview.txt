Project Spothify
Author: Aaron Hsu
Last Updated: 2/10/2014
	

:::::::::::::::File Source Tree:::::::::::::::

src/Test_Music	: tests how to play/copy music data using GO
	- play_music.go		: GO test file
	- wmplayer.exe		: Windows Media Player Executable (can really be any exe)


src/peernode	: Source code for PEER node functionalities (includes client and server functionalities)
	/impl/client_impl.go		: Implementation of Peer Client API functions
	/impl/server_impl.go		: Implementation of Peer Server API functions
	/proto/peernodeproto.go		: Definitions of Peer RPC Parameters
	/rpc/peernoderpc.go		: Registers Objects for RPCs to be accessible remotely
	/testclient/testclient.go	: Build & Run! Tests an instance of a Peer Client (*)
						Parameters: -port=(#) -bridge=("") [username] [cmd] [args]
	/testserver/testserver.go	: Build & Run! Tests an instance of a Peer Server
						Parameters: -port=(#) -bridge=("") [username]

	(*) = Commands & Arguments for PeerClient:
	ap	"Add Playlist" (playlist_name)
	dp	"Delete Playlist" (playlist_name)
	rp	"Rename Playlist" (playlist_name) (new_playlist_name)
	vp	"View Playlists" ()
	dlp	"Download Playlist" (playlist_name)
	ps	"Play Song"	(song_name)
	as	"Add Song"	(playlist_name) (song_name)
	ds	"Delete Song"	(playlist_name) (song_name)
	vs	"View Songs"	(playlist_name)
	q	"Quit"		()

	More to come...
	sp	"Sort Playlist"
	ss	"Search Song"
	sa	"Search Artist"
	

src/bridge	: Source code for BRIDGE node functionalities
	/impl/bridge_impl.go		: Implementation of BRIDGE Server API functions
	/proto/bridgeproto.go		: Definitions of BRIDGE RPC Parameters
	/rpc/bridgerpc.go		: Registers Objects for RPCs to be accessible remotely
	/testbridge/testbridge.go	: Build & Run! Tests an instance of a BRIDGE Server
						Parameters: -port=(#)


src/storage	: Source code for STORAGE node functionalities
	/impl/storage_impl.go		: Implementation of STORAGE Server API functions
	/proto/storageproto.go		: Definitions of STORAGE RPC Parameters
	/rpc/storagerpc.go		: Registers Objects for RPCs to be accessible remotely
	/testbridge/teststorage.go	: Build & Run! Tests an instance of a STORAGE Server
						Parameters: -port=(#)


src/monitor	: Source code for MONITOR node functionalities
	/impl/monitor_impl.go		: Implementation of MONITOR Server API functions
	/proto/monitorproto.go		: Definitions of MONITOR RPC Parameters
	/rpc/monitorrpc.go		: Registers Objects for RPCs to be accessible remotely
	/testmonitor/testmonitor.go	: Build & Run! Tests an instance of a MONITOR Server
						Parameters: -port=(#) -n=(#num of storage servers)


src/util	: Supporting Data Structures / Functions for universal utility
	/filetraverser.go	: Directory/File Traversing/Handling... not entirely built into GO...
	/songinfo.go		: Struct for Song Information (Name, Artist, Album, etc.)
	/songinfolist.go	: Doubly LL / Hash structure(**) for handling SongInfo (used for caching)
				     Supports Move-To-Front, Front, Push, Pop, Contains, ...
	/stringlist.go		: Doubly LL / Hash structure(**) for handling strings (used to keep lists of hostports)
				     Supports Insert, Remove, Contains, ...
	/uint32list.go		: Doubly LL / Hash structure(**) for handling uint32 (used to keep lists of ids/ports?)
				     Supports Insert, Remove, Contains, ...

	(**) = structure has O(n) size, but has O(1) operations. GO does not have a strictly defined linked list interface



src/consts	: Universal Definitions for entire Project
	/consts.go		: Constant Definitions



tests/		: Test Scripts



:::::::::::::::Sample Test::::::::::::::: (2 peers, 1 bridge, 1 server, monitor defaults to port 9009)

testclient.exe -port=8888 ahsu [command & args]
testserver.exe -port=8888 -bridge=localhost:9876 ahsu
testclient.exe -port=8889 jfan [command & args]
testserver.exe -port=8889 -bridge=localhost:9876 jfan

testbridge.exe -port=9876
teststorage.exe -port=9321
testmonitor.exe -n=1


testclient.exe -port=8889 jfan [command & args]
testserver.exe -port=8889 -bridge=localhost:9876 jfan