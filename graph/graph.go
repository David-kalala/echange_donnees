package graph

import (
	"ed/model"
	"strconv"
	"strings"
)

// BuildNameIndex construit un map "nom en lowercase" -> ID
func BuildNameIndex(entities map[int]*model.Entity) map[string]int {
	index := make(map[string]int)
	for _, e := range entities {
		lower := strings.ToLower(e.Nom)
		index[lower] = e.ID
	}
	return index
}

// BuildReverseGraph construit le graphe inversé : pour chaque relation
//
//	(origine -> cible), on stocke dans reverseG[cible] = append(..., {TargetID=origine, Percent=...})
func BuildReverseGraph(entities map[int]*model.Entity, relations []*model.Relation) map[int][]model.OwnershipEdge {
	nameIndex := BuildNameIndex(entities)
	reverseG := make(map[int][]model.OwnershipEdge)

	for _, r := range relations {
		// trouver ID de r.Origine
		origineID := -1
		if id, ok := nameIndex[strings.ToLower(r.Origine)]; ok {
			origineID = id
		} else {
			continue
		}
		// trouver ID de r.Cible
		cibleID := -1
		if id, ok := nameIndex[strings.ToLower(r.Cible)]; ok {
			cibleID = id
		} else {
			continue
		}
		// calculer pourcentage
		p := parseValeurAsFloat(r.Valeur)

		// on ajoute dans reverseG[cibleID]
		reverseG[cibleID] = append(reverseG[cibleID], model.OwnershipEdge{
			TargetID: origineID,
			Percent:  p,
		})
	}

	return reverseG
}

// parseValeurAsFloat => Convertit la chaîne "100" => 100.0, ">50" => 50, "contrôle" => 100, etc.
func parseValeurAsFloat(val string) float64 {
	v := strings.ToLower(strings.TrimSpace(val))
	if v == "contrôle" || v == "participe" {
		return 100.0
	}
	if strings.HasPrefix(v, ">") {
		// par ex. ">50"
		f, err := strconv.ParseFloat(v[1:], 64)
		if err == nil {
			return f
		}
		return 100
	}
	// essayer de parse en float
	f, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return f
	}
	return 100
}

// GatherFinalOwners => fonction récursive pour trouver tous les propriétaires finaux
// entityID: entité qu’on cherche à remonter
// reverseG: graphe inverse
// currentPct: pourcentage (initialement 100 quand c’est un média qu’on examine)
// finalOwners : map[ownerID] => pourcentage cumulé
func GatherFinalOwners(entityID int, reverseG map[int][]model.OwnershipEdge, finalOwners map[int]float64, currentPct float64) {
	// Qui détient entityID ?
	owners := reverseG[entityID]
	if len(owners) == 0 {
		// Personne ne possède entityID => c’est un "final owner"
		finalOwners[entityID] += currentPct
		return
	}
	// Sinon, pour chacun des owners, on multiplie les pourcentages
	for _, edge := range owners {
		newPct := (edge.Percent / 100.0) * (currentPct / 100.0) * 100.0
		GatherFinalOwners(edge.TargetID, reverseG, finalOwners, newPct)
	}
}
