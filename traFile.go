package main

import (
	"errors"
	"log"
	"strconv"
)

//ANMOIS;DEP;ARR;NVOLS;PAX_FS;FSC;ZON

type TrafficInfo struct {
	FlightCount int
	PaxCount    int
}

type TrafficRecord struct {
	Date          MonthPeriod
	DepartureICAO string
	ArrivalICAO   string
	Range         Range
	Zone          Zone
	Value         TrafficInfo
}

func (t *TrafficRecord) initFrom(record []string) error {

	if len(record) != 7 {
		log.Println(record)
		return errors.New("Invalid record")
	}

	err := t.Date.initFrom(record[0])
	if err != nil {
		return err
	}
	t.DepartureICAO = record[1]
	t.ArrivalICAO = record[2]

	t.Value.FlightCount, err = strconv.Atoi(record[3])
	if err != nil {
		return err
	}
	t.Value.PaxCount, err = strconv.Atoi(record[4])
	if err != nil {
		return err
	}
	t.Range = Range(record[5])
	t.Zone = Zone(record[6])

	return nil
}

func ReadTraFile(traFile string) ([]TrafficRecord, error) {

	all, err := readCsv(traFile)
	if err != nil {
		return nil, err
	}

	var traffic []TrafficRecord

	for i, record := range all {
		if i == 0 {
			continue //header
		}
		var t TrafficRecord
		err := t.initFrom(record)
		if err != nil {
			return nil, err
		}
		traffic = append(traffic, t)
	}

	return traffic, nil
}
