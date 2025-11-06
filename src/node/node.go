package main

import (
	proto "Mandatory4/grpc"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// change states
type NodeState int
const (
	StateReleased NodeState = iota
	StateHeld
	StateWanted
)

type Node struct { 
	proto.UnimplementedNodeServer
	lamport_clock 		int64
	cs_access			bool
	state				NodeState 		


}


func critical_section(){
	fmt.Println("Doing something in the critical section høhø")
}	


func loopOfLife(){
	for{
		// do some stuff that calculates if we want to access CS or not. 
	}
}


func main() { 
	server := &Node{}
	server.state = StateReleased
	server.lamport_clock = 0;
	server.cs_access = false;
	go loopOfLife();
	server.start_server()
}


func (s *Node) Request(ctx context.Context, timestamp *proto.TimeStamp ) (*proto.Empty, error) {
	return &proto.Empty{}, nil
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

