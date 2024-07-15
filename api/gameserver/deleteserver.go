package gameserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// StopOptions représente les options pour arrêter un container
type StopOptions struct {
	Signal  string `json:",omitempty"`
	Timeout *int   `json:",omitempty"`
}

type RemoveOptions struct {
	RemoveVolumes bool `json:",omitempty"`
	RemoveLinks   bool `json:",omitempty"`
	Force         bool `json:",omitempty"`
}

// ServerInfo représente les informations du serveur à arrêter
type ServerDelInfo struct {
	ContainerID   string        `json:"container_id"`
	StopOptions   StopOptions   `json:"stopoptions"`
	RemoveOptions RemoveOptions `json:"removeoptions"`
}

// deleteGameServer arrête un container de jeu donné
func deleteGameServer(infos ServerDelInfo) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	stopOptions := container.StopOptions{
		Signal:  infos.StopOptions.Signal,
		Timeout: infos.StopOptions.Timeout,
	}

	err = cli.ContainerStop(context.Background(), infos.ContainerID, stopOptions)
	if err != nil {
		return err
	}

	removeOptions := container.RemoveOptions{
		RemoveVolumes: infos.RemoveOptions.RemoveVolumes,
		RemoveLinks:   infos.RemoveOptions.RemoveLinks,
		Force:         infos.RemoveOptions.Force,
	}
	err = cli.ContainerRemove(context.Background(), infos.ContainerID, removeOptions)
	if err != nil {
		return err
	}

	return nil
}

// DeleteServerAPIHandler gère les requêtes pour arrêter un serveur de jeu
func DeleteServerAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var serverInfo ServerDelInfo
	err := json.NewDecoder(r.Body).Decode(&serverInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = deleteGameServer(serverInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success"}
	json.NewEncoder(w).Encode(response)
}
