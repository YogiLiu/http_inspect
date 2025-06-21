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

type representedCountry struct {
	nameWithIso
	Type              string `json:"type"`
	IsInEuropeanUnion bool   `json:"IsInEuropeanUnion"`
}

type country struct {
	nameWithIso
	IsInEuropeanUnion bool `json:"IsInEuropeanUnion"`
}

type info struct {
	City               name               `json:"city"`
	PostalCode         string             `json:"postalCode"`
	Continent          continent          `json:"continent"`
	Subdivisions       []nameWithIso      `json:"subdivisions"`
	RepresentedCountry representedCountry `json:"representedCountry"`
	Country            country            `json:"country"`
	RegisteredCountry  country            `json:"registeredCountry"`
	Location           struct {
		TimeZone       string  `json:"timeZone"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		MetroCode      uint    `json:"metroCode"`
		AccuracyRadius uint16  `json:"accuracyRadius"`
	} `json:"location"`
	Traits struct {
		IsAnonymousProxy    bool `json:"isAnonymousProxy"`
		IsAnycast           bool `json:"isAnycast"`
		IsSatelliteProvider bool `json:"isSatelliteProvider"`
	} `json:"traits"`
}
