package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Entity is a structure representing an observation,
type Entity struct {
	Id           int         `json:"id" form:"id"`
	TaxonId      int         `json:"taxon_id" form:"taxon_id"`
	Uuid         uuid.UUID   `json:"uuid" form:"uuid"`
	PlaceGuess   string      `json:"place_guess" form:"place_guess"`
	SpeciesGuess string      `json:"species_guess" form:"species_guess"`
	Latitude     string      `json:"latitude" form:"latitude"`
	Longitude    string      `json:"longitude" form:"longitude"`
	ObservedOn   pgtype.Date `json:"observed_on" form:"observed_on"`
	TimeZone     string      `json:"time_zone" form:"time_zone"`
}

// EntityQuery is designed for gathering data about entities for search queries with either taxon_id, or in between to dates
type EntityQuery struct {
	TaxonId int         `json:"taxon_id" form:"taxon_id"`
	Date1   pgtype.Date `json:"date1" form:"date1"`
	Date2   pgtype.Date `json:"date2" form:"date2"`
}

// DateQuery type is meant to search in between 2 specific dates
type DateQuery struct {
	Date1 pgtype.Date `json:"date1" form:"date1"`
	Date2 pgtype.Date `json:"date2" form:"date2"`
}
