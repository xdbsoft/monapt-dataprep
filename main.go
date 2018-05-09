package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/xdbsoft/olap"
)

func main() {

	maxYear := 2016

	//Read all files
	//APT
	airportRecords, err := ReadAptFile("./data/DataViz_APT.csv")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("APT file read. Records count: ", len(airportRecords))

	//TRA
	var traffic []TrafficRecord
	for year := 1990; year < maxYear+1; year++ {

		t, err := ReadTraFile(fmt.Sprintf("./data/DataViz_TRA_%d.csv", year))
		if err != nil {
			log.Fatal(err)
		}
		traffic = append(traffic, t...)
	}
	log.Println("TRA file read. Records count: ", len(traffic))

	//RET
	var delays []DelayRecord
	for year := 2012; year < maxYear+1; year++ {

		d, err := ReadRetFile(fmt.Sprintf("./data/DataViz_RET_%d.csv", year))
		if err != nil {
			log.Fatal(err)
		}
		delays = append(delays, d...)
	}
	log.Println("RET file read. Records count: ", len(delays))

	//EMI
	var emissions []EmissionRecord
	for year := 2000; year < maxYear+1; year++ {

		e, err := ReadEmiFile(fmt.Sprintf("./data/DataViz_EMI_%d.csv", year))
		if err != nil {
			log.Fatal(err)
		}
		emissions = append(emissions, e...)
	}
	log.Println("EMI file read. Records count: ", len(emissions))

	// Export airport info
	airports := airportInfo(airportRecords, "I")
	encode(airports, "json/airportInfo.json")
	encode(airportInfo(airportRecords, ""), "json/allAirportInfo.json")

	//Create the cube at airport/destination level for traffic and delays
	cubeLink := olap.Cube{
		Dimensions: []string{"year", "month", "icao", "country", "zone", "range", "direction", "dest_icao", "dest_country", "dest_zone", "dest_range"},
		Fields:     []string{"flights", "pax", "delay_15_pct_dep", "delay_15_pct_arr", "delay_mean_dep", "delay_mean_arr"},
	}

	aggregateCubeLink := func(aggregate, value []interface{}) []interface{} {

		flights, ok := aggregate[0].(int)
		if !ok {
			flights = 0
		}
		flights += value[0].(int)
		aggregate[0] = flights

		pax, ok := aggregate[1].(int)
		if !ok {
			pax = 0
		}
		pax += value[1].(int)
		aggregate[1] = pax

		//TODO: aggregate

		return aggregate
	}

	toPoint := func(date MonthPeriod, icaoDep string, icaoArr string) []interface{} {

		escaleDep := getEscale(icaoDep, airportRecords)
		escaleArr := getEscale(icaoArr, airportRecords)

		pt := make([]interface{}, len(cubeLink.Dimensions))
		pt[0] = date.Year
		pt[1] = date.Month
		pt[2] = escaleDep.ICAO
		pt[3] = escaleDep.CountryCode
		pt[4] = escaleDep.Zone
		pt[5] = escaleDep.Range
		pt[6] = "D"
		pt[7] = escaleArr.ICAO
		pt[8] = escaleArr.CountryCode
		pt[9] = escaleArr.Zone
		pt[10] = escaleArr.Range

		return pt
	}

	switchDepArr := func(pt []interface{}) []interface{} {

		ptArr := make([]interface{}, len(pt))
		copy(ptArr, pt)

		ptArr[6] = "A"
		ptArr[2], ptArr[7] = ptArr[7], ptArr[2]
		ptArr[3], ptArr[8] = ptArr[8], ptArr[3]
		ptArr[4], ptArr[9] = ptArr[9], ptArr[4]
		ptArr[5], ptArr[10] = ptArr[10], ptArr[5]

		return ptArr
	}

	for _, t := range traffic {

		pt := toPoint(t.Date, t.DepartureICAO, t.ArrivalICAO)
		ptReverse := switchDepArr(pt)

		data := make([]interface{}, len(cubeLink.Fields))
		data[0] = t.Value.FlightCount
		data[1] = t.Value.PaxCount

		dataReverse := make([]interface{}, len(cubeLink.Fields))
		dataReverse[0] = t.Value.FlightCount
		dataReverse[1] = t.Value.PaxCount

		cubeLink.Points = append(cubeLink.Points, pt, ptReverse)
		cubeLink.Data = append(cubeLink.Data, data, dataReverse)
	}

	mapKeyIndex := make(map[string]int)
	key := func(pt []interface{}) string {
		return fmt.Sprint(pt...)
	}
	for i, pt := range cubeLink.Points {

		k := key(pt)

		_, found := mapKeyIndex[k]
		if found {
			log.Println("Duplicated key", k)
		}

		mapKeyIndex[k] = i
	}

	for _, d := range delays {

		debug := false
		if d.DepartureICAO == "LFKJ" &&
			d.ArrivalICAO == "LFMN" &&
			d.Date.Year == 2016 &&
			d.Date.Month == 3 {
			debug = true
		}

		if debug {
			log.Println("Delay", d)
		}

		pt := toPoint(d.Date, d.DepartureICAO, d.ArrivalICAO)
		ptReverse := switchDepArr(pt)

		k := key(pt)
		kReverse := key(ptReverse)

		idx, found := mapKeyIndex[k]
		if !found {
			log.Println("Not found", d)
		}
		idxReverse, found := mapKeyIndex[kReverse]
		if !found {
			log.Println("Not found reverse", d)
		}

		if debug {
			log.Println("Key", k, kReverse, idx, idxReverse)
		}

		if d.Value.InfoAtDeparture {
			cubeLink.Data[idx][2] = d.Value.Percent15AtDeparture
			cubeLink.Data[idx][4] = d.Value.MeanDelayAtDeparture
			cubeLink.Data[idxReverse][2] = d.Value.Percent15AtDeparture
			cubeLink.Data[idxReverse][4] = d.Value.MeanDelayAtDeparture
		}
		if d.Value.InfoAtArrival {
			cubeLink.Data[idx][3] = d.Value.Percent15AtArrival
			cubeLink.Data[idx][5] = d.Value.MeanDelayAtArrival
			cubeLink.Data[idxReverse][3] = d.Value.Percent15AtArrival
			cubeLink.Data[idxReverse][5] = d.Value.MeanDelayAtArrival
		}

		if debug {
			log.Println("cubeLink.Data[idx]", cubeLink.Data[idx])
			log.Println("cubeLink.Data[idxReverse]", cubeLink.Data[idxReverse])
		}
	}

	//Slice per airport and export
	// for _, airport := range airportRecords {

	// 	if airport.CountryZone != "I" {

	cubeLinkAirport := cubeLink.Slice("icao", "LFKJ").RollUp(
		// 		cubeLinkAirport := cubeLink.Slice("icao", airport.ICAO).RollUp(
		[]string{"year", "month", "direction", "dest_icao", "dest_country", "dest_zone", "dest_range"},
		cubeLink.Fields,
		aggregateCubeLink,
		make([]interface{}, len(cubeLink.Fields)),
	)

	encode(cubeLinkAirport, fmt.Sprintf("json/%s_cube_links.json", "LFKJ"))
	encode(cubeLinkAirport.Rows(), fmt.Sprintf("json/%s_cube_rows.json", "LFKJ"))
	//		encode(cubeLinkAirport, fmt.Sprintf("json/%s_cube_links.json", airport.ICAO))
	// 	}
	// }

	//APT_OACI;APT_IATA;APT_NOM;APT_ISO2;APT_PAYS;PAYS_ZON;PAYS_FSC;APT_LAT;APT_LONG

	//ANMOIS;APT;FSC;ZON;CAT;MVT;PAX;PEQ;CO2;NOX;COVNM;TSP

	//ANMOIS;DEP;ARR;NVOLS;PAX_FS;FSC;ZON
	//ANMOIS;DEP;ARR;PC15_D;PC15_A;RETARD_D;RETARD_A

	// DEP/ARR -> APT/DIR/ESC
	// ESC -> ESC/FSC/ZON/

	//Create the cube at airport/destination level for traffic and delays
	//Dim: ANMOIS;APT;DIR;ESC;ESCFSC;ESCZON;APTFSC;APTZON
	//Fields: NVOLS;PAX_FS;PC15_D;PC15_A;RETARD_D;RETARD_A

	//Slice per airport and export

	//Create the cube at airport level only for traffic and emission
	//Dim: ANMOIS;APT;FSC;ZON;CAT
	//Fields: MVT;PAX;PEQ;CO2;NOX;COVNM;TSP

	//Slice per airport and export
}

func getEscale(icao string, airports []AirportRecord) Escale {

	var found *AirportRecord
	for i := range airports {
		if airports[i].ICAO == icao {
			found = &airports[i]
		}
	}

	if found == nil {
		panic(errors.New("Airport not found " + icao))
	}

	var e Escale
	e.Range = found.CountryRange
	e.Zone = found.CountryZone
	e.CountryCode = found.CountryCode
	e.ICAO = found.ICAO

	return e
}

func encode(d interface{}, path string) {

	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	e := json.NewEncoder(out)
	err = e.Encode(d)
	if err != nil {
		log.Fatal(err)
	}
}
