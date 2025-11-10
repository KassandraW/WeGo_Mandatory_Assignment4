package main

import (
	proto "Mandatory4/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// change states
type NodeState int

const (
	Released NodeState = iota
	Held
	Wanted
)

type Node struct {
	proto.UnimplementedNodeServer
	lamport_clock    int64
	cs_access        bool
	state            NodeState
	node_port        string
	node_connections []client
	l                sync.Mutex
	reply_count      int
	request_queue    []*proto.Msg_Info
}

type client struct {
	port       string
	nodeclient proto.NodeClient
}

// representation of the critical section
func critical_section() {
	fmt.Println("Doing something in the critical section høhø")
	log.Println("Doing something in the critical section høhø")
}

// increments lamport clock
func (s *Node) inc_clock() {
	s.l.Lock()
	s.lamport_clock++
	s.l.Unlock()
}

// starts the process of accessing the critical section
func (s *Node) request_access() {
	s.inc_clock()
	s.state = Wanted
	s.reply_count = 0

	//send request to all nodes
	fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Sending a request to everyone")
	log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Sending a request to everyone")
	for _, c := range s.node_connections {
		_, err := c.nodeclient.Request(context.Background(), &proto.Msg_Info{Ts: s.lamport_clock, Port: s.node_port[len(s.node_port)-4:]})
		check(err, "Could not send request to port "+c.port)
	}

	//wait for replies from all nodes
	for s.reply_count < len(s.node_connections) {
		time.Sleep(time.Millisecond * 100)
	}

	//enter critical section
	s.inc_clock()
	s.state = Held
	fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Entering critical section")
	log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Entering critical section")
	critical_section()

	//exit critical section after some time
	time.Sleep(time.Second * 15)
	s.inc_clock()
	fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Exiting critical section")
	log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Exiting critical section")
	s.inc_clock()
	s.state = Released

	//reply to all queued requests
	s.inc_clock()
	for _, req := range s.request_queue {
		for _, c := range s.node_connections {
			if c.port == req.Port {
				fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Replying to queued request from " + req.Port)
				log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Replying to queued request from " + req.Port)
				_, err := c.nodeclient.Reply(context.Background(), &proto.Msg_Info{Ts: s.lamport_clock, Port: s.node_port})
				check(err, "Could not send reply to port "+c.port)
			}
		}
	}
}

// currently just a loop that waits for user input to request access to critical section
func (s *Node) loopOfLife() {
	for {
		var input string
		fmt.Scan(&input)

		if input == "cs" {
			if s.state != Released {
				fmt.Println("You either already have access, or is currently requesting access, dumbass")
				continue
			}
			go s.request_access()
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// helper function for error checks
func check(e error, msg string) {
	if e != nil {
		log.Fatalf(msg+": %v", e)
	}
}

func main() {
	// Log
	filepath := "../files/log.log"
	Log_File, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	check(err, "could not open log file client")
	defer Log_File.Close()

	log.SetOutput(Log_File)
	log.SetFlags(0)

	// Server Setup
	server := &Node{}
	server.state = Released
	server.lamport_clock = 0
	server.cs_access = false
	fmt.Println("Please input the nodes port:")
	fmt.Scanln(&server.node_port)
	log.SetPrefix("[" + server.node_port + "] ")

	// Read ports from Nodes.txt
	file, err := os.Open("../files/nodes.txt")
	check(err, "could not open nodes.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("scanner failed")
		}

		port := scanner.Text()
		last4_port := port[len(port)-4:] //get last 4 digits of the port to identify the node
		if last4_port == server.node_port {
			continue //skip connecting to self
		}
		conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials())) //connects to server. Insecure.newcredentials is used to skip TLS encryption for simplification
		check(err, "failed to connect to port "+port)
		server.node_connections = append(server.node_connections, client{port: last4_port, nodeclient: proto.NewNodeClient(conn)})
	}
	//start the loop of life in a separate goroutine
	go server.loopOfLife()

	server.start_server()

}

func (s *Node) Request(ctx context.Context, req_info *proto.Msg_Info) (*proto.Empty, error) {
	//Update lamport clock
	s.lamport_clock = max(s.lamport_clock, req_info.Ts) + 1

	//log
	fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Request from " + req_info.Port)
	log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Request from " + req_info.Port)

	//We only reply when we arent currently in "HELD", or "WANTED"  and the request has a lower timestamp
	if s.state == Held || (s.state == Wanted && s.lamport_clock > req_info.Ts) {
		s.inc_clock() //local event i guess
		fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Putting the request from " + req_info.Port + " in the queue")
		log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Putting the request from " + req_info.Port + " in the queue")
		//defer reply
		s.request_queue = append(s.request_queue, req_info)

	} else {
		//send reply
		s.inc_clock()
		fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Replying to request from " + req_info.Port)
		log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Replying to request from " + req_info.Port)
		for _, c := range s.node_connections {
			if c.port == req_info.Port {
				_, err := c.nodeclient.Reply(context.Background(), &proto.Msg_Info{Ts: s.lamport_clock, Port: s.node_port})
				check(err, "Could not send reply to port "+c.port)
				return &proto.Empty{}, nil
			}
		}

	}
	return &proto.Empty{}, nil
}

func (s *Node) Reply(ctx context.Context, name *proto.Msg_Info) (*proto.Empty, error) {
	s.inc_clock()
	s.lamport_clock = max(s.lamport_clock, name.Ts) + 1
	s.reply_count++
	fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Reply from " + name.Port)
	log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Reply from " + name.Port)
	return &proto.Empty{}, nil
}

func (s *Node) start_server() {
	grpcServer := grpc.NewServer()                      //Creates a new gRPC server instance
	listener, err := net.Listen("tcp", ":"+s.node_port) //Listens on TCP port with the given port using net.Listen
	if err != nil {
		log.Fatalf("failed to start node on port %s: %v", s.node_port, err)
	}
	s.inc_clock()
	log.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Node is up and running")
	fmt.Println("T" + strconv.Itoa(int(s.lamport_clock)) + " : Node is up and running")

	proto.RegisterNodeServer(grpcServer, s) //registers the server implementation with gRPC

	err = grpcServer.Serve(listener) // starts serving incoming requests

	if err != nil {
		log.Fatalf("Did not work")
	}

}
