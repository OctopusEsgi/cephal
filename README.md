# Cephal - Projet Octopus

Voici le module Cephal du projet annuel de l'ESGI de Reims (équipe Théo, Adrien, Mael et Yanis) le projet "Octopus".

Cephal permet de réception des requêtes API Rest pour manipuler un envirronement docker swarm.


## Structure du projet

-   **api/** : Contient les différents gestionnaires d'API pour les conteneurs, les nœuds, les services et la création de serveurs.
-   **front/** : Contient les fichiers statiques de l'interface Administrateur (de test).
-   **main.go** : Point d'entrée principal de l'application.

## Exemple d'utilisation de l'API

En cURL:
```c
curl --location 'http://docker.home:8080/api/createserver' \
--header 'Content-Type: application/json' \
--data '{
    "game": "mindustryesgi",
    "alias": "srv100",
    "env": ["VERSION=v146", "MAP=Tendrils", "MODE=sandbox", "PLAYERSLIMIT=2"],
    "ram": 512,
    "cpu": 2,
    "portsTCP": ["25010"],
    "portsUDP": ["25010"]
}'
```

en
 
## Installation

Vous pouvez installer Cephal sur n'importe quel environnement Docker (Docker en Standalone ou Docker en mode Swarm).

L'application sera disponible sur `http://*:8080`.

