package utils

import "regexp"

func IsValidTaxCode(taxCode string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(taxCode)
}
func IsValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\d{1,10}$`)
	return re.MatchString(phone)
}
