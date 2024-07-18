package apiPort

import (
	"cephal/utils/portmanager"
	"encoding/json"
	"net/http"
)

type PortInfo struct {
	TCP []string `json:"tcp"`
	UDP []string `json:"udp"`
}

// GetUsedPorts récupère les ports TCP et UDP utilisés
func GetUsedPorts() (PortInfo, error) {
	tcpPorts, udpPorts := portmanager.FindUsedPort()
	return PortInfo{
		TCP: tcpPorts,
		UDP: udpPorts,
	}, nil
}

// PortsAPIHandler est un gestionnaire HTTP pour obtenir les ports disponibles
func PortsAPIHandler(w http.ResponseWriter, r *http.Request) {
	ports, err := GetUsedPorts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ports)
}
