package services

import "fmt"

func GetYearLevelString(year int) string {
	switch year {
	case 1:
		return "First"
	case 2:
		return "Second"
	case 3:
		return "Third"
	case 4:
		return "Fourth"
	default:
		return fmt.Sprintf("%dth", year)
	}
}
