# PRR_Labo4
Respository du laboratoire 04 pour le cours PRR

# Étudiants
- Forestier Quentin
- Herzig Melvyn

# Introduction
Ce laboratoire a pour but de comprendre les paradigmes de programmation dans une réseau réparti. Pour la réalisation, nous avons dû utiliser le protocole UDP unicast. Nous avons démontré de paradigmes:
* Algorithme ondulatoire 
* Algorithme sondes et échos

Chacune de ces parties effectue la recherche de plus courts chemins dans un réseau en échangeant des topologies (matrices d'adjacence).

La donnée complète du laboratoire est disponible [ici](https://github.com/MelvynHerzig/PRR_Labo4/blob/main/Labo_4_PRR_donnee.pdf).

# État
## Fonctionne
Toutes les fonctionnalités demandées dans la [donnée](https://github.com/MelvynHerzig/PRR_Labo4/blob/main/Labo_4_PRR_donnee.pdf) ont été implémentées avec succès.

De plus, nous avons pris la liberté d'ajouter un mécanisme d'attente au démarrage des serveurs. Les serveurs n'acceptent aucune demande client tant que tous n'ont pas été démarrés. Si le client effectue une demande alors que les serveurs ne sont pas complétement démarrée, celle-ci est renouvellé automatiquement jusqu'à l'attente du démarrage complet des serveurs. 

## Ne fonctionne pas
En se basant sur la [donnée](https://github.com/MelvynHerzig/PRR_Labo4/blob/main/Labo_4_PRR_donnee.pdf), tout fonctionne.

## Améliorations possibles
Pour l'algorithme ondulatoire, si une demande est en cours d'exécution (tous les serveurs ont déjà reçu le signal de départ), aucune autre demande ne peut être démarrée. Toutefois, si deux demandes initiales sont émises en même temps, le comportement est indéfini. 

Il faudrait améliorer cet aspect, par exemple avec un mécanisme de section critique.

Le réseau est considéré sans panne, sans erreur et ne change pas au fil du temps. Dans une situation réelle, ces éléments devraient être pris en compte.

# Protocole de communication UDP de démarrage des serveurs (serveur - serveur)
Nous avons mis en place un protocole de démarrage des serveurs. Les serveurs doivent s'attendre avant de démarrer le traitement des demandes.

Durant le protocole de démarrage uniquement, tous les messages sont acquittés. Si aucun acquittement n'est reçu après 1 seconde, le message est re-envoyé. 

## Comment un serveur trouve un autre serveur (adresses et ports)?
Le serveur interroge le fichier _config.json_.

## Qui parle et quand ? 
Au démarrage, un serveur de numéro N ( > 0 ) commence par attendre que le serveur de numéro N - 1 lui envoie le signal qu'il a démarré "OK". 

Dès que le serveur N a reçu le "OK" de N - 1, il envoie à son tour un "OK" à N + 1, seulement si N n'est pas le dernier serveur dans _config.json_

Lorsque le dernier serveur dans _config.json_ (de numéro M) reçoit "OK", cela signifie que tous les serveurs ont été allumés. Il envoie alors à M-1 "GO". M-1 transfert à M-2 et ainsi de suite. Recevoir le message "GO" signifie que le serveur peut désormais accepter des demandes clientes.

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
`ACK`

## Tests
Les tests ont été réalisés avec la fichier de configuration de ce readme.
### Attente des serveurs non démarrés
__Description__\
Les serveurs attendent que tous les serveurs soient démarrés avant d'accepter les demandes. Les serveur démarrés doivent afficher:\
![l1](https://user-images.githubusercontent.com/34660483/147829116-e89c04c5-c6b0-499c-aca3-84846f42691c.png)\
__Résultat__\
<span style="color:green">Succès</span>

### Les serveurs s'attendent avant de démarrer
__Description__\
Les serveurs sont démarrées dans un ordre aléatoire. Après un instant, ils doivent afficher:\
![cmd](https://user-images.githubusercontent.com/34660483/147826831-11b7c10d-1aa8-401a-ac24-50f55c371e9e.png)\
__Résultat__\
<span style="color:green">Succès</span>

### Un message indésiré ne perturbe pas le processus de démarrage
__Description__\
Lors du démarrage le client de l'algorithme ondulatoire est démarré, il envoie "START" à tout les serveurs. Les serveurs ignorent le message et démarrent correctement\
__Résultat__\
<span style="color:green">Succès</span>

# Protocole de communication UDP de lancement d'une demande (client - serveur)
Ce protocole est utilisé par le client pour démarrer une recherche des plus courts chemins

## Comment le client trouve le(s) serveur(s) (adresses et ports)?
Le client interroge le fichier _config.json_.

## Qui parle et quand ? 
Le client envoie le message "EXEC" au(x) serveur(s). Si l'algorithme configuré est ondulatoire, le client contacte tous les serveurs. Sinon si l'algorithme est sondes et échos, il contacte le serveur reçu en paramètre.

Le client attend de recevoir un "ACK" durant 1 seconde, sinon la demande est renouvelée.

## Qu'est ce qui se passe quand le message est reçu par le serveur ? 
Si le serveur est en cours de démarrage, le message est ignoré, sinon il acquitte le message et transmet le message au processus d'exécution de l'algorithme.

> Le message peut être ignoré par le processus d'exécution de l'algorithme si une demande est déjà en cours de traitement.

## Syntaxe des messages de réplication
### Requête
| Utilité | Syntaxe |
|---|----|
| Effectuer une demande | "EXEC" CRLF |
| Acquitter |"ACK" CRLF  |

## Exemple d'une conversation entre un client et 3 serveurs (algorithme ondulatoire)

Client -> Server 1:\
`EXEC`\
Server 1 -> Client:\
`ACK`\
Client -> Server 2:\
`EXEC`

_Pas de ack après une seconde_

Client -> Server 2:\
`EXEC`\
Server 2 -> Client:\
`ACK`\
Client -> Server 3:\
`EXEC`\
Server 3 -> Client:\
`ACK`

## Exemple d'une conversation entre un client et un serveur (algorithme sondes et échos)

Client -> Server 5:\
`EXEC`

_Pas de ack après une seconde_

Client -> Server 5:\
`EXEC`\
Server 5 -> Client:\
`ACK`

## Tests
Les tests ont été réalisés avec la fichier de configuration de ce readme.
### Réitération de la demande
__Description__\
Si un serveur n'acquitte pas la demande, cette dernière est réitérée\
__Résultat__\
<span style="color:green">Succès</span>

### Contacte tous les serveurs avec l'algorithme ondulatoire
__Description__\
Si l'algorithme configuré est ondulatoire, le client ne termine pas avant d'avoir pû transmettre sa demande à tous les serveurs.\
__Résultat__\
<span style="color:green">Succès</span>

### Contacte du serveur désigné avec l'algorithme sondes et échos
__Description__\
Si l'algorithme configuré est sondes et échos, le client ne termine pas avant d'avoir pû transmettre sa demande au serveur désigné en argument.\
__Résultat__\
<span style="color:green">Succès</span>

# Première partie: algorithme ondulatoire

## Installation et utilisation

* Cloner le répertoire.
> `$ git clone https://github.com/MelvynHerzig/PRR_Labo4.git`

* Remplir le fichier de configuration _config.json_ à la racine du projet.
  * debug ( booléen, true/false ): Pour lancer les serveurs en mode debug (affiche les messages entrants et sortants)
  * version ( string, indiquer "wave" ou "probe" ): Définit l'algorithme, utiliser "wave" pour l'algorithme ondulatoire.
  * servers ( ip, port et numéros des voisins [0, Nb serveurs - 1]) Définition du réseau. Au minimum 1 serveur. Si le réseau contient plus d'un serveur, il doit être **connexe** pour l'algorithme ondulatoire. Dans le cas contraire, l'algorithme ne fonctionnera pas.
```
{
  "debug": true,
  "version": "wave",
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
> Sous Windows, il est également possible de démarrer les serveurs automatiquement avec le scripte `server/win_start_servers.bat`.\
> \
> L'ordre de démarrage n'est pas important. Durant cette étape, les serveurs s'inter-connectent. En conséquence, tant que tous ne sont pas allumés et connectés, ils n'acceptent que des connexions ayant une adresse IP source appartenant au fichier de configuration. De plus toute demande initiale sont refusées tant que le réseau n'est pas prêt.\
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

## Protocole UDP des ondes
Les ondes échangées antres les serveurs ont le format suivant:\
\<matrice d'adjacence> <numéro du serveur source> <numéro de l'onde> \<activité de la source>

Par exemple, si le serveur 2 émet sa troisième onde, il n'est plus actif et possède la matrice d'adjacence:
```
[false, true, false]
[true, false, true]
[false, true, false]
```
Le message sérialisé sera "0-1-0_1-0-1_0-1-0 2 3 0"
> Les booléens sont traduits en 0 (=false) et 1 (=vrai)

> Les lignes de la matrice sont séparées par des '_'

> Les colonnes de la matrice sont séparées par des '-'

> Les messages ne doivent pas excéder 1024 bytes. Il n'est pas conseillé de faire un réseau avec plus de 20 noeuds

## Tests
Les tests ont été réalisés avec la fichier de configuration de ce readme.

Pour ces tests nous n'allons pas analyser la sortie de toutes les consoles. Nous allons nous concentrer sur des situations spécifiques 

Pour les réceptions, les heures affichées sont les heures de traitement et non pas les heures effectives de réception.

### Les ondes sont straitées séquentiellement
__Description__\
Les ondes doivent être traitées les unes à la suite des autres, sans chevauchement, indépendament de la vitesse de chaque noeud.

Ce test vise à tester la mise en cache des ondes reçues mais qui ne doivent pas être encore traîtées à cette itération.

Nous allons analyser la sortie de la console 6  
`1  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 1 1`\
`2  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 1 1`\
`3  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 1 1`\
`4  DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 1 1`\
`5  DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 3 1 1`\
`6  DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 1 1`\
`7  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 2 1`\
`8  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 2 1`\
`9  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 2 1`\
`10 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 3 2 1`\
`11 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 2 1`\
`12 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 2 1`\
`13 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 3 1`\
`14 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 3 1`\
`15 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 3 1`\
`16 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 3 1`\
`17 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 3 1`\
`18 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 3 3 0`\
`19 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 4 0`\
`20 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 4 0`\
`21 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 4 0`\
`22 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 4 0`\
`[Résultat final omis]`

> Pour rappel, l'avant dernier numéro représente le numéro de l'onde.

Comme nous pouvons le voir, par exemple, aux lignes 4-5-6 ou 10-11-12 les ondes ne se chevauchent pas. Le serveur attend d'avoir reçu les ondes correspondantes à l'onde locale en cours avant de passer à l'onde suivante.

Pour effectuer ce test, nous avons volontairement ralenti le server 5. Ce qui a laissé aux serveurs 3 et 7 l'opportunité de potentiellement voler la place du serveur 5 durant l'onde en cours.

__Résultat__\
<span style="color:green">Succès</span>

### Les ondes sont traitées séquentiellement
__Description__\
Si un serveur se désigne comme inactif, il ne reçoit et n'émet plus de messages à la prochaine onde.

Nous allons analyser la sortie de la console 6  
`1  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 1 1`\
`2  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 1 1`\
`3  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 1 1`\
`4  DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 1 1`\
`5  DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 3 1 1`\
`6  DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 1 1`\
`7  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 2 1`\
`8  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 2 1`\
`9  DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 2 1`\
`10 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 3 2 1`\
`11 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 2 1`\
`12 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 2 1`\
`13 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 3 1`\
`14 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 3 1`\
`15 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 3 1`\
`16 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 3 1`\
`17 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 3 1`\
`18 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 3 3 0`\
`19 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 4 0`\
`20 DEBUG >> Dec 31 16:07:52 SENDED) [matrice omise] 6 4 0`\
`21 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 7 4 0`\
`22 DEBUG >> Dec 31 16:07:52 RECEIVED) [matrice omise] 5 4 0`
`[Résultat final omis]`

> Pour rappel, l'antépénultième et le dernier nombres sont respectivement le serveur source et son état (actif ou non)

Comme nous pouvons le voir à la ligne 18, le serveur 3 s'annonce innactif. Le serveur 6 ne lui envoie plus de message (lignes 19 et 20) et il n'attend pas de recevoir ses messages pour continuer son exécution (lignes 21 et 22).

__Résultat__\
<span style="color:green">Succès</span>

### Fusion correcte des topologies
__Description__\
Le serveur modifie correctement sa topologie avec les informations reçues.

Nous allons analyser 4 lignes de la console 7\
`1 [ligne omise]`\
`2 DEBUG >> Dec 31 16:07:52 SENDED)`\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-1-0-0-0-0-1-0 `[fin omises]`\
`3 DEBUG >> Dec 31 16:07:52 RECEIVED)`\
0-0-0-0-0-0-0-0_\
1-0-0-0-0-0-0-1_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0 `[fin omises]`\
`4 DEBUG >> Dec 31 16:07:52 RECEIVED)`\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-1-0-1-0-1_\
0-0-0-0-0-0-0-0 `[fin omises]`\
`5 DEBUG >> Dec 31 16:07:52 SENDED)`\
0-0-0-0-0-0-0-0_\
1-0-0-0-0-0-0-1_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-1-0-1-0-1_\
0-1-0-0-0-0-1-0 `[fin omises]`\
`[lignes omises]`\
`[Résultat final omis]`

Comme nous pouvons le voir, les matrices reçues en lignes 3 et 4 ont correctement été fusionnées avec la première matrice locale envoyée en ligne 2.

Matrice ligne 2 || Matrice ligne 3 || Matrice ligne 4 = Matrice ligne 5

__Résultat__\
<span style="color:green">Succès</span>

### Les plus courts chemins sont corrects
__Description__\
À la fin de l'exécution, les plus courts chemin sont des plus courts chemins

Résultat du noeud 6:\
`Shortest path to 0, length: 2, Path: 4 -> 3 -> 0`\
`Shortest path to 1, length: 3, Path: 4 -> 3 -> 0 -> 1`\
`Shortest path to 2, length: 3, Path: 4 -> 3 -> 0 -> 2`\
`Shortest path to 3, length: 1, Path: 4 -> 3`\
`Shortest path to 4, length: 0, Path: 4`\
`Shortest path to 5, length: 2, Path: 4 -> 3 -> 5`\
`Shortest path to 6, length: 2, Path: 4 -> 3 -> 6`\
`Shortest path to 7, length: 3, Path: 4 -> 3 -> 6 -> 7`

Conformément à la représentation du graphe, l'algorithme a trouvé les plus courts chemins du noeud 6 au reste

__Résultat__\
<span style="color:green">Succès</span>

### Réseau non connexe dysfonctionnel
__Description__\
Si le réseau n'est pas connexe, la recherche du plus court chemin ne fonctionne pas car l'algorithme ne s'arrête jamais. En effet, la condition d'arrêt des ondes n'est jamais remplie.

Pour ce test, nous avons retiré les liens du noeud 3 et de ses voisins.

__Résultat__\
<span style="color:green">Succès</span>

# Deuxième partie: algorithme sondes et échos

## Installation et utilisation
* Cloner le répertoire.
> `$ git clone https://github.com/MelvynHerzig/PRR_Labo4.git`

* Remplir le fichier de configuration _config.json_ à la racine du projet.
  * debug ( booléen, true/false ): Pour lancer les serveurs en mode debug (affiche les messages entrants et sortants)
  * versions ( string, indiquer "wave" ou "probe" ): Définit l'algorithme, utiliser "probe" pour l'algorithme sondes et échos
  * servers ( ip, port et numéros des voisins [0, Nb serveurs - 1]) Définition du réseau. Au minimum 1 serveur.
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
> Sous windows, il est également possible de démarrer les serveurs automatiquement avec le scripte `server/win_start_servers.bat`.\
> \
> L'ordre de démarrage n'est pas important. Durant cette étape, les serveurs s'inter-connectent. En conséquence, tant que tous ne sont pas allumés et connectés, ils n'acceptent que des connexions ayant une adresse IP source appartenant au fichier de configuration. De plus toute demande initiale sont refusées tant que le réseau n'est pas prêt.\
>\
> Lorsque un serveur a complétement démarré, il affiche les lignes suivantes:\
>![cmd](https://user-images.githubusercontent.com/34660483/147826831-11b7c10d-1aa8-401a-ac24-50f55c371e9e.png)\
> Si tous les serveurs affichent ces messages, le réseau est démarré et prêt à accepter les demandes clientes.

* Démarrer le client pour lancer une demande.  Un argument est nécessaire.
  * Numéro du serveur distant auquel se connecter.

> <u>Depuis le dossier _client_</u>\
>\
> Pour lancer un client qui démarre la recherche depuis le serveur 3\
> `$ go run . 3`\
>\
>Il peut être lancé avant que tous le serveur soit prêt. Sa demande est renouvellée automatiquement jusqu'à ce qu'il soit prêt.

* Lorsque la demande a été prise en compte chaque serveur affiche sa connaissance partielle des plus courts chemins jusqu'aux autres serveurs. Le noeud source affichera la liste réelle des plus courts chemins.

> Par exemple, le résultat du serveur 0 dans la config de ce _readme_:\
> ![probe](https://user-images.githubusercontent.com/34660483/147831762-c6c1d019-02ac-4940-93eb-1be21dd355ce.png)
> 

## Protocole UDP des sondes et échos
Les messages échangées antres les serveurs ont le format suivant:\
\<type> \<id> \<serveur source> \<matrice d'adjacence>

> Le type est soit "probe" pour les sondes ou "echo" pour les échos

> Id est l'identifiant unique de la série de sonde est écho en cours. En général, id à la valeur du serveur qui a déclanché la demande.

> Source est le numéro du serveur qui a émit le message.

> Chemin les plus courts est un tableau à deux dimensions des chemins les plus courts connus.

Par exemple, si le serveur 1 émet une sonde, initiée par une demande sur le serveur 5 et que sa topologie est:
```
[false, true, false]
[true, false, true]
[false, true, false]
```

Le message sérialisé sera "probe 5 1 0-1-0_1-0-1_0-1-0"

> Les lignes de la matrice sont séparées par des '_'

> Les colonnes de la matrice sont séparées par des '-'

> Si le plus court chemin n'est pas connu, il est traduit par un 'N'

> Les messages ne doivent pas excéder 1024 bytes. Il n'est pas conseillé de faire un réseau avec plus de 20 noeuds

## Tests
Les tests ont été réalisés avec la fichier de configuration de ce readme.

Pour ces tests nous n'allons pas analyser la sortie de toutes les consoles. Nous allons nous concentrer sur des situations spécifiques 

Les extraits suivant ont été tirés d'une demande effectuée sur le noeud 0.
### Un noeud se comporte correctement lorsqu'il reçoit sa première sonde
__Description__\
Lorsqu'un noeud reçoit une sonde, il envoie une sonde à tous ses voisins sauf à l'émetteur de la sonde. Dès qu'il a reçu toutes ses réponses, il envoie un écho
à l'émetteur de la première sonde.

Nous allons analyser la sortie de la console 3\
`1 DEBUG >> Dec 31 17:07:04 RECEIVED) probe 0 0 [matrice omise] from 0`\
`2 DEBUG >> Dec 31 17:07:04 SENDED) probe 0 3 [matrice omise] to 4`\
`3 DEBUG >> Dec 31 17:07:04 SENDED) probe 0 3 [matrice omise] to 5`\
`4 DEBUG >> Dec 31 17:07:04 SENDED) probe 0 3 [matrice omise] to 6`\
`5 DEBUG >> Dec 31 17:07:04 RECEIVED) echo 0 4 [matrice omise] from 4`\
`6 DEBUG >> Dec 31 17:07:04 RECEIVED) echo 0 5 [matrice omise] from 5`\
`7 DEBUG >> Dec 31 17:07:04 RECEIVED) echo 0 6 [matrice omise] from 6`\
`8 DEBUG >> Dec 31 17:07:04 SENDED) echo 0 3 [matrice omise] to 0`\
`[Résultat final omis]`

En ligne 1, le serveur 3 reçoit sa première sonde. Comme prévu en ligne 2-3-4, il envoie une sonde à tous ses voisins sauf à son "parent". Ligne 5-6-7, il récupère tous les échos, dès lors, il peut envoyer son écho à son parent en ligne 8.

__Résultat__\
<span style="color:green">Succès</span>

### Une sonde peut être considérée comme un écho
__Description__\
Dans un graphe avec un cycle, deux sondes peuvent se croiser sur le même lien, il devient alors nécessaire de pouvoir considérer une sonde comme un écho.

Ce phénomène arrive entre les noeuds 6 et 7.

Nous allons analyser la sortie de la console 6\
`1 DEBUG >> Dec 31 17:07:04 RECEIVED) probe 0 3 [matrice omise] from 3`\
`2 DEBUG >> Dec 31 17:07:04 SENDED) probe 0 6 [matrice omise] to 5`\
`3 DEBUG >> Dec 31 17:07:04 SENDED) probe 0 6 [matrice omise] to 7`\
`4 DEBUG >> Dec 31 17:07:04 RECEIVED) probe 0 5[matrice omise] from 5`\
`5 DEBUG >> Dec 31 17:07:04 RECEIVED) probe 0 7 [matrice omise] from 7`\
`6 DEBUG >> Dec 31 17:07:04 SENDED) echo 0 6 [matrice omise] to 3`\
`[Résultat final omis]`

En lignes 1-2-3, le serveur reçoit une sonde du serveur 3 ce qui déclanche l'envoi de ses sondes aux serveur 5 et 7. Ces deux derniers ont déjà reçu une sonde qu'ils ont eux même propagée. En conséquence, en lignes 4-5, le serveur 6 considère leurs sondes comme des échos. Finalement, en ligne 6, le serveur peut répondre à son parent.

__Résultat__\
<span style="color:green">Succès</span>

## Un noeud sans "enfant" répond directement à son parent
__Description__\
Il peut arriver qu'un noeud soit en extrémité d'une branche du graphe, il n'a donc aucun enfant à transmettre la sonde. En conséquence, il répond directement avec un écho.

Ce phénomène arrive sur le noeud 2.

`1 DEBUG >> Dec 31 17:07:04 RECEIVED) probe 0 0 [matrice omise] from 0`\
`2 DEBUG >> Dec 31 17:07:04 SENDED) echo 0 2 [matrice omise] to 0`\
`[Résultat final omis]`

En ligne 1, le noeud 1 reçoit la sonde de 0. Comme il n'a aucun autre voisin, il lui rend directement un écho avec la même liste des plus courts chemins.

__Résultat__\
<span style="color:green">Succès</span>

### Fusion correcte des topologies
__Description__\
Lorsque un serveur reçoit une sonde ou un écho, il met à jour sa matrice d'adjacence.

Nous allons analyser 4 lignes de la console 6\
`[ligne omise]`\
`3 DEBUG >> Dec 31 22:12:22 SENDED) probe 0 6`\
 0-1-1-1-0-0-0-0_\
 1-0-0-0-0-0-0-1_\
 0-0-0-0-0-0-0-0_\
 0-0-0-0-0-0-0-0_\
 0-0-0-0-0-0-0-0_\
 0-0-0-0-0-0-0-0_\
 0-0-0-1-0-1-0-1_\
 0-1-0-0-0-0-1-0 `to 5`\
`4 DEBUG >> Dec 31 22:12:22 RECEIVED) echo 0 5`\
0-1-1-1-0-0-0-0_\
1-0-0-0-0-0-0-1_\
0-0-0-0-0-0-0-0_\
1-0-0-0-1-1-1-0_\
0-0-0-0-0-0-0-0_\
0-0-0-1-0-0-1-0_\
0-0-0-1-0-1-0-1_\
0-1-0-0-0-0-1-0 `from 5`\
`5 DEBUG >> Dec 31 22:12:22 RECEIVED) probe 0 3`\
0-1-1-1-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
1-0-0-0-1-1-1-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0_\
0-0-0-0-0-0-0-0 `from 3`\
`6 DEBUG >> Dec 31 22:12:22 SENDED) echo 0 6`\
 0-1-1-1-0-0-0-0_\
 1-0-0-0-0-0-0-1_\
 0-0-0-0-0-0-0-0_\
 1-0-0-0-1-1-1-0_\
 0-0-0-0-0-0-0-0_\
 0-0-0-1-0-0-1-0_\
 0-0-0-1-0-1-0-1_\
 0-1-0-0-0-0-1-0 `to 7`\
`[Résultat final omis]`

Comme nous pouvons le voir, les matrices reçues en lignes 4 et 5 ont correctement été fusionnées avec la première matrice locale envoyée en ligne 3.

Matrice ligne 2 || Matrice ligne 3 || Matrice ligne 4 = Matrice ligne 5

__Résultat__\
<span style="color:green">Succès</span>

### Les plus courts chemins sont corrects
__Description__\
À la fin de l'exécution, les plus courts chemin sont des plus courts chemins

Résultat du noeud 0:\
`Shortest path to 0, length: 0, Path: 0`\
`Shortest path to 1, length: 1, Path: 0 -> 1`\
`Shortest path to 2, length: 1, Path: 0 -> 2`\
`Shortest path to 3, length: 1, Path: 0 -> 3`\
`Shortest path to 4, length: 2, Path: 0 -> 3 -> 4`\
`Shortest path to 5, length: 2, Path: 0 -> 3 -> 5`\
`Shortest path to 6, length: 2, Path: 0 -> 3 -> 6`\
`Shortest path to 7, length: 2, Path: 0 -> 1 -> 7`

Conformément à la représentation du graphe, l'algorithme a trouvé les plus courts chemins du noeud 0 au reste.

__Résultat__\
<span style="color:green">Succès</span>

### Réseau non connexe fonctionnel
__Description__\
Si le réseau n'est pas connexe, la recherche du plus court chemin fonctionne. Cependant,
certains chemin seront désignés comme unknown.

Pour ce test, nous avons retiré les liens du noeud 3 et de ses voisins.

`Shortest path to 0, length: 0, Path: 0`\
`Shortest path to 1, length: 1, Path: 0 -> 1`\
`Shortest path to 2, length: 1, Path: 0 -> 2`\
`unknown`\
`unknown`\
`Shortest path to 5, length: 4, Path: 0 -> 1 -> 7 -> 6 -> 5`\
`Shortest path to 6, length: 3, Path: 0 -> 1 -> 7 -> 6`\
`Shortest path to 7, length: 2, Path: 0 -> 1 -> 7`

En conséquence, les noeuds 3 et 4 ne peuvent plus être atteints. Ils sont inconnus sur le résultat final.

__Résultat__\
<span style="color:green">Succès</span>


## Note complémentaire

L'implémentation a été designée de sorte à pouvoir gérer plusieurs demandes simulatanément si elles sont émises par des serveurs différents (c'est à dire si elles ont des identifiants différents).

Comme les identifiants des sondes et échos correspondent au numéro du serveur source initial, chaque serveur intermédiaire stocke l'état spécifique à l'identifiant dans un tableau.

Par exemple, le serveur 0 peut gérer en même temps une demande effectuée par le serveur 1 et le serveur 5.
