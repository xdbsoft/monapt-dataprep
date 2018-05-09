package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

//ANMOIS;APT;FSC;ZON;CAT;MVT;PAX;PEQ;CO2;NOX;COVNM;TSP

type EmissionInfo struct {
	Movements                      int
	PaxCount                       int
	PaxEquivalentCount             int
	Co2KiloTons                    float64
	NoxTons                        float64
	NonMethaneVOCTons              float64
	TotalParticlesInSuspensionTons float64
}

type EmissionRecord struct {
	Date     MonthPeriod
	ICAO     string
	Range    Range
	Zone     Zone
	Category string
	Value    EmissionInfo
}

func (e *EmissionRecord) initFrom(record []string) error {

	if len(record) != 12 {
		log.Println(record)
		return errors.New("Invalid record")
	}

	err := e.Date.initFrom(record[0])
	if err != nil {
		return err
	}
	e.ICAO = record[1]
	e.Range = Range(record[2])
	e.Zone = Zone(record[3])
	e.Category = record[4]

	e.Value.Movements, err = strconv.Atoi(record[5])
	if err != nil {
		return err
	}
	e.Value.PaxCount, err = strconv.Atoi(record[6])
	if err != nil {
		return err
	}
	e.Value.PaxEquivalentCount, err = strconv.Atoi(record[7])
	if err != nil {
		return err
	}
	e.Value.Co2KiloTons, err = strconv.ParseFloat(strings.Replace(record[8], ",", ".", -1), 64)
	if err != nil {
		return err
	}
	e.Value.NoxTons, err = strconv.ParseFloat(strings.Replace(record[9], ",", ".", -1), 64)
	if err != nil {
		return err
	}
	e.Value.NonMethaneVOCTons, err = strconv.ParseFloat(strings.Replace(record[10], ",", ".", -1), 64)
	if err != nil {
		return err
	}
	e.Value.TotalParticlesInSuspensionTons, err = strconv.ParseFloat(strings.Replace(record[11], ",", ".", -1), 64)
	if err != nil {
		return err
	}

	return nil
}

func ReadEmiFile(emiFile string) ([]EmissionRecord, error) {

	all, err := readCsv(emiFile)
	if err != nil {
		return nil, err
	}

	var emissions []EmissionRecord

	for i, record := range all {
		if i == 0 {
			continue //header
		}
		var e EmissionRecord
		err := e.initFrom(record)
		if err != nil {
			return nil, err
		}
		emissions = append(emissions, e)
	}

	return emissions, nil
}
