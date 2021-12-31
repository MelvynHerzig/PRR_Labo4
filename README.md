# PRR_Labo4
Respository du laboratoire 03 pour le cours PRR

# Étudiants
- Forestier Quentin
- Herzig Melvyn

# Introduction
Ce laboratoire a pour but de comprendre les paradigmes de programmation dans une réseau réparti. Pour la réalisation, nous avons dû utiliser le protocole UDP unicast. Nous avons démontré de paradigmes:
* Algorithme ondulatoire 
* Algorithme sondes et échos

Chacune de ces parties effectue la recherche de plus courts chemins dans un réseau __connexe__.

La donnée complète du laboratoire est disponible [ici](https://github.com/MelvynHerzig/PRR_Labo4/blob/main/Labo_4_PRR_donnee.pdf).

# État
## Fonctionne
Toutes les fonctionnalités demandées dans la [donnée](https://github.com/MelvynHerzig/PRR_Labo4/blob/main/Labo_4_PRR_donnee.pdf) ont été implémentées avec succès.

De plus, nous avons pris la liberté d'ajouter un mécanisme d'attente au démarrage des serveurs. Les serveurs n'acceptent aucune demande client tant que tous n'ont pas été démarrés. Si le client effectue une demande alors que les serveurs ne sont pas complétement démarrée, celle-ci est renouvellé automatiquement jusqu'à l'attente du démarrage complet des serveurs. 

## Ne fonctionne pas
En se basant sur la [donnée](https://github.com/MelvynHerzig/PRR_Labo4/blob/main/Labo_4_PRR_donnee.pdf), tout fonctionne.

## Améliorations possibles
Pour l'algorithme ondulatoire, si une demande est en cours d'exécution (tous les serveurs ont déjà reçu le signal de départ), aucune autre demande ne peut être démarrée. Toutefois, si deux demandes initiales sont émises en même temps, le comportement est indéfini. Il faudrait améliorer cet aspect avec un mécanisme de section critique par exemple. 

Le réseau est considéré sans panne, sans erreur et ne change pas au fil du temps. Dans une situation réelle, ces éléments devraient être pris en compte.

# Protocole de communication UDP de démarrage des serveurs
Nous avons mis en place un protocole de démarrage des serveurs. Les serveurs doivent s'attendre avant de démarrer le traitement des demandes.

Durant le protocole de démarrage uniquement, tous les messages sont acquittés. Si aucun acquittement n'est reçu après 1s, le message est re-envoyé. 

## Comment un serveur trouve un autre serveur (adresses et ports)?
Le serveur interroge le fichier _config.json_.

## Qui parle et quand ? 
Au démarrage, un serveur de numéro N ( > 0 ) commence par attendre que le serveur de numéro N - 1 lui envoie le signal qu'il a démarré "OK". 

Dès que le serveur N a reçu le "OK" de N - 1, il envoie à son tour un "OK" à N + 1, seulement si N n'est pas le dernier serveur dans _config.json_

Lorsque le dernier serveur dans _config.json_ (de numéro M) reçoit "OK", cela signifie que tous les serveurs ont été allumés. Il envoie alors à M-1 "GO". M-1 transfert à M-2 et ainsi de suite. recevoir le message go signifie que le serveur peut désormais accepter des demandes clientes.

## Qu'est ce qui se passe quand un message est reçu ? 
Le serveur N vérifie si le contenu du message correspond à ce qu'il devrait recevoir durant sa phase de démarrage et si la source est N +/- 1 (dans un context local, la vérification de la source ne fait pas de sens). Si c'est le cas, il acquitte le message en envoyant un "ACK", sinon le message est abandonné.

## Syntaxe des messages de réplication
### Requête
| Utilité | Syntaxe |
|---|----|
| Signaler la mise en ligne | "OK" CRLF |
| Autorisation de démarrer | "GO" CRLF  |
| Acquitter |"ACK" CRLF  |

## Exemple d'une conversation entre 3 serveurs, Server 3 discute avec Server 2 et Server 4

_Démarrage de Server 2_ 

Server 2 -> 3 : <br>
`OK`\
Server 3 -> 2 :\
`ACK`\
Server 3 -> 4 : <br>
`OK`\
Server 4 -> 3 :\
`ACK`

_Server 4 a reçu OK_ 

Server 4 -> 3 : <br>
`GO`\
Server 3 -> 4 :\
`ACK`\
Server 3 -> 2 : <br>
`GO`\
Server 2 -> 3 :\
`ACK`\

# Première partie: algorithme ondulatoire

## Installation et utilisation

* Cloner le répertoire.
> `$ git clone https://github.com/MelvynHerzig/PRR_Labo4.git`

* Remplir le fichier de configuration _config.json_ à la racine du projet.
  * debug ( booléen, true/false ): Pour lancer les serveurs en mode debug (affiche les messages entrants et sortants)
  * versions ( string, indiquer "wave" ou "probe" ): Pour définir l'algorithme à utiliser "wave" pour ondulatoire ou "probe" pour sondes et échos
  * servers ( ip, port et numéros des voisins [0, Nb serveurs - 1]) Définition du réseau. Au minimum 1 serveur. Si le réseau contient plus d'un serveur, il doit être **connexe**, au quel cas, l'algorithme ne fonctionnerait pas.
```
{
  "debug": true,
  "version": "probe",
  "servers": [
    {
      "ip": "127.0.0.1",
      "port": 3000,
      "neighbors": [1, 2, 3]
    },
    {
      "ip": "127.0.0.1",
      "port": 3001,
      "neighbors": [0, 7]
    },
    {
      "ip": "127.0.0.1",
      "port": 3002,
      "neighbors": [0]
    },
    {
      "ip": "127.0.0.1",
      "port": 3003,
      "neighbors": [0, 4, 5, 6]
    },
    {
      "ip": "127.0.0.1",
      "port": 3004,
      "neighbors": [3]
    },
    {
      "ip": "127.0.0.1",
      "port": 3005,
      "neighbors": [3, 6]
    },
    {
      "ip": "127.0.0.1",
      "port": 3006,
      "neighbors": [3, 5, 7]
    },
    {
      "ip": "127.0.0.1",
      "port": 3007,
      "neighbors": [1, 6]
    }
  ]
}
```
> La configuration précédente est un exemple avec huit serveurs.\
> Le premier serveur dans la liste est le serveur 0 et le dernier, le serveur 7.\
> La configuration précédente représente le graphe suivant:\
> ![arborescence](https://user-images.githubusercontent.com/34660483/147826308-e62ec851-9d5e-4dd8-83b7-b6c104e5f928.png)

* Démarrer le(s) serveur(s). Un argument est nécessaire.
  * Entre 0 et N-1 avec N = nombres de serveurs configurés dans _config.json_

> <u>Depuis le dossier _server_.</u>\
> En admettant le fichier de configuration précédent:\
> `$ go run . 0`\
> `$ go run . 2`\
> `$ go run . 4`\
> `$ go run . 1`\
> `$ go run . 5`\
> `$ go run . 7`\
> `$ go run . 6`\
> `$ go run . 3`
>
> Il est également possible de démarrer les serveurs automatiquement avec les scriptes `server/win_start_servers.bat` pour windows ou `server/lin_start_server.sh` pour linux.\
> \
> L'ordre de démarrage n'est pas important. Durant cette étape, les serveurs s'inter-connectent. En conséquence, tant que tous ne sont pas allumés et connectés, ils n'acceptent que des connexions ayant une adresse IP source appartenant au fichier de configuration. De plus toute demande initiale se voit refusée tant que le réseau n'est pas prêt.\
>\
> Lorsque un serveur a complétement démarré, il affiche les lignes suivantes:\
>![cmd](https://user-images.githubusercontent.com/34660483/147826831-11b7c10d-1aa8-401a-ac24-50f55c371e9e.png)\
> Si tous les serveurs affichent ces messages, le réseau est démarré et prêt à accepter les demandes clientes.

* Démarrer le client pour lancer une demande. Aucun argument

> <u>Depuis le dossier _client_</u>\
>\
> `$ go run .`\
>\
>Le client lance une demande initiale à tous les serveurs.\
>Il peut être lancé avant que tous les serveurs soient prêts. Sa demande est renouvellée automatiquement jusqu'à ce qu'ils soient prêts.

* Lorsque la demande a été prise en compte chaque serveur affiche les plus courts chemins jusqu'aux autres serveurs. 

> Par exemple, le résultat du serveur 0 dans la config de ce _readme_:\
> ![sp](https://user-images.githubusercontent.com/34660483/147827174-952960be-977a-4ef6-bbaf-9a0b4a573055.png)
> 