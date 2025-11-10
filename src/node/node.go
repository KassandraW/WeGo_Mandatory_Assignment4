package main

import (
	proto "Mandatory4/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	//"time"

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
	lamport_clock    int64
	cs_access        bool
	state            NodeState
	node_connections []proto.NodeClient
}

func critical_section() {
	fmt.Println("Doing something in the critical section høhø")
}

func loopOfLife() {
	// do some stuff that calculates if we want to access CS or not.
	for {
		//time.Sleep(time.Millisecond * 500)
	}

}

func main() {
	// set up log info
	filepath := "../grpc/Log_info"
	Log_File, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("could not open log file client: %v", err)
	}
	if err := os.Truncate(filepath, 0); // clear the log file on each run
	err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	defer Log_File.Close()

	log.SetOutput(Log_File)
	log.SetFlags(log.Lshortfile)

	//get server port from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("to add a Node pleas enter adress info of 1 node from Adress_file")
	log.Println("to add a Node pleas enter adress info of 1 node from Adress_file")
	adressLine, err := reader.ReadString('\n')
	adressLine = strings.TrimSuffix(adressLine, "\n")
	//fmt.Println("this is the adress line :" + adressLine)

	// Set server up
	server := &Node{}
	server.state = StateReleased
	server.lamport_clock = 0
	server.cs_access = false

	if err != nil {
		fmt.Print("read input not avalable")
		log.Fatalf("scanner failed")
	}
	// repeat for every node in a loop, adding them to the node_connections list
	for i := 0; i < 3; i++ {
	}
	fmt.Print("loooping!")

	conn, err := grpc.NewClient("localhost:5052", grpc.WithTransportCredentials(insecure.NewCredentials())) //connects to server. Insecure.newcredentials is used to skip TLS encryption for simplification
	if err != nil {
		log.Fatalf("conection start Not working")
	}

	server.node_connections = append(server.node_connections, proto.NewNodeClient(conn))

	go loopOfLife()

	server.start_server()
}

func (s *Node) Request(ctx context.Context, timestamp *proto.TimeStamp) (*proto.Empty, error) {
	s.state = StateWanted

	return &proto.Empty{}, nil
}

func (s *Node) start_server() {
	grpcServer := grpc.NewServer()              //Creates a new gRPC server instance
	listener, err := net.Listen("tcp", ":5050") //Listens on TCP port 5050 using net.Listen
	log.Println("Node server is up and running ")
	fmt.Println("Node server is up and running ")
	if err != nil {
		log.Fatalf("failed to start server")
	}

	proto.RegisterNodeServer(grpcServer, s) //registers the server implementation with gRPC

	err = grpcServer.Serve(listener) // starts serving incoming requests

	if err != nil {
		log.Fatalf("Did not work")
	}
}
