package econet

import (
	"reflect"
	"testing"
	"time"
)

// TestCreateEconetDate_Table verifies the bit packing performed by CreateEconetDate.
// According to the current implementation in time.go:
//   - Byte0: day (1..31)
//   - Byte1: lower nibble = month (time.Month: 1..12)
//     upper nibble = low 4 bits of the year (byte(date.Year()) >> 4 when reversed),
//     i.e., Byte1 = (byte(Year) << 4) | byte(Month)
//
// This test locks in that behavior using a table-driven approach.
func TestCreateEconetDate_Table(t *testing.T) {
	mk := func(y int, m time.Month, d int, loc *time.Location) time.Time {
		if loc == nil {
			loc = time.UTC
		}
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	}

	tests := []struct {
		name string
		t    time.Time
		want []byte
	}{
		{
			name: "TypicalDate_2025_09_14",
			t:    mk(2025, time.September, 14, time.UTC),
			want: []byte{byte(14), byte(4)<<4 | byte(9)},
		},
		{
			name: "MonthBoundary_January",
			t:    mk(2020, time.January, 1, time.UTC),
			want: []byte{byte(1), byte(0)<<4 | byte(1)},
		},
		{
			name: "MonthBoundary_December",
			t:    mk(1999, time.December, 31, time.UTC),
			want: []byte{byte(31), byte(0)<<4 | byte(12)},
		},
		{
			name: "YearLowByteWrap_256",
			// 256 decimal -> 0x0100; byte(256) == 0x00, so upper nibble becomes 0
			t:    mk(256, time.July, 4, time.UTC),
			want: []byte{byte(4), byte(0)<<4 | byte(7)},
		},
		{
			name: "YearLowByte_255",
			// 255 decimal -> 0xFF; upper nibble 0xF
			t:    mk(255, time.March, 15, time.UTC),
			want: []byte{byte(15), byte(0)<<4 | byte(3)},
		},
		{
			name: "NonUTC_LocationIgnoredForDate",
			// Ensure location/timezone doesn\'t influence the day/month/year extraction
			t:    mk(2023, time.November, 5, time.FixedZone("X", 5*3600)),
			want: []byte{byte(5), byte(2023-2021)<<4 | byte(11)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CreateEconetDate(tc.t)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("CreateEconetDate(%v) mismatch. want=%#v got=%#v", tc.t, tc.want, got)
			}
		})
	}
}
