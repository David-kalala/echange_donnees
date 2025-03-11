package main

import (
	"ed/graph"
	"ed/model"
	"ed/parser"
	"fmt"
	"log"
	"sort"
	"strings"
)

func main() {
	// 1) Lire les entités
	entities, err := parser.ReadEntities("data/medias_francais.tsv")
	if err != nil {
		log.Fatalf("Erreur lecture entities: %v", err)
	}

	// 2) Lire les relations
	relations, err := parser.ReadRelations("data/relations_medias_francais.tsv")
	if err != nil {
		log.Fatalf("Erreur lecture relations: %v", err)
	}

	// 3) Construire le graphe inversé
	reverseG := graph.BuildReverseGraph(entities, relations)

	// 4) Exécuter la requête sur un média
	mediaName := "Télérama"
	finalOwners := QueryMediaOwners(mediaName, entities, reverseG)

	// 5) Affichage formaté
	if len(finalOwners) == 0 {
		fmt.Printf("Pas de propriétaires finaux trouvés pour %s.\n", mediaName)
		return
	}

	fmt.Printf("%s est possédé par:\n", mediaName)
	// tri par pourcentage décroissant
	type pair struct {
		id  int
		pct float64
	}
	var pairs []pair
	for id, pct := range finalOwners {
		pairs = append(pairs, pair{id, pct})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].pct > pairs[j].pct
	})

	for _, p := range pairs {
		e := entities[p.id]
		if e != nil {
			fmt.Printf("  %.2f%% par %s (%s)\n", p.pct, e.Nom, e.TypeLibelle)
		}
	}
}

// QueryMediaOwners cherche l’ID du média par son nom, puis appelle GatherFinalOwners
func QueryMediaOwners(mediaName string, entities map[int]*model.Entity, reverseG map[int][]model.OwnershipEdge) map[int]float64 {
	// trouver l’ID du média
	mediaID := -1
	lowerMedia := strings.ToLower(mediaName)
	for _, e := range entities {
		if strings.ToLower(e.Nom) == lowerMedia && e.TypeCode == 3 {
			mediaID = e.ID
			break
		}
	}
	if mediaID < 0 {
		// Pas trouvé
		return nil
	}

	finalOwners := make(map[int]float64)
	// on part de currentPct=100 pour le média lui-même
	graph.GatherFinalOwners(mediaID, reverseG, finalOwners, 100.0)
	return finalOwners
}
