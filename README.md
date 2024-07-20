# Cephal - Projet Octopus

Voici le module Cephal du projet annuel de l'ESGI de Reims (équipe Théo, Adrien, Mael et Yanis) le projet "Octopus".

Cephal permet de réception des requêtes API Rest pour manipuler un environement docker / docker swarm.

Cephal peut être installer sur n'importe quel environement basé sur docker et même docker swarm ;
Features :
- Création de serveur de jeux volatile via API
- Gestion automatique des ports / ressources machine
- Destrution automatique des containers
- Reconstruction de son environement en cas de panne ou reset de docker

Cephal peut fonctionner par lui même grâce à une page d'administration mais il est recommandé d'utiliser l'API.


DISCLAIMER:
Pour l'instant seul le jeu [Mindustry](https://mindustrygame.github.io/) est disponible dans cephal par defaut, il est cepedent possible d'ajouter d'autre images de jeux via le fichier de configuration principal (/etc/cephal/cephal.yml). J'ai fait en sorte que l'ajout de futur soit simplifié ! Cela reste cependant très peu testé et très experimental /!\


## Structure du projet

-   **api/** : Contient les différents gestionnaires d'API pour les conteneurs, les nœuds, les services et la création de serveurs.
    -   **containers/** : 
    -   **gameserver/** : 
    -   **nodes/** : 
    -   **ports/** : 
    -   **services/** : 
-   **config/** : Contient les config par défault de Cephal
-   **utils/** : Gestionnaire de config / builder d'images / manager de port


## Exemple d'utilisation de l'API

**Créer un serveur de jeux**

```bash
curl --location 'http://10.255.0.13:8080/api/createserver' \
--header 'Content-Type: application/json' \
--data '{
    "game": "mindustryesgi",
    "alias": "srv1ef01",
    "env": [
        "VERSION=v146",
        "MAP=Tendrils",
        "MODE=sandbox",
        "PLAYERSLIMIT=2"
    ]
}'
```

En cas de succès l'API renvoie :

```json
{
    "Id": "ID_LONG_CONTAINER",
    "Warnings": []
}
```

**Supprimer un serveur de jeu**

```bash
curl --location --request DELETE 'http://localhost:8080/api/deleteserver' \
--header 'Content-Type: application/json' \
--data '{
    "container_id": "ton_giga_mega_id_de_container",
    "stopoptions": {
        "Signal": "SIGTERM",
        "Timeout": 10
    },
    "removeoptions": {
        "RemoveVolumes": true,
        "RemoveLinks": false,
        "Force": true
    }
}'
```
- `container_id` (string) : ID du conteneur à supprimer.
- `stopoptions` (objet) : Options de l'ârret du conteneur.

    - `Signal` (string) : Le signal à envoyer pour arrêter le conteneur. SIGTERM ou SIGKILL par exemple.
    - `Timeout` (int) : Le temps d'attente maximum en secondes pour que le conteneur s'arrête. 0 = stop direct (ungracefully stop), -1 = infini (aka il attend que l'app soit bien stoppé)

- `removeoptions` (objet) : Options pour la suppression du conteneur.

    - `RemoveVolumes` (booléen) : Si true, les volumes associés au conteneur seront également supprimés.
    - `RemoveLinks` (booléen) : TOUJOURS SUR FALSE *pour l'instant*
    - `Force` (booléen) : Si true, la suppression du conteneur sera forcée même s'il est en cours d'exécution.

En cas de succès l'API renvoie :

```json
{
    "status": "success"
}
```

**Récupérer la liste des containers présent sur la machine**

```bash
curl --location 'http://localhost:8080/api/containers'
```

Pour un seul container avec son id:

```bash
curl --location 'http://localhost:8080/api/containers?id=CAFE001122AA'
```

**Récupérer la liste des nœuds Swarm**

```bash
curl --location 'http://localhost:8080/api/nodes'
```

*Recupération par ID par encore implémenté !*

**Récupérer les ports disponibles** <- ??

```bash
curl --location 'http://localhost:8080/api/getusedports'
```

> Ce handle devrait disparaitre 🦖 -Yanis

## Installation

OUTDATED

