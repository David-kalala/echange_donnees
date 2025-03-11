package parser

import (
	"ed/model"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ReadEntities lit medias_francais.tsv et retourne un map d'Entity indexé par ID
func ReadEntities(filePath string) (map[int]*model.Entity, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1 // pour accepter des lignes de taille variable

	// Lire la première ligne (entête)
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de l'entête: %w", err)
	}
	fmt.Printf("Header Entities: %v\n", header)

	entities := make(map[int]*model.Entity)

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("erreur lecture ligne: %w", err)
		}
		// On attend au moins 4 colonnes [id, nom, typeLibelle, typeCode]
		if len(record) < 4 {
			continue // ignore la ligne
		}

		idStr := record[0]
		nomStr := record[1]
		typeLibelleStr := record[2]
		typeCodeStr := record[3]

		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		typeCodeInt, _ := strconv.Atoi(typeCodeStr)

		e := &model.Entity{
			ID:          idInt,
			Nom:         nomStr,
			TypeLibelle: typeLibelleStr,
			TypeCode:    typeCodeInt,
		}
		entities[idInt] = e
	}
	return entities, nil
}

// ReadRelations lit relations_medias_francais.tsv et renvoie un slice de Relation
func ReadRelations(filePath string) ([]*model.Relation, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	// Lire l'entête
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("erreur lecture entête relations: %w", err)
	}

	fmt.Printf("Header Relations: %v\n", header)

	var relations []*model.Relation
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if len(record) < 4 {
			continue
		}
		idStr := record[0]
		origineStr := record[1]
		valeurStr := record[2]
		cibleStr := record[3]

		var sourceStr, datePub, dateCons string
		if len(record) > 4 {
			sourceStr = record[4]
		}
		if len(record) > 5 {
			datePub = record[5]
		}
		if len(record) > 6 {
			dateCons = record[6]
		}

		idInt, _ := strconv.Atoi(idStr)

		r := &model.Relation{
			ID:               idInt,
			Origine:          origineStr,
			Valeur:           valeurStr,
			Cible:            cibleStr,
			Source:           sourceStr,
			DatePublication:  datePub,
			DateConsultation: dateCons,
		}
		relations = append(relations, r)
	}

	return relations, nil
}
