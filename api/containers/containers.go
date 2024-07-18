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

// ContainerInfo est la structure contenant les informations sur le conteneur
type ContainerInfo struct {
	ID      string `json:"id"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Name    string `json:"name"`
	Ports   string `json:"ports"`
	Created string `json:"created"`
}

// getContainers récupère la liste des conteneurs
func GetContainers() ([]ContainerInfo, error) {
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

// getContainerByID récupère les informations d'un conteneur par son ID
func getContainerByID(id string) (*ContainerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	containerJSON, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	ports := ""
	for port, bindings := range containerJSON.NetworkSettings.Ports {
		for _, binding := range bindings {
			ports += fmt.Sprintf("%s:%s -> %s, ", binding.HostIP, binding.HostPort, port)
		}
	}

	createdTime, err := time.Parse(time.RFC3339Nano, containerJSON.Created)
	if err != nil {
		return nil, err
	}

	containerInfo := &ContainerInfo{
		ID:      containerJSON.ID[:10],
		Image:   containerJSON.Config.Image,
		Status:  containerJSON.State.Status,
		Name:    containerJSON.Name,
		Ports:   ports,
		Created: createdTime.Format(time.RFC3339),
	}

	return containerInfo, nil
}

// ContainersapiHandler gère les requêtes API pour les conteneurs
func ContainersapiHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		// Si un ID est fourni, récupérer les informations du conteneur correspondant
		container, err := getContainerByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(container)
		return
	}

	// Sinon, récupérer la liste de tous les conteneurs
	containers, err := GetContainers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}
