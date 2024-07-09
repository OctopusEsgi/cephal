package createserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServerInfo struct {
	Game string `json:"game"`
	Env  string `json:"env"`
	Ram  int    `json:"ram"`
	CPU  int    `json:"cpu"`
}

func createServer(srvInfo ServerInfo) {
	fmt.Println("reçu:", srvInfo)
}

func CreateServerAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	jsonServer := &ServerInfo{}
	err := json.NewDecoder(r.Body).Decode(jsonServer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createServer(*jsonServer)
	w.WriteHeader(http.StatusCreated)
}
