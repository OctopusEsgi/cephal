package main

import (
	"cephal/api/containers"
	"cephal/api/gameserver"
	"cephal/api/nodes"
	"cephal/api/services"
	"cephal/imagesinit"
	"os"

	"fmt"
	"html/template"
	"log"
	"net/http"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"tls"`
	Port       int    `yaml:"port"`
	RootDirSRV string `yaml:"root_directory"`
}

type GameImage struct {
	Nom   string `yaml:"nom"`
	Tag   string `yaml:"tag"`
	Ports struct {
		TCP []string `yaml:"tcp"`
		UDP []string `yaml:"udp"`
	} `yaml:"ports"`
	Spec struct {
		Core int `yaml:"core"`
		RAM  int `yaml:"ram"`
	} `yaml:"spec"`
}

type ConfigCephal struct {
	Server     ServerConfig `yaml:"server"`
	GameImages []GameImage  `yaml:"gameimages"`
}

func loadConfig(filename string) (*ConfigCephal, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config ConfigCephal
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func frontHandler(config *ConfigCephal) http.HandlerFunc {
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
	configCephal, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Printf("Configuration chargée [%s]", configPath)
	}

	// -- TEST --
	images := []imagesinit.ImagePath{
		{ImageName: "mindustryesgi:latest", Dockerfile: "path/to/your/Dockerfile1"},
	}

	err = imagesinit.EnsureImagesList(images)
	if err != nil {
		log.Fatal(err)
	}
	// --fe
	//
	//
	// FRONT
	http.Handle("/front/", http.StripPrefix("/front/", http.FileServer(http.Dir("front"))))
	// API
	http.HandleFunc("/api/containers", containers.ContainersapiHandler)
	http.HandleFunc("/api/nodes", nodes.NodesAPIHandler)
	http.HandleFunc("/api/services", services.ServicesAPIHandler)
	http.HandleFunc("/api/createserver", gameserver.CreateServerAPIHandler)
	http.HandleFunc("/api/deleteserver", gameserver.DeleteServerAPIHandler)
	http.HandleFunc("/", frontHandler(configCephal))
	log.Printf("Lancement du serveur sur le port %d", configCephal.Server.Port)
	// LANCEMENT SRV
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configCephal.Server.Port), nil))
}
