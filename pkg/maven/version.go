package maven

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a Maven version with semantic versioning
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Qualifier  string // e.g., "alpha", "beta", "RC1"
	Original   string // Original version string
}

// ParseVersion parses a Maven version string into a Version struct
func ParseVersion(v string) (*Version, error) {
	if v == "" {
		return nil, fmt.Errorf("empty version string")
	}

	original := v
	version := &Version{Original: original}

	// Split on dash to separate qualifier
	parts := strings.SplitN(v, "-", 2)
	numericPart := parts[0]
	if len(parts) > 1 {
		version.Qualifier = parts[1]
	}

	// Split numeric part by dots
	numbers := strings.Split(numericPart, ".")
	if len(numbers) == 0 || len(numbers) > 3 {
		return nil, fmt.Errorf("invalid version format: %s", v)
	}

	// Parse major version
	major, err := strconv.Atoi(numbers[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", numbers[0])
	}
	version.Major = major

	// Parse minor version if present
	if len(numbers) > 1 {
		minor, err := strconv.Atoi(numbers[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", numbers[1])
		}
		version.Minor = minor
	}

	// Parse patch version if present
	if len(numbers) > 2 {
		patch, err := strconv.Atoi(numbers[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", numbers[2])
		}
		version.Patch = patch
	}

	return version, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.Original
}

// Compare compares two versions
// Returns: -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	// Compare major
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	// Compare minor
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	// Compare patch
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Compare qualifier
	return compareQualifier(v.Qualifier, other.Qualifier)
}

// compareQualifier compares version qualifiers
func compareQualifier(q1, q2 string) int {
	// No qualifier is considered higher than any qualifier (release > RC > beta > alpha)
	if q1 == "" && q2 == "" {
		return 0
	}
	if q1 == "" {
		return 1 // No qualifier wins
	}
	if q2 == "" {
		return -1
	}

	// Lexicographic comparison for qualifiers
	if q1 < q2 {
		return -1
	}
	if q1 > q2 {
		return 1
	}
	return 0
}

// IsValid checks if version string is valid
func IsValid(v string) bool {
	_, err := ParseVersion(v)
	return err == nil
}

// MatchesPrefix checks if version matches a prefix (e.g., "3.8" matches "3.8.6")
func (v *Version) MatchesPrefix(prefix string) bool {
	prefixVer, err := ParseVersion(prefix)
	if err != nil {
		return false
	}

	// Check major
	if v.Major != prefixVer.Major {
		return false
	}

	// If prefix has minor, check it
	if prefix != fmt.Sprintf("%d", prefixVer.Major) {
		if v.Minor != prefixVer.Minor {
			return false
		}
	}

	return true
}

// SortVersions sorts versions in descending order (newest first)
func SortVersions(versions []string) ([]string, error) {
	parsed := make([]*Version, 0, len(versions))
	for _, v := range versions {
		pv, err := ParseVersion(v)
		if err != nil {
			return nil, fmt.Errorf("invalid version %s: %w", v, err)
		}
		parsed = append(parsed, pv)
	}

	// Bubble sort (simple, sufficient for typical use case)
	for i := 0; i < len(parsed); i++ {
		for j := i + 1; j < len(parsed); j++ {
			if parsed[i].Compare(parsed[j]) < 0 {
				parsed[i], parsed[j] = parsed[j], parsed[i]
			}
		}
	}

	// Convert back to strings
	result := make([]string, len(parsed))
	for i, v := range parsed {
		result[i] = v.String()
	}

	return result, nil
}
