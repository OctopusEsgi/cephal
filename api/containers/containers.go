package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ID      string `json:"ID"`
	Image   string `json:"Image"`
	Status  string `json:"Status"`
	Name    string `json:"Name"`
	Ports   string `json:"Ports"`
	Created string `json:"Created"`
}

func getContainers() ([]ContainerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containerInfos []ContainerInfo
	for _, ctr := range containers {
		ports := ""
		for _, p := range ctr.Ports {
			ports += fmt.Sprintf("%d/%s -> %d, ", p.PrivatePort, p.Type, p.PublicPort)
		}

		containerInfos = append(containerInfos, ContainerInfo{
			ID:      ctr.ID[:10],
			Image:   ctr.Image,
			Status:  ctr.Status,
			Name:    ctr.Names[0],
			Ports:   ports,
			Created: time.Unix(ctr.Created, 0).Format(time.RFC3339),
		})
	}

	return containerInfos, nil
}

func ContainersapiHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := getContainers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}
