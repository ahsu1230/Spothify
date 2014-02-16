Project Spothify
Author: Aaron Hsu
Last Updated: 2/10/2014

::::::::::::::: Accomplished / To-Do(*) Features :::::::::::::::

Peer Client
 - Most Bridge & Storage Interaction: Add/Delete/Rename/View Playlists, Add/Delete/View Songs, Quit
 - (*) PlaySong
 - (*) DownloadPlaylist & Saving Song Files
 - (*) PlaySong (Offline Mode)
 - (*) Offline Mode
 - (*) SortPlaylist, SearchSong, SearchArtist implementations (only sends messages to bridge / storage)

Peer Server
 - Responds to PeerClient Requests
 - Receives Bridge Responses
 - Recognizes if request comes from a wrong username
 - Asks Bridge Server about nearby peers and updates list of peers
 - (*) P2P PlaySong - look for song between peers


Bridge Server
 - Waits until all required Storage Server(s) are up
 - Responds to PeerServer Requests
 - Receives Storage Responses
 - Give Peers a list of their nearby peers
 - Consistent Hashing of requests to send to correct storage in Distributed Storage Cluster
 - (*) Cache popular song requests
 - (*) Cache highly active user info?
 - (*) What if server fails?


Storage Server
 - Waits to start until all required Storage Server(s) are up
 - Responds to BridgeServer Requests
 - Save User Information
 - Stores Song Data
 - Consistent Hashing & Split what each server handles
 - (*) Is it better to have a lead storage server (other storages/bridges register with leader storage)
 - (*) What if server fails?
 - (*) Need Replication mechanism


Monitor Server (?)
 - Waits until enough servers connect with monitor
 - Maintains Connection with all Storage servers
 - (*) What happens 
 - (*) What if server fails? Single point of failure!