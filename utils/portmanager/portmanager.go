package portmanager

import (
	"cephal/api/containers"
	"cephal/utils/config"

	"encoding/json"
	"fmt"
	"strings"
)

func splitPorts(portString string) []string {
	var ports []string
	portEntries := strings.Split(portString, ",")
	for _, entry := range portEntries {
		port := strings.TrimSpace(entry)
		ports = append(ports, port)
	}
	return ports
}

// extractPort nettoie "->"
func extractPort(portString string) string {
	parts := strings.Split(portString, "->")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(parts[0])
}

// FindAndPrintPort trouve et imprime les ports, divisés en TCP et UDP
func FindUsedPort() ([]string, []string) {
	body, err := containers.GetContainers()
	if err != nil {
		fmt.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return nil, nil
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	var containers []containers.ContainerInfo

	err = json.Unmarshal(jsonData, &containers)
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return nil, nil
	}

	var tcpPorts []string
	var udpPorts []string
	for _, container := range containers {
		portList := splitPorts(container.Ports)
		for _, port := range portList {
			cleanPort := extractPort(port)
			if strings.Contains(port, "tcp") {
				tcpPorts = append(tcpPorts, cleanPort)
			} else if strings.Contains(port, "udp") {
				udpPorts = append(udpPorts, cleanPort)
			}
		}
	}
	return tcpPorts, udpPorts
}

func isPortUsed(port string, usedPorts []string) bool {
	for _, usedPort := range usedPorts {
		if port == usedPort {
			return true
		}
	}
	return false
}

func AssignPorts(nbtcp, nbudp int, confCephal *config.ConfigCephal) (tcpPorts []string, udpPorts []string, err error) {
	usedTcpPorts, usedUdpPorts := FindUsedPort()
	portRange := confCephal.Global.Portrange

	tcpPorts = make([]string, 0, nbtcp)
	udpPorts = make([]string, 0, nbudp)

	for port := portRange.Min; port <= portRange.Max && len(tcpPorts) < nbtcp; port++ {
		portStr := fmt.Sprintf("%d", port)
		if !isPortUsed(portStr, usedTcpPorts) {
			tcpPorts = append(tcpPorts, portStr)
		}
	}

	for port := portRange.Min; port <= portRange.Max && len(udpPorts) < nbudp; port++ {
		portStr := fmt.Sprintf("%d", port)
		if !isPortUsed(portStr, usedUdpPorts) {
			udpPorts = append(udpPorts, portStr)
		}
	}
	if len(tcpPorts) < nbtcp {
		err = fmt.Errorf("impossible d'assigner %d ports TCP, seuls %d disponibles", nbtcp, len(tcpPorts))
	}

	if len(udpPorts) < nbudp {
		err = fmt.Errorf("impossible d'assigner %d ports UDP, seuls %d disponibles", nbudp, len(udpPorts))
	}

	return tcpPorts, udpPorts, err
}
