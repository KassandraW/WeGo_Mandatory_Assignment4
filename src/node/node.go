package main

import (
	proto "Mandatory4/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

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
	node_port        string
	node_connections []proto.NodeClient
}

func critical_section() {
	fmt.Println("Doing something in the critical section høhø")
}

func loopOfLife() {
	// do some stuff that calculates if we want to access CS or not.
	for {
		time.Sleep(time.Millisecond * 500)
	}

}

func check(e error, msg string) {
	if e != nil {
		log.Fatalf(msg+": %v", e)
	}
}

func main() {
	// Log
	filepath := "../files/Log_info"
	Log_File, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	check(err, "could not open log file client")

	err = os.Truncate(filepath, 0) //clear the log file on each run
	check(err, "Failed to truncate")
	defer Log_File.Close()

	log.SetOutput(Log_File)
	log.SetFlags(log.Lshortfile)

	// Server Setup
	server := &Node{}
	server.state = StateReleased
	server.lamport_clock = 0
	server.cs_access = false
	fmt.Println("Please input the nodes port:")
	fmt.Scanln(&server.node_port)

	// Read ports from Nodes.txt
	file, err := os.Open("../files/Nodes.txt")
	check(err, "could not open Nodes.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("scanner failed")
		}

		port := scanner.Text()
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials())) //connects to server. Insecure.newcredentials is used to skip TLS encryption for simplification
		check(err, "failed to connect to port "+port)
		server.node_connections = append(server.node_connections, proto.NewNodeClient(conn))
		fmt.Println("Connected to port " + conn.Target())
		log.Println("Connected to port " + conn.Target())
	}

	go loopOfLife() //starts the loop of life in a separate goroutine

	server.start_server()

}

func (s *Node) Request(ctx context.Context, timestamp *proto.TimeStamp) (*proto.Empty, error) {
	s.state = StateWanted

	return &proto.Empty{}, nil
}

func (s *Node) start_server() {
	grpcServer := grpc.NewServer()                      //Creates a new gRPC server instance
	listener, err := net.Listen("tcp", ":"+s.node_port) //Listens on TCP port with the given port using net.Listen
	if err != nil {
		log.Fatalf("failed to start node on port %s: %v", s.node_port, err)
	}
	log.Println("The node on port " + s.node_port + " is nowup and running")
	fmt.Println("The node on port " + s.node_port + " is nowup and running")

	proto.RegisterNodeServer(grpcServer, s) //registers the server implementation with gRPC

	err = grpcServer.Serve(listener) // starts serving incoming requests

	if err != nil {
		log.Fatalf("Did not work")
	}

}
