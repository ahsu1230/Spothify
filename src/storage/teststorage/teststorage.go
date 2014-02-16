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
var monitorport *int = flag.Int("monitor", 0, "Monitor node to connect to - defaults to 9009.")

func main() {

	flag.Parse()
	if *portnum == 0 {
		log.Fatal("Need a portnum!")
	}
	if *monitorport == 0 {
		log.Println("Monitor defaulting to 9009")
		*monitorport = 9009
	}
	id := flag.Arg(0)
	storageID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Println("Invalid Storage ID!");
	}
	
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	_, listenport, _ := net.SplitHostPort(l.Addr().String())
	log.Println("Server listening at ", listenport)
	*portnum, _ = strconv.Atoi(listenport)

	log.Println("Arguments: [portnum:%d] [monitorport:%d] [id:%d]\n", *portnum, *monitorport, storageID)
	ss := storagenode.NewStorageServer(*portnum, *monitorport, uint32(storageID))

	srpc := storagenoderpc.NewStorageRPC(ss)
	rpc.Register(srpc)
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