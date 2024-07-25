package apiterraform

import (
	"cephal/utils/terraforminit"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ServerInfo struct {
	Action string `json:"action"`
}

func traiteMoiLePostDeTerraFormPOurAdrienQUiLeVoulaisAuDernierMomentCarSinonIlNaimaitPasLeProjet(action string) (string, error) {
	switch action {
	case "init":
		return terraforminit.InitTerraform()
	case "apply":
		return terraforminit.ApplyTerraform()
	case "destroy":
		return terraforminit.DestroyTerraform()
	default:
		return "", fmt.Errorf("action inconnue: %s", action)
	}
}

func NodesTerraform(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Reception de [IP:%s], data : %s", r.RemoteAddr, *jsonServer)
	response, err := traiteMoiLePostDeTerraFormPOurAdrienQUiLeVoulaisAuDernierMomentCarSinonIlNaimaitPasLeProjet(jsonServer.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERREUR: %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
