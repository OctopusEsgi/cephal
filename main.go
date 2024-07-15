package main

import (
	"cephal/api/containers"
	"cephal/api/gameserver"
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
	// Serve static files (CSS, JS, images, etc.) from the "front" directory
	http.Handle("/front/", http.StripPrefix("/front/", http.FileServer(http.Dir("front"))))

	http.HandleFunc("/api/containers", containers.ContainersapiHandler)
	// http.HandleFunc("/api/container/", )
	http.HandleFunc("/api/nodes", nodes.NodesAPIHandler)
	http.HandleFunc("/api/services", services.ServicesAPIHandler)
	http.HandleFunc("/api/createserver", gameserver.CreateServerAPIHandler)
	http.HandleFunc("/", frontHandler)
	fmt.Println("Lancement du serveur sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
