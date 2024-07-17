package gameserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	initconf "cephal/utils/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Paramètre nécessaire pour créer un conteneur
type ServerInfo struct {
	Game     string   `json:"game"`
	Alias    string   `json:"alias"`
	Env      []string `json:"env"`
	PortsTCP []string `json:"portsTCP"` // Liste des ports TCP externes
	PortsUDP []string `json:"portsUDP"` // Liste des ports UDP externes
}

type ContainerInfo struct {
	ID      string `json:"ID"`
	Image   string `json:"Image"`
	Status  string `json:"Status"`
	Name    string `json:"Name"`
	Ports   string `json:"Ports"`
	Created string `json:"Created"`
}

func getGameImageConfig(game string, confCephalIMG []initconf.GameImage) (*initconf.GameImage, error) {
	for _, img := range confCephalIMG {
		if img.Nom == game {
			return &img, nil
		}
	}
	return nil, fmt.Errorf("image de jeux non définie: %s", game)
}

func createServer(srvInfo ServerInfo, confCephal *initconf.ConfigCephal) (*container.CreateResponse, error) {

	gameConfig, err := getGameImageConfig(srvInfo.Game, confCephal.GameImages)
	if err != nil {
		return nil, err
	}

	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	// Lier les ports TCP
	if len(srvInfo.PortsTCP) != len(gameConfig.Ports.TCP) {
		return nil, fmt.Errorf("mismatch in the number of TCP ports")
	}
	for i, externalPort := range srvInfo.PortsTCP {
		internalPort, err := nat.NewPort("tcp", gameConfig.Ports.TCP[i])
		if err != nil {
			return nil, err
		}
		exposedPorts[internalPort] = struct{}{}
		portBindings[internalPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: externalPort,
			},
		}
	}

	// Lier les ports UDP
	if len(srvInfo.PortsUDP) != len(gameConfig.Ports.UDP) {
		return nil, fmt.Errorf("mismatch in the number of UDP ports")
	}
	for i, externalPort := range srvInfo.PortsUDP {
		internalPort, err := nat.NewPort("udp", gameConfig.Ports.UDP[i])
		if err != nil {
			return nil, err
		}
		exposedPorts[internalPort] = struct{}{}
		portBindings[internalPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: externalPort,
			},
		}
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Resources: container.Resources{
			Memory:   int64(gameConfig.Spec.RAM) * 1024 * 1024, // Convertir MB en bytes
			NanoCPUs: int64(gameConfig.Spec.Core) * 1e9,        // Convertir CPU en nanosecondes
		},
	}

	config := &container.Config{
		Image:        srvInfo.Game,
		Env:          srvInfo.Env,
		ExposedPorts: exposedPorts,
		Hostname:     fmt.Sprintf("%s-%s", srvInfo.Game, srvInfo.Alias),
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	ctn, err := cli.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		nil,
		nil,
		fmt.Sprintf("%s-%s", srvInfo.Game, srvInfo.Alias),
	)
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(context.Background(), ctn.ID, container.StartOptions{}); err != nil {
		return nil, err
	}

	log.Printf("Container %s est créé et démarré", ctn.ID)

	return &ctn, nil
}

func CreateServerAPIHandler(confCephal *initconf.ConfigCephal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		response, err := createServer(*jsonServer, confCephal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERREUR: %s", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
