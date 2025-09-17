package econet

import (
	"time"
)

// CreateEconetDate The date is stored in two bytes. The first byte is the day of the month, the second.
func CreateEconetDate(date time.Time) []byte {

	var (
		d, m, y byte
	)
	/*
	   Byte   Bits        Meaning
	   ---------------------------
	   1       0-7         days
	   2       0 to 3      months
	           4 to 7      years
	                       undefined
	*/

	d = byte(date.Day())
	m = byte(date.Month())

	// 2021 is the minimum date allowed
	if date.Year() >= 2021 {
		y = byte(date.Year() - 2021)
	} else {
		y = 0 // i.e. 2021 which is the minimum date allowed
	}

	// shift the year into bits 4-7 and month into bits 0-3
	my := y<<4 | m
	return []byte{d, my}
}
