package main

import (
	"cephal/api/containers"
	"cephal/api/gameserver"
	"cephal/api/nodes"
	"cephal/api/services"
	"cephal/imagesinit"

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
	// -- TEST --
	fmt.Println(imagesinit.GetImagesList())
	// --fe
	//
	//
	// FRONT
	http.Handle("/front/", http.StripPrefix("/front/", http.FileServer(http.Dir("front"))))
	// API
	http.HandleFunc("/api/containers", containers.ContainersapiHandler)
	// http.HandleFunc("/api/container/", )
	http.HandleFunc("/api/nodes", nodes.NodesAPIHandler)
	http.HandleFunc("/api/services", services.ServicesAPIHandler)
	http.HandleFunc("/api/createserver", gameserver.CreateServerAPIHandler)
	http.HandleFunc("/api/deleteserver", gameserver.DeleteServerAPIHandler)
	http.HandleFunc("/", frontHandler)
	fmt.Println("Lancement du serveur sur le port 8080")
	// LANCEMENT SRV
	log.Fatal(http.ListenAndServe(":8080", nil))
}
