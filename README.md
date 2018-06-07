# monapt-dataprep

English version below

Cette application a été réalisé par [Simon HEGE](https://twitter.com/simonhege) dans le cadre du concours 
Datavizz organisé lors des assises du transport aérien par la [Direction Générale de l'Aviation Civile 
(DGAC)](https://www.ecologique-solidaire.gouv.fr/direction-generale-laviation-civile-dgac). 


## License
Le code source de cette application est disponible sous licence MIT.

## Installation

### Prérequis
Les outils suivants doivent être installés préalablement:
 - git https://git-scm.com/
 - go https://golang.org

### Procédure d'installation
```bash
# Réccupérer les dépendances
go get github.com/xdbsoft/olap

# Cloner le dépot
git clone https://github.com/xdbsoft/monapt-dataprep.git 

# Compiler l'exéutable
go build .
```
L'exécutable est généré dans le dossier courant.

## Utilisation
Créer 2 sous-dossiers:
- `data`
- `json`

Dans le sous-dossier `data`, placer les fichiers CSV. 

Puis lancer la comande suivante.
```bash
./monapt-dataprep
```

Les fichiers JSON requis par https://github.com/xdbsoft/monaeroport sont générés dans le sous-dossier `json`.

# Engish version

## License
This software is distributed under MIT License

## Install and use

```bash
go get github.com/xdbsoft/olap

git clone https://github.com/xdbsoft/monapt-dataprep.git 

./monapt-dataprep
```
