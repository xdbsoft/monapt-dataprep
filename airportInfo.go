package main

type AirportInfo struct {
	ICAO        string `json:"icao,omitempty"`
	IATA        string `json:"iata,omitempty"`
	Name        string `json:"name,omitempty"`
	CountryZone Zone   `json:"zone,omitempty"`
	Position    Point  `json:"pos,omitempty"`
}

func airportInfo(airports []AirportRecord, excludedZone Zone) []AirportInfo {

	var infos []AirportInfo
	for _, a := range airports {

		if a.CountryZone != excludedZone && a.ICAO[0] != 'Z' {
			var i AirportInfo
			i.ICAO = a.ICAO
			i.IATA = a.IATA
			i.Name = a.Name
			i.CountryZone = a.CountryZone
			i.Position = a.Position

			infos = append(infos, i)
		}
	}

	return infos
}
