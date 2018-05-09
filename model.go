package main

import "strconv"

type Point struct {
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`
}

type MonthPeriod struct {
	Year  int
	Month int
}

func (d *MonthPeriod) initFrom(s string) error {
	var err error
	d.Year, err = strconv.Atoi(s[0:4])
	if err != nil {
		return err
	}
	d.Month, err = strconv.Atoi(s[4:6])
	if err != nil {
		return err
	}
	return nil
}

type Range string //LC/MC/CC

type Zone string //M/O/I
