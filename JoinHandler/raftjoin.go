package handle

import (
	"social_media_app-golang/grpc/client"

	"encoding/json"
	"log"
	"net/http"
)

func Join(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering join ------>>>>>>>>>>>")
	m := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("-->aaaaaaa------>")
	if len(m) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("-->bbbbbbbbb------>")
	remoteAddr, ok := m["addr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("-->cccccc------>")
	nodeID, ok := m["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("-->Finally joining----->")
	client.RpcJoin(nodeID, remoteAddr)
}
