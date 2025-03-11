# echange_donnees

Ce projet illustre un exemple d'exploitation de données tabulaires au format TSV (fichiers `medias_francais.tsv` et `relations_medias_francais.tsv`) pour remonter la chaîne de détention de médias. L'application est écrite en Go et procède comme suit :

1. **Parser les entités** (médias, personnes physiques, personnes morales, etc.) depuis le fichier `medias_francais.tsv`.
2. **Parser les relations** de détention (pourcentage, participation, contrôle, etc.) depuis le fichier `relations_medias_francais.tsv`.
3. **Construire un graphe inversé** (reverse graph) permettant de remonter des cibles (médias) vers leurs propriétaires.
4. **Calculer les propriétaires finaux** d'un média donné, en cumulant les pourcentages de détention.

Le code prend l’exemple du média "Télérama" et affiche la liste de ses propriétaires finaux avec leur pourcentage respectif de participation.

---

## Structure du projet

```text
.
├── data/
│   ├── medias_francais.tsv            # Données des entités
│   └── relations_medias_francais.tsv  # Données des relations capitalistiques
├── graph/
│   └── graph.go                       # Construction du graphe inversé et fonction récursive
├── model/
│   └── model.go                       # Définitions des structures de données (Entity, Relation, etc.)
├── parser/
│   └── parser.go                      # Fonctions de lecture/parsing des fichiers TSV
├── main.go                            # Point d'entrée principal : chargement, requête sur le média
├── go.mod                             # Fichier Go modules
├── app5_donnees_sequence2_TSV_...odt  # Exemple de document (non utilisé directement dans le code)
└── README.md                          # Présent fichier
```

---

## Prérequis

- **Go** version ≥ 1.23.3 (au minimum Go 1.18+ devrait fonctionner, mais ce `go.mod` indique 1.23.3).
- Un environnement pouvant exécuter les commandes `go build`, `go run`, etc.

---

## Installation & exécution

1. **Cloner** le dépôt ou copier l'ensemble des fichiers localement.
2. **Vérifier** que le module Go est initialisé (fichier `go.mod` présent).
3. **Compiler** le projet :
   ```bash
   go build -o echange_donnees main.go
   ```
4. **Exécuter** l'exécutable :
   ```bash
   ./echange_donnees
   ```
   Par défaut, l'application va chercher à lire les fichiers :
   - `data/medias_francais.tsv`
   - `data/relations_medias_francais.tsv`
   
   Puis, elle effectuera la requête de la chaîne de détention pour "Télérama" et affichera la liste de ses propriétaires finaux et leurs pourcentages.

---

## Explications principales

### 1. Parsing des entités et relations

- **`ReadEntities`** (dans `parser/parser.go`) : lit `medias_francais.tsv`, retourne un `map[int]*model.Entity`.
- **`ReadRelations`** (dans `parser/parser.go`) : lit `relations_medias_francais.tsv`, retourne un slice de `model.Relation`.

### 2. Construction du graphe inversé

- **`BuildReverseGraph`** (dans `graph/graph.go`) : à partir de la liste de relations, crée une structure `reverseG` dans laquelle chaque clé (ID de cible) pointe vers un tableau d'edges `{ TargetID, Percent }` représentant la détention de la cible par `TargetID` avec un pourcentage donné.

### 3. Calcul des propriétaires finaux

- **`GatherFinalOwners`** (fonction récursive dans `graph/graph.go`) : pour un média de base, on remonte récursivement dans le graphe inversé.  
  - S’il n’y a pas de propriétaire (pas de clé), c’est un propriétaire final.  
  - Sinon, on multiplie les pourcentages à chaque niveau de remontée.

### 4. Affichage

- **Exemple** dans `main.go` : on interroge les propriétaires de "Télérama". Les résultats sont ensuite triés et affichés.

---

## Personnalisation

- Pour changer le média dont on veut afficher la chaîne de détention, il suffit de modifier la variable `mediaName` dans `main.go`.
- Pour intégrer ce code dans un autre projet, on peut directement réutiliser les fonctions de `parser` et `graph` en adaptant la logique de requête.

---

## Limitations et évolutions possibles

- Le calcul suppose que tous les pourcentages s’additionnent simplement et que les lignes de détention ne s’excluent pas mutuellement.
- Les champs `contrôle`, `participe`, `>50`, etc. sont simplifiés à 100 ou une valeur numérique ; il est possible de raffiner cette logique pour d’autres cas.
- Pour analyser d'autres médias ou d'autres entrées, modifier le code ou ajouter une interface CLI serait envisageable.

---

## Licence

Ce projet est fourni à titre d’exemple pédagogique. Vous pouvez l’adapter et l’étendre suivant vos besoins. Aucune licence explicite n’a été fournie dans les fichiers du dépôt.

