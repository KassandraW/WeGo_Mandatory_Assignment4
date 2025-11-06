package main

import (
	proto "Mandatory4/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"bufio"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	node_connections			[]proto.NodeClient
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
	// set up log info 
	filepath := "../grpc/Log_info"
	Log_File, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("could not open log file client: %v", err)
	}
	defer Log_File.Close()


	// Set server up
	server := &Node{}
	server.state = StateReleased
	server.lamport_clock = 0
	server.cs_access = false
	// make reader avalable 
	reader := bufio.NewReader(os.Stdin)
	var line, _ = reader.ReadString('\n')

	// repeat for every node in a loop, adding them to the node_connections list
	for i := 0 ; i < 3; i++{
		fmt.Println("to add a Node pleas enter adress info of 1 node from Adress_file")
		log.Println("to add a Node pleas enter adress info of 1 node from Adress_file")

		conn, err := grpc.NewClient(line, grpc.WithTransportCredentials(insecure.NewCredentials( ))) //connects to server. Insecure.newcredentials is used to skip TLS encryption for simplification
	if err != nil {
		log.Fatalf("conection start Not working")
	}
	server.node_connections = append(server.node_connections, proto.NewNodeClient(conn))
	} 
	

	go loopOfLife();




	server.start_server()
}


func (s *Node) Request(ctx context.Context, timestamp *proto.TimeStamp ) (*proto.Empty, error) {
	s.state = StateWanted;

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

