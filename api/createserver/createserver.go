package createserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Paramètre nécessaire pour créer un conteneur
type ServerInfo struct {
	Game string `json:"game"` // IMAGE
	Env  string `json:"env"`  // ENV
	Ram  int    `json:"ram"`  // RAM
	CPU  int    `json:"cpu"`  // CPU
}

type ContainerInfo struct {
	ID      string `json:"ID"`
	Image   string `json:"Image"`
	Status  string `json:"Status"`
	Name    string `json:"Name"`
	Ports   string `json:"Ports"`
	Created string `json:"Created"`
}

func createServer(srvInfo ServerInfo) error {
	fmt.Println("reçu:", srvInfo)

	// Configure les ports TCP et UDP
	portTCP := "5012"
	portUDP := "5012"
	newPortTCP, err := nat.NewPort("tcp", portTCP)
	if err != nil {
		return nil
	}
	newPortUDP, err := nat.NewPort("udp", portUDP)
	if err != nil {
		return nil
	}

	// Configure HostConfig
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			newPortTCP: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: portTCP,
				},
			},
			newPortUDP: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: portUDP,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{},
		},
		Resources: container.Resources{
			Memory:   int64(srvInfo.Ram) * 1024 * 1024, // Convertir MB en bytes
			NanoCPUs: int64(srvInfo.CPU) * 1e9,         // Convertir CPU en nanosecondes
		},
	}

	// Configure NetworkingConfig
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	gatewayConfig := &network.EndpointSettings{
		Gateway: "gatewayname",
	}
	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

	// Configure les ports à exposer
	exposedPorts := map[nat.Port]struct{}{
		newPortTCP: {},
		newPortUDP: {},
	}

	// Configure le conteneur
	config := &container.Config{
		Image:        srvInfo.Game,
		Env:          []string{srvInfo.Env},
		ExposedPorts: exposedPorts,
		Hostname:     fmt.Sprintf("%s-hostnameexample", srvInfo.Game),
	}

	// Créer le client Docker
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Créer le conteneur
	ctn, err := cli.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		networkConfig,
		nil,
		"TEST",
	)
	if err != nil {
		return err
	}

	// Démarrer le conteneur
	if err := cli.ContainerStart(context.Background(), ctn.ID, container.StartOptions{}); err != nil {
		return err
	}

	log.Printf("Container %s is created and started", ctn.ID)

	return nil
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

	err = createServer(*jsonServer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
