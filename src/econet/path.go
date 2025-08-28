package econet

import (
	"regexp"
)

// leadingNamedDiskRE matches a leading ":<name><digit>." at the start, e.g., ":DISK0." where <name>
// may include letters, digits, and any of the following characters: ! % & = - ~ ^ | \ @ { [ £ _ + ; } ] < > ? /
// It captures the trailing single digit which indicates the disk number and consumes the trailing dot.
var leadingNamedDiskRE = regexp.MustCompile(`^:[A-Za-z0-9!%&=\-~^|\\@{\[£_+;}\]<>?/\.]+([0-9])\.`)

// HasDiskPrefix reports whether s contains a colon-prefixed disk name (e.g., ":XYZ", ":DISK0").
func HasDiskPrefix(s string) bool {

	return leadingNamedDiskRE.MatchString(s)
}

// StripDiskPrefix returns the disk number and the remainder of the string after
// the leading disk prefix. It supports either formats:
//   - "<digit>:" (e.g., "0:$...")
//   - ":<letters><digit>." (e.g., ":DISK0.$...")
//
// If no such prefix exists, ok is false.
func StripDiskPrefix(s string) (disk int, rest string, ok bool) {
	if m := leadingNamedDiskRE.FindStringSubmatchIndex(s); m != nil {
		g1Start := m[2]
		disk = int(s[g1Start] - '0')
		rest = s[m[1]:]
		return disk, rest, true
	}
	if m := leadingNamedDiskRE.FindStringSubmatchIndex(s); m != nil {
		// For this pattern, group 1 is the digit after the letters; the full match ends after the dot.
		g1Start := m[2]
		disk = int(s[g1Start] - '0')
		rest = s[m[1]:]
		return disk, rest, true
	}
	return 0, s, false
}
