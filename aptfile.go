package main

import (
	"errors"
	"strconv"
	"strings"
)

//APT_OACI;APT_IATA;APT_NOM;APT_ISO2;APT_PAYS;PAYS_ZON;PAYS_FSC;APT_LAT;APT_LONG

type Escale struct {
	Range       Range
	Zone        Zone
	CountryCode string
	ICAO        string
}

type AirportRecord struct {
	ICAO         string
	IATA         string
	Name         string
	CountryCode  string
	CountryName  string
	CountryRange Range
	CountryZone  Zone
	Position     Point
}

func (a *AirportRecord) initFrom(record []string) error {

	if len(record) != 9 {
		return errors.New("Invalid record")
	}

	a.ICAO = record[0]
	a.IATA = record[1]
	a.Name = record[2]
	a.CountryCode = record[3]
	a.CountryName = record[4]
	a.CountryZone = Zone(record[5])
	a.CountryRange = Range(record[6])
	if len(record[7]) > 0 || len(record[8]) > 0 {
		var err error
		a.Position.Latitude, err = strconv.ParseFloat(strings.Replace(record[7], ",", ".", -1), 64)
		if err != nil {
			return err
		}
		a.Position.Longitude, err = strconv.ParseFloat(strings.Replace(record[8], ",", ".", -1), 64)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadAptFile(aptFile string) ([]AirportRecord, error) {

	all, err := readCsv(aptFile)
	if err != nil {
		return nil, err
	}

	var airports []AirportRecord

	for i, record := range all {
		if i == 0 {
			continue //header
		}
		var a AirportRecord
		err := a.initFrom(record)
		if err != nil {
			return nil, err
		}
		airports = append(airports, a)
	}

	return airports, nil
}
