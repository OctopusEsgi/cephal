# Utilise l'image officielle de base OpenJDK 11
FROM openjdk:11-jre-slim

# Définit le répertoire de travail
WORKDIR /app

# Installe wget
RUN apt-get update && apt-get install -y wget && rm -rf /var/lib/apt/lists/*

# Expose le port par défaut utilisé par le serveur Mindustry
EXPOSE 6567/tcp
EXPOSE 6567/udp

# Commande pour télécharger le fichier jar et lancer le serveur
ENTRYPOINT ["sh", "-c", "wget https://github.com/Anuken/Mindustry/releases/download/${VERSION}/server-release.jar -O server.jar && java -jar server.jar host \"${MAP}\" ${MODE},playerlimit ${PLAYERSLIMIT}"]