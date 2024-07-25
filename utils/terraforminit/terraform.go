package terraforminit

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func remoteRun(user string, addr string, privateKeyPath string, cmd string) (string, error) {
	// Lire la clé privée à partir du chemin spécifié
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	// Configuration SSH
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connexion SSH
	fmt.Println("Connecting to server:", addr)
	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, "22"), config)
	if err != nil {
		return "", fmt.Errorf("unable to connect: %v", err)
	}
	defer client.Close()

	// Création d'une session
	fmt.Println("Creating session")
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("unable to create session: %v", err)
	}
	defer session.Close()

	// Exécution de la commande
	fmt.Println("Running command:", cmd)
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("command execution failed: %v", err)
	}

	return b.String(), nil
}

func ApplyTerraform() (string, error) {
	output, err := remoteRun("octopus", "192.168.69.1", "/home/octopus/.ssh/id_rsa", "cd terraform/ && terraform apply -auto-approve")
	if err != nil {
		return "", err
	}
	return output, nil
}

func DestroyTerraform() (string, error) {
	output, err := remoteRun("octopus", "192.168.69.1", "/home/octopus/.ssh/id_rsa", "cd terraform/ && terraform destroy -auto-approve")
	if err != nil {
		return "", err
	}
	return output, nil
}

func InitTerraform() (string, error) {
	output, err := remoteRun("octopus", "192.168.69.1", "/home/octopus/.ssh/id_rsa", "cd terraform/ && terraform init")
	if err != nil {
		return "", err
	}
	return output, nil
}
