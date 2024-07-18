package main

import (
	"cephal/api/containers"
	"cephal/api/gameserver"
	"cephal/api/nodes"
	apiPort "cephal/api/ports"
	"cephal/api/services"
	initconf "cephal/utils/config"
	"cephal/utils/imagesinit"
	"os"
	"path/filepath"

	"fmt"
	"html/template"
	"log"
	"net/http"
)

func frontHandler(config *initconf.ConfigCephal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { //On fait ça pour que le handler recoit que
		tmplPath := config.Server.RootDirSRV
		t, err := template.ParseFiles(fmt.Sprintf("%s/front.html", tmplPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	//CONFIG LOAD
	if len(os.Args) < 2 {
		log.Fatal("Vous devez spécifié un chemin en tant qu'argument de cephal.")
	}
	configPath := os.Args[1]
	configCephal, err := initconf.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Printf("Configuration chargée [%s]", configPath)
	}

	// -- TEST --
	images := []imagesinit.ImagePath{
		{ImageName: "mindustryesgi:latest", Dockerfile: "/home/octopus/godev/cephal/config/builder/mindustryesgi/Dockerfile"},
	}

	err = imagesinit.EnsureImagesList(images)
	if err != nil {
		log.Fatal(err)
	}
	// --fe
	//
	//
	// FRONT
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(configCephal.Server.RootDirSRV, "static")))))
	// API
	http.HandleFunc("/api/containers", containers.ContainersapiHandler)
	http.HandleFunc("/api/nodes", nodes.NodesAPIHandler)
	http.HandleFunc("/api/services", services.ServicesAPIHandler)
	http.HandleFunc("/api/createserver", gameserver.CreateServerAPIHandler(configCephal))
	http.HandleFunc("/api/deleteserver", gameserver.DeleteServerAPIHandler)
	http.HandleFunc("/api/getfreeports", apiPort.PortsAPIHandler)
	http.HandleFunc("/", frontHandler(configCephal))
	log.Printf("Lancement du serveur sur le port %d", configCephal.Server.Port)
	// --

	// LANCEMENT SRV
	addr := fmt.Sprintf(":%d", configCephal.Server.Port)
	if configCephal.Server.TLS.Enabled {
		log.Fatal(http.ListenAndServeTLS(addr, configCephal.Server.TLS.CertFile, configCephal.Server.TLS.KeyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}
