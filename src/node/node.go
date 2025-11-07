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
    lamport_clock       int64
    cs_access           bool
    state               NodeState 
    node_connections            []proto.NodeClient
}


func critical_section(){
    fmt.Println("Doing something in the critical section høhø")
}   


func loopOfLife(){
    for{
        time.Sleep(time.Millisecond * 500) // is there a smarter way to do this? this seems wasteful.
        // do some stuff that calculates if we want to access CS or not. 
    }
}


func main() { 
    // set up log info 
    filepath := "../grpc/Log_info"
    Log_info, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("could not open log file client: %v", err)}
    if err := os.Truncate(filepath, 0); // clear the log file on each run
    err != nil {
        log.Printf("Failed to terminate File: %v", err)
    }
    defer Log_info.Close()
    
    log.SetOutput(Log_info)
    log.SetFlags(log.Lshortfile)



    
    // Set server up
    server := &Node{}
    fmt.Println("to add a Node pleas enter adress info of 1 node from Adress_file")
    log.Println("to add a Node pleas enter adress info of 1 node from Adress_file")
    server.state = StateReleased
    server.lamport_clock = 0
    server.cs_access = false
    
    
    // repeat for every node in a loop, adding them to the node_connections list 
    //for i := 0; i<3; i++{}
    // get server port from user 
    reader := bufio.NewReader(os.Stdin)
    adressLine, err := reader.ReadString('\n')
    adressLine = strings.TrimSuffix(adressLine,"\n") // tims input so that it is usable to start a new node up
    if err != nil {
        fmt.Print("read input not avalable") 
        log.Fatalf("scanner failed")
    }
        conn, err := grpc.NewClient(adressLine, grpc.WithTransportCredentials(insecure.NewCredentials( ))) //connects to server. Insecure.newcredentials is used to skip TLS encryption for simplification
    if err != nil {
        log.Fatalf("conection start Not working")
    }
    server.node_connections = append(server.node_connections, proto.NewNodeClient(conn))
    fmt.Println("Node running on port " + adressLine)

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
        log.Fatalf("Did not work(start server)")
    }
    fmt.Println("Node servers is up and running")
    log.Println("Node servers is up and running")

    proto.RegisterNodeServer(grpcServer, s) //registers the server implementation with gRPC

    err = grpcServer.Serve(listener) // starts serving incoming requests

    if err != nil {
        log.Fatalf("Did not work (incoming requests)")
 
	}

}

