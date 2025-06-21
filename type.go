package main

type name map[string]string

func (n name) Filter(codes ...string) map[string]string {
	m := make(map[string]string)
	for _, code := range codes {
		if v, ok := n[code]; ok {
			m[code] = v
		}
	}
	return m
}

type nameWithIso struct {
	Names   name   `json:"names"`
	IsoCode string `json:"isoCode"`
}

type continent struct {
	Names name   `json:"names"`
	Code  string `json:"code"`
}

type country struct {
	nameWithIso
	IsInEuropeanUnion bool `json:"isInEuropeanUnion"`
}

type representedCountry struct {
	country
	Type string `json:"type"`
}

type info struct {
	Continent    continent     `json:"continent"`
	Country      country       `json:"country"`
	Subdivisions []nameWithIso `json:"subdivisions"`
	City         name          `json:"city"`
	Location     struct {
		TimeZone       string  `json:"timeZone"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		MetroCode      uint    `json:"metroCode"`
		AccuracyRadius uint16  `json:"accuracyRadius"`
	} `json:"location"`
	PostalCode         string             `json:"postalCode"`
	RepresentedCountry representedCountry `json:"representedCountry"`
	RegisteredCountry  country            `json:"registeredCountry"`
	Traits             struct {
		IsAnonymousProxy    bool `json:"isAnonymousProxy"`
		IsAnycast           bool `json:"isAnycast"`
		IsSatelliteProvider bool `json:"isSatelliteProvider"`
	} `json:"traits"`
}
