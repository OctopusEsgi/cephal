package createserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Paramètre nécessaire pour créer un conteneur
type ServerInfo struct {
	Game    string   `json:"game"`    // IMAGE
	Alias   string   `json:"alias"`   // Alias pour debug
	Env     []string `json:"env"`     // ENV
	Ram     int      `json:"ram"`     // RAM
	CPU     int      `json:"cpu"`     // CPU
	PortTCP string   `json:"portTCP"` // Port TCP
	PortUDP string   `json:"portUDP"` // Port UDP
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

	// Ports Docker internes doivent être 6567
	internalPortTCP, err := nat.NewPort("tcp", "6567")
	if err != nil {
		return err
	}
	internalPortUDP, err := nat.NewPort("udp", "6567")
	if err != nil {
		return err
	}

	// Configure les ports à exposer dans le conteneur
	exposedPorts := nat.PortSet{
		internalPortTCP: struct{}{},
		internalPortUDP: struct{}{},
	}

	// Configure HostConfig pour mapper les ports externes vers les ports internes
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			internalPortTCP: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: srvInfo.PortTCP,
				},
			},
			internalPortUDP: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: srvInfo.PortUDP,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Resources: container.Resources{
			Memory:   int64(srvInfo.Ram) * 1024 * 1024, // Convertir MB en bytes
			NanoCPUs: int64(srvInfo.CPU) * 1e9,         // Convertir CPU en nanosecondes
		},
	}

	// Configure le conteneur
	config := &container.Config{
		Image:        srvInfo.Game,
		Env:          srvInfo.Env,
		ExposedPorts: exposedPorts,
		Hostname:     fmt.Sprintf("%s-%s", srvInfo.Game, srvInfo.Alias),
	}

	// Créer le client Docker
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Créer le conteneur
	ctn, err := cli.ContainerCreate(
		context.Background(),
		config,     // Config du conteneur
		hostConfig, // Config de l'hôte
		nil,        // NetworkingConfig
		nil,        // Platform -- Pas utile pour l'instant ??
		fmt.Sprintf("%s-%s", srvInfo.Game, srvInfo.Alias), // Nom du conteneur
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
