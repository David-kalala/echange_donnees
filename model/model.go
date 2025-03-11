package model

type Entity struct {
	ID          int
	Nom         string
	TypeLibelle string
	TypeCode    int
}

type Relation struct {
	ID               int
	Origine          string
	Valeur           string
	Cible            string
	Source           string
	DatePublication  string
	DateConsultation string
}

// OwnershipEdge pour construire un graphe inversé (cible -> [owner...])
type OwnershipEdge struct {
	TargetID int     // entité propriétaire
	Percent  float64 // pourcentage de détention (0..100)
}
