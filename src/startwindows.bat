@echo off
echo "Starting Windows..."

start "Monitor" 	cmd /K "cd monitor/testmonitor"
start "StorageA" 	cmd /K "cd storage/teststorage"
start "StorageB" 	cmd /K "cd storage/teststorage"
start "Bridge" 		cmd /K "cd bridge/testbridge"
start "PeerA - Server"	cmd /K "cd peernode/testserver"
start "PeerA - Client" 	cmd /K "cd peernode/testclient"
start "PeerB - Server"	cmd /K "cd peernode/testserver"
start "PeerB - Client"	cmd /K "cd peernode/testclient"