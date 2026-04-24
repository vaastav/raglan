package workflow

func IsEu(country string) bool {
	_, exists := EUCountries[country]
	return exists
}

func IsNa(country string) bool {
	_, exists := NACountries[country]
	return exists
}

func IsSa(country string) bool {
	_, exists := SACountries[country]
	return exists
}

func IsAs(country string) bool {
	_, exists := ASCountries[country]
	return exists
}

func IsAf(country string) bool {
	_, exists := AFCountries[country]
	return exists
}

func IsOc(country string) bool {
	_, exists := OCCountries[country]
	return exists
}
