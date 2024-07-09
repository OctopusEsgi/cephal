package main

import (
	"cephal/api/containers"
	"cephal/api/createserver"
	"cephal/api/nodes"
	"cephal/api/services"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func frontHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("front", "front.html")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/containers", containers.ContainersapiHandler)
	http.HandleFunc("/api/nodes", nodes.NodesAPIHandler)
	http.HandleFunc("/api/services", services.ServicesAPIHandler)
	http.HandleFunc("/api/createserver", createserver.CreateServerAPIHandler)
	http.HandleFunc("/", frontHandler)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
