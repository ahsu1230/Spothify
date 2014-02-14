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
	if *portnum == *monitorport {
		log.Fatal("Portnum cannot match monitor's portnum!")
	}
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	_, listenport, _ := net.SplitHostPort(l.Addr().String())
	log.Println("Server listening at ", listenport)
	*portnum, _ = strconv.Atoi(listenport)

	fmt.Printf("Arguments: [portnum:%d] [monitorport:%d]\n", *portnum, *monitorport)
	bs := bridgenode.NewBridgeServer(*portnum, *monitorport)

	brpc := bridgenoderpc.NewBridgeRPC(bs)
	rpc.Register(brpc)
	log.Println("Server starting HTTP")
	rpc.HandleHTTP()
	//go runServer(ns)
	log.Println("Server serving HTTP")
	http.Serve(l, nil)
	log.Println("Server done")
}