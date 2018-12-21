package main

import (
	"bytes"
	pb "social_media_app-golang/gRpc/protobuff"
	"social_media_app-golang/gRpc/server/serverDB"
	"social_media_app-golang/storage"

	// "context"
	"encoding/json"
	"flag"
	"fmt"
	// "github.com/hashicorp/raft"
	// "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// "io"
	"log"
	"net"
	"net/http"
	"os"
	// "path/filepath"
	// "sync"
	// "time"
)

var port = ":9091"
var storageDir string
var rpcPort string
var raftPort string
var nodeName string
var isLeader bool

func join(nodeID, raftAddr string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	log.Println("-------> Enter join")
	// http.Post(request_url, "application/x-www-form-urlencoded", body)
	resp, err := http.Post(fmt.Sprintf("http://:9090/join"), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	log.Println("======")
	defer resp.Body.Close()

	return nil
}

func init() {
	flag.StringVar(&storageDir, "storageDir", "/tmp/dir1", "Set the storage directory")
	flag.StringVar(&rpcPort, "rpcPort", "9090", "Set Rpc bind address")
	flag.StringVar(&raftPort, "raftPort", "9091", "Set Raft bind address")
	flag.StringVar(&nodeName, "nodeName", "node0", "Set the name of server")
	flag.BoolVar(&isLeader, "isLeader", false, "node is leader or not")
	flag.Usage = func() {
		fmt.Println("Usage: go run server.go [options] <data-path>")
		flag.PrintDefaults()
	}
}

func New() *serverDB.DB {
	fmt.Println(storageDir)
	fmt.Println(rpcPort)
	fmt.Println(raftPort)
	fmt.Println(nodeName)
	WebDB := &serverDB.DB{
		Inmem:     false,
		RaftDir:   storageDir,
		RaftBind:  raftPort,
		Logger:    log.New(os.Stderr, "[store2] ", log.LstdFlags),
		UsersInfo: make(map[string]storage.Twitter_User),
	}
	return WebDB
}
func connectRpc() (*grpc.Server, net.Listener) {
	port = rpcPort
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	return s, lis
}
func main() {
	flag.Parse()
	s, lis := connectRpc()
	WebDB := New()
	log.Println("------->", isLeader)
	if err := WebDB.Open(isLeader, nodeName); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	// time.Sleep(5 * time.Second)
	log.Println("------->", isLeader)
	if !isLeader {
		join(nodeName, raftPort)
	}

	pb.RegisterWebServer(s, WebDB)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
