# Social_media_app-golang
Distributed Systems: Final Project

Team Members:
Shradha Ahuja (sa4741)
Sneha Munden (sm7352)

Architecture Diagram:

![Screenshot](Arch.jpg)

Setup:
Clone the repository and change the name of the cloned folder to twitterPt2 
Now go to path ./grpc/server to run the Raft Backend Servers. 
Run the following commands in 3 seperate terminal windows to run 3 nodes as backend server
      go run server.go -storageDir /tmp/nodeA -nodeName nodeA -rpcPort :9093 -raftPort :12000 -isLeader=true
      go run server.go -storageDir /tmp/nodeB -nodeName nodeB -rpcPort :9094 -raftPort :13000 -isLeader=false
      go run server.go -storageDir /tmp/nodeC -nodeName nodeC -rpcPort :9095 -raftPort :14000 -isLeader=false


Now in a separate terminal window run the client server with the following command
      go run main.go


  Go to - http://localhost:9090/login







