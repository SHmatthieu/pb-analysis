# app
The "process behavior analysis" project aims to study the usual operation of a Linux-type computer to detect unusual behavior that could correspond to a potential security problem.
## installation

-   go (le projet à été développé avec la version 1.17) https://go.dev/dl/ (réquis pour compiler le projet pas l'éxecution)
-   systemTap via un gestionnaire de package ou https://sourceware.org/systemtap/ftp/releases/ (nécessaire au fonctionnement)

## build

-   go build . 

## lancer

stap à besoin d'etre root sur la machine, si le terminal n'est pas sudo il faut rentrer le
mot de passe root dans la console pour que le script se lance.

(il est possible que des données s'affiche dans la console car stap et app sont lancer en
meme temps, il faut quand meme rentrer le mot de passe root)

Au lancement du mode surveillance beaucoup d'information s'affiche dans la console et on peut
parfois ne pas voir le message qui demande le mot de passe root

systemTap prend quelque seconde à démarrer.

### commande pour démarrer le programme (que le binaire app soit la)

-   sudo stap ./systapscript.stp | ./app 1 //mode apprentissage
-   sudo stap ./systapscript.stp | ./app 0 // mode monitor

## structure

Le projet est en 2 parties: l'application en go et le script systemeTap
les deux parties communiques grace à un pipe.

### systemTap

Le sript systemTap est sans systapscript.stp

### packages et code en go

Un programme en go est organisé en package, chaque package represente un
teme particulier (il ne doit pas y avoir d'interdépendance). Le point d'entré du code des dans le fichier main.go à la racine du projet.

-   data : contient tout se qui est rélatif à la persistance des données (interaction avec SQLite)
-   log : s'occupe de générer et de formater les logs.
-   proc : contient les structures Proc, Syscall et Prog et toute les méthodes liées
-   tap : s'occuper de parser les données de systemtap
-   tools : contient les deux methodes correspondantes au deux modes de fonctionnement.

### convention de code

La convention de code utilisé est celle imposé par go (car cela à un impacte par exemple sur la porté des fonctions et variables).

## fichier de configuration

dans le dossier config

-   PID.CONFIG permet de lister les pid à analyser (1 par ligne)
-   PROG.CONFIG permet de lister le chemin vers les applications à suivre
