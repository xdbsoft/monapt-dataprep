package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

//ANMOIS;DEP;ARR;PC15_D;PC15_A;RETARD_D;RETARD_A

type DelayInfo struct {
	InfoAtDeparture      bool
	InfoAtArrival        bool
	Percent15AtDeparture float64
	Percent15AtArrival   float64
	MeanDelayAtDeparture float64
	MeanDelayAtArrival   float64
}

type DelayRecord struct {
	Date          MonthPeriod
	DepartureICAO string
	ArrivalICAO   string
	Value         DelayInfo
}

func (d *DelayRecord) initFrom(record []string) error {

	if len(record) != 7 {
		log.Println(record)
		return errors.New("Invalid record")
	}

	err := d.Date.initFrom(record[0])
	if err != nil {
		return err
	}
	d.DepartureICAO = record[1]
	d.ArrivalICAO = record[2]

	d.Value.InfoAtDeparture = (len(record[3]) > 0)
	d.Value.InfoAtArrival = (len(record[4]) > 0)

	if d.Value.InfoAtDeparture {
		d.Value.Percent15AtDeparture, err = strconv.ParseFloat(strings.Replace(record[3], ",", ".", -1), 64)
		if err != nil {
			return err
		}
		d.Value.MeanDelayAtDeparture, err = strconv.ParseFloat(strings.Replace(record[5], ",", ".", -1), 64)
		if err != nil {
			return err
		}
	}
	if d.Value.InfoAtArrival {
		d.Value.Percent15AtArrival, err = strconv.ParseFloat(strings.Replace(record[4], ",", ".", -1), 64)
		if err != nil {
			return err
		}
		d.Value.MeanDelayAtArrival, err = strconv.ParseFloat(strings.Replace(record[6], ",", ".", -1), 64)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadRetFile(retFile string) ([]DelayRecord, error) {

	all, err := readCsv(retFile)
	if err != nil {
		return nil, err
	}

	var delays []DelayRecord

	for i, record := range all {
		if i == 0 {
			continue //header
		}
		var d DelayRecord
		err := d.initFrom(record)
		if err != nil {
			return nil, err
		}
		delays = append(delays, d)
	}

	return delays, nil
}
