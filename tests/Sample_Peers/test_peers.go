package main

import (
    "fmt"
    "time"
)

func Sleeper(c chan int) {
    time.Sleep(1500 * time.Millisecond)
    c <- 1
}

func LookForSong(peerID int, songlist []string, song string, c chan int) {
    time.Sleep(1200 * time.Millisecond)
    for _,s := range songlist {
        if s == song {
            c <- peerID
            return
        }
    }
    c <- (-peerID)
}

func SearchPeers(PeerList map[int] ([]string), song string, c chan int) int {
    sleeper := make(chan int)
    go Sleeper(sleeper)
    for p := range PeerList {
		go LookForSong(p, PeerList[p], song, c)        
	}
    cntY := 0
    cntN := 0
    numY := 3
    // stop when we have enough peers
    // or when all peers return with Y/N
    for (cntY < numY) && ((cntY + cntN) < len(PeerList)){
        select {
        case i:=<-c:
            if i > 0 {
                cntY++
            	fmt.Printf("Peer %d has song %s!\n", i, song)
            } else {
                cntN++
                fmt.Printf("Peer %d does not have song...\n", -i)
            }
        case <-sleeper:
            fmt.Println("Took too long!")
            close(c)
            close(sleeper)
            return cntY
        }
	}
    return cntY
}

func main() {
	PeerList := make(map[int] []string)
	PeerList[1] = []string{"A", "B", "C"}
	PeerList[2] = []string{"D", "E"}
	PeerList[3] = []string{"A", "G", "H"}
	PeerList[4] = []string{}
	PeerList[5] = []string{"G", "E", "I", "C"}
    c := make(chan int)
	
	// current peer p is looking for some song among its peerlist, and passes a channel
    n := SearchPeers(PeerList, "C", c)
    fmt.Println("Num peers with song...", n)
    
    fmt.Println("Program Finished")
}