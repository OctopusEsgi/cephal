package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ServiceInfo struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Mode     string `json:"Mode"`
	Replicas uint64 `json:"Replicas"`
	VIP      string `json:"VIP"`
	Labels   string `json:"Labels"`
}

func getServices() ([]ServiceInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceInfos []ServiceInfo
	for _, service := range services {
		serviceInfos = append(serviceInfos, ServiceInfo{
			ID:       service.ID[:10],
			Name:     service.Spec.Name,
			Mode:     fmt.Sprint(service.Spec.Mode.Replicated),
			Replicas: *service.Spec.Mode.Replicated.Replicas,
			VIP:      fmt.Sprint(service.Spec.EndpointSpec.Ports),
			Labels:   fmt.Sprint(service.Spec.Labels),
		})
	}

	return serviceInfos, nil
}

func ServicesAPIHandler(w http.ResponseWriter, r *http.Request) {
	services, err := getServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}
