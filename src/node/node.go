package main

import (
	proto "Mandatory4/grpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Node struct { 
	proto.UnimplementedNodeServer
	lamport_clock 		int64
	cs_access			bool

}

func critical_section(){
	log.Printf("Doing something in the critical section høhø")
}	


func main() { //initializes server
	server := &Node{}
	server.lamport_clock = 0;
	server.cs_access = false;
	server.start_server() //starts the gRPC server
}

func (s *Node) start_server() {
	grpcServer := grpc.NewServer() //Creates a new gRPC server instance
	listener, err := net.Listen("tcp", ":5050") //Listens on TCP port 5050 using net.Listen
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterNodeServer(grpcServer, s) //registers the server implementation with gRPC

	err = grpcServer.Serve(listener) // starts serving incoming requests

	if err != nil {
		log.Fatalf("Did not work")
	}
}

