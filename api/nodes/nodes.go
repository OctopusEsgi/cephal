package nodes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type NodeInfo struct {
	ID       string `json:"ID"`
	Hostname string `json:"Hostname"`
	Status   string `json:"Status"`
	Role     string `json:"Role"`
	State    string `json:"State"`
}

func getNodes() ([]NodeInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeInfos []NodeInfo
	for _, node := range nodes {
		nodeInfos = append(nodeInfos, NodeInfo{
			ID:       node.ID[:10],
			Hostname: node.Description.Hostname,
			Status:   string(node.Status.State),
			Role:     string(node.Spec.Role),
			State:    string(node.Spec.Availability),
		})
	}

	return nodeInfos, nil
}

func NodesAPIHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := getNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}
