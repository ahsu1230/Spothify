package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	//"time"
	"../impl"
	"../rpc"
)

var portnum *int = flag.Int("port", 0, "port # to listen on.  Monitor nodes default to 9009.")
var bridgeHP *string = flag.String("bridge", "", "bridge port # to connect to.")

func main() {

	flag.Parse()
	if *portnum == 0 || *bridgeHP == "" {
		log.Fatal("Need to specify port number and/or bridgeHP!")
	}
	if (flag.NArg() < 1) {
		fmt.Printf("arguments: [username]\n")
		log.Fatal("Insufficient arguments for server")
	}
	username := flag.Arg(0)
	
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	_, listenport, _ := net.SplitHostPort(l.Addr().String())
	log.Println("Server listening at ", listenport)
	*portnum, _ = strconv.Atoi(listenport)

	fmt.Printf("Arguments: [username: %s] [portnum:%d] [bridgeport:%s]\n", username, *portnum, *bridgeHP)
	ps := peernode.NewPeerServer(username, *portnum, *bridgeHP)

	prpc := peernoderpc.NewPeerRPC(ps)
	rpc.Register(prpc)
	log.Println("Server starting HTTP")
	rpc.HandleHTTP()
	//go runServer(ns)
	log.Println("Server serving HTTP\n")
	http.Serve(l, nil)
	log.Println("Server done")
}
/*
func runServer(ns *serverimpl.Server) {
	time.Sleep(time.Duration(2) * time.Second)
	for !ms.GameData.GameEnded {
		// Dial to All Servers and check their connections
		err := ms.DialAllServers()
		if err != nil {
			log.Println("ERROR DIALING!", err.Error())
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	log.Println("Game Ended!")
}
*/