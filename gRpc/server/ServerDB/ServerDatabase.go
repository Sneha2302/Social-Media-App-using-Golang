package serverDB

import (
	pb "social_media_app-golang/gRpc/protobuff"
	"social_media_app-golang/storage"
	// "bytes"
	"context"
	"encoding/json"
	// "flag"
	"fmt"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	// "net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DB struct {
	RaftDir  string
	RaftBind string
	Inmem    bool

	mu        sync.Mutex
	UsersInfo map[string]storage.Twitter_User
	Raft      *raft.Raft // The consensus mechanism

	Logger *log.Logger
}
type command struct {
	Op   string
	Name string
	Info storage.Twitter_User
}

func (db *DB) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserReply, error) {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	var uName string = in.Uname
	//var temp storage.Twitter_User = db.UsersInfo[uName]
	temp, _ := db.Get(uName)
	user := storage.UserToProto(temp)

	return &pb.GetUserReply{Userinfo: user}, nil
}

func (db *DB) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.BoolReply, error) {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	uName := in.UName    //storage
	password1 := in.Pwd1 //storage
	password2 := in.Pwd2 //storage
	if password1 != password2 {
		return &pb.BoolReply{T: false}, nil
	}
	if uName == "" || password1 == "" {
		return &pb.BoolReply{T: false}, nil
	}
	StoreUser := storage.Twitter_User{uName, password1, storage.Tweets{}, []string{uName}}
	if _, ok := db.Get(uName); ok {
		return &pb.BoolReply{T: false}, nil
	}

	db.Set(uName, StoreUser)
	return &pb.BoolReply{T: true}, nil
}

func (db *DB) Tweet(ctx context.Context, in *pb.TweetRequest) (*pb.BoolReply, error) {
	// db.mu.Lock()
	// defer db.mu.Unlock()
	uName := in.UName
	usr := storage.ProtoToUser(in.ProtoUser)
	if uName != usr.UName {
		return &pb.BoolReply{T: false}, nil
	}
	// if _, ok := db.UsersInfo[uName]; ok != true {
	if _, ok := db.Get(uName); ok != true {
		return &pb.BoolReply{T: false}, nil
	}

	db.Set(uName, usr)
	return &pb.BoolReply{T: true}, nil
}

func (db *DB) UserExists(ctx context.Context, in *pb.UserExistsRequest) (*pb.BoolReply, error) {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	uName := in.UName  //storage
	password := in.Pwd //storage
	if uName == "" || password == "" {
		return &pb.BoolReply{T: false}, nil
	}
	// Check Whether User in usersInfo
	//user, exist := db.UsersInfo[uName]
	user, exist := db.Get(uName)
	if exist && user.Pwd == password {
		return &pb.BoolReply{T: true}, nil
	}
	return &pb.BoolReply{T: false}, nil
}

func (db *DB) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.BoolReply, error) {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	uName := in.UName         //form
	otherName := in.Othername //form
	if twitter_user, ok := db.Get(uName); ok {
		if storage.Contains(twitter_user.Follow, otherName) {
			return &pb.BoolReply{T: false}, nil
		}
		twitter_user.Follow = append(twitter_user.Follow, otherName)
		//db.UsersInfo[uName] = twitter_user
		db.Set(uName, twitter_user)
		return &pb.BoolReply{T: true}, nil
	}
	return &pb.BoolReply{T: false}, nil
}

func (db *DB) UnFollow(ctx context.Context, in *pb.FollowRequest) (*pb.BoolReply, error) {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	uName := in.UName         //storage
	otherName := in.Othername //storage
	if twitter_user, ok := db.Get(uName); ok {
		if !storage.Contains(twitter_user.Follow, otherName) {
			return &pb.BoolReply{T: false}, nil
		}
		twitter_user.Follow = storage.Deletes(twitter_user.Follow, otherName)
		//db.UsersInfo[uName] = twitter_user
		db.Set(uName, twitter_user)
		return &pb.BoolReply{T: true}, nil
	}
	return &pb.BoolReply{T: false}, nil
}

func GetTweets(arr storage.Tweets) []string {
	var ret []string
	for _, tweet := range arr {

		tmp := tweet.User + ": " + tweet.Contents
		ret = append(ret, tmp)
	}
	return ret
}

func (db *DB) GetUserInfo(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	//db.mu.Lock()
	//defer db.mu.Unlock()
	uName := in.Username //from memory
	user, _ := db.UsersInfo[uName]

	UserName := user.UName
	Following := user.Follow

	var UnFollowed []string
	var Posts storage.Tweets
	// Get all Posts information
	for name, userInfo := range db.UsersInfo {
		if storage.Contains(Following, name) {
			for _, post := range userInfo.Post {
				Posts = append(Posts, post)
			}
		} else {
			UnFollowed = append(UnFollowed, name)
		}
	}
	userPosts := storage.GetUserInfo(Posts)
	// Remove the user itself from following list (just not shown in screen but in memory)
	Following = storage.Deletes(Following, uName)
	log.Printf("This User is %s", UserName)
	log.Printf("User is following %s", Following)
	log.Printf("User is not following: %s", UnFollowed)
	log.Printf("User's post %s", userPosts)
	var tweet = &pb.Page{UName: UserName, Unfollow: UnFollowed, Follow: Following, Posts: userPosts}
	return &pb.GetUserInfoReply{UName: tweet.UName, Unfollow: tweet.Unfollow, Follow: tweet.Follow, Posts: tweet.Posts}, nil

}

func (db *DB) IsLeader(ctx context.Context, in *pb.IsLeaderRequest) (*pb.BoolReply, error) {
	log.Println("------>", db.Raft.State())
	var temp1 bool = db.Raft.State() == raft.Leader
	log.Println("raft leader --->", temp1)
	return &pb.BoolReply{T: temp1}, nil
}

type fsm DB

func (s *DB) Join(ctx context.Context, in *pb.JoinRequest) (*pb.BoolReply, error) {
	nodeID := in.NodeID
	addr := in.RemoteAddr

	log.Println("------> Join.....")
	s.Logger.Printf("Join request received from remote node %s at %s", nodeID, addr)

	configFuture := s.Raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.Logger.Printf("Unabe to get raft configuration: %v", err)
		return &pb.BoolReply{T: false}, err
	}
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				s.Logger.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return &pb.BoolReply{T: false}, nil
			}
			future := s.Raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return &pb.BoolReply{T: false}, fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}
	f := s.Raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		fmt.Println(f.Error())
		return &pb.BoolReply{T: false}, f.Error()
	}
	s.Logger.Printf("node %s at %s joined successfully", nodeID, addr)
	return &pb.BoolReply{T: true}, nil
}

func (s *DB) Open(enableSingle bool, localID string) error {
	// Setting up Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	// Raft communication being setup
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	// Creating raft snapshot store
	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, 2, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}
	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	if s.Inmem {
		logStore = raft.NewInmemStore()
		stableStore = raft.NewInmemStore()
	} else {
		boltDB, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
		if err != nil {
			return fmt.Errorf("new bolt store: %s", err)
		}
		logStore = boltDB
		stableStore = boltDB
	}
	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	s.Raft = ra
	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}
	return nil
}

type fsmSnapshot struct {
	UsersInfo map[string]storage.Twitter_User
	mp        map[string]string
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	// Clone the map.
	o1 := make(map[string]storage.Twitter_User)
	//o2 := make(map[string]string)
	for k, v := range f.UsersInfo {
		o1[k] = v
	}

	return &fsmSnapshot{UsersInfo: o1}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	return nil
}

// Get returns the value for the given key.
func (s *DB) Get(name string) (storage.Twitter_User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.UsersInfo[name]
	return info, ok
}

// Set sets the value for the given key.
func (s *DB) Set(name string, info storage.Twitter_User) error {
	if s.Raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	c := &command{
		Op:   "set",
		Name: name,
		Info: info,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := s.Raft.Apply(b, 2*time.Second)
	return f.Error()
}

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("Unable to marshal the command: %s", err.Error()))
	}
	switch c.Op {
	case "set":
		return f.applySet(c.Name, c.Info)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

func (f *fsm) applySet(name string, info storage.Twitter_User) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	fmt.Println(name)
	fmt.Println(info)
	f.UsersInfo[name] = info
	return nil
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (f *fsmSnapshot) Release() {}
