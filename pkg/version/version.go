// Package version provides Solidity version detection and comparison
package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a Solidity version
type Version struct {
	Major int
	Minor int
	Patch int
}

// New creates a new Version
func New(major, minor, patch int) Version {
	return Version{Major: major, Minor: minor, Patch: patch}
}

// String returns the version as a string
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare compares two versions
// Returns -1 if v < other, 0 if equal, 1 if v > other
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// LessThan returns true if v < other
func (v Version) LessThan(other Version) bool {
	return v.Compare(other) < 0
}

// LessThanOrEqual returns true if v <= other
func (v Version) LessThanOrEqual(other Version) bool {
	return v.Compare(other) <= 0
}

// GreaterThan returns true if v > other
func (v Version) GreaterThan(other Version) bool {
	return v.Compare(other) > 0
}

// GreaterThanOrEqual returns true if v >= other
func (v Version) GreaterThanOrEqual(other Version) bool {
	return v.Compare(other) >= 0
}

// Equal returns true if v == other
func (v Version) Equal(other Version) bool {
	return v.Compare(other) == 0
}

// IsZero returns true if version is unset (0.0.0)
func (v Version) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

// Parse parses a version string like "0.8.20" or "0.8"
func Parse(s string) (Version, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return Version{}, fmt.Errorf("invalid version format: %s", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch := 0
	if len(parts) == 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return Version{}, fmt.Errorf("invalid patch version: %s", parts[2])
		}
	}

	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// MustParse parses a version string and panics on error
func MustParse(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// DetectedVersion represents the version info extracted from source code
type DetectedVersion struct {
	Raw        string   // Raw pragma string, e.g., "^0.8.0"
	Constraint string   // Constraint operator, e.g., "^", ">=", etc.
	Version    Version  // Parsed version
}

// Detect extracts Solidity version information from source code
// Returns the detected version info, or error if no pragma found
func Detect(source string) (*DetectedVersion, error) {
	// Find pragma solidity statement
	pragmaRe := regexp.MustCompile(`pragma\s+solidity\s+([^;]+);`)
	matches := pragmaRe.FindStringSubmatch(source)

	if matches == nil {
		return nil, fmt.Errorf("no pragma solidity found")
	}

	raw := strings.TrimSpace(matches[1])
	
	// Parse the constraint
	constraintRe := regexp.MustCompile(`^(\^|~|>=|<=|>|<|=)?(\d+\.\d+(\.\d+)?)`)
	constraintMatches := constraintRe.FindStringSubmatch(raw)
	
	if constraintMatches == nil {
		return nil, fmt.Errorf("invalid pragma version: %s", raw)
	}

	constraint := constraintMatches[1]
	versionStr := constraintMatches[2]

	version, err := Parse(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in pragma: %w", err)
	}

	return &DetectedVersion{
		Raw:        raw,
		Constraint: constraint,
		Version:    version,
	}, nil
}

// DetectAll extracts all Solidity version pragmas from source code
// Useful for files with multiple pragma statements
func DetectAll(source string) ([]*DetectedVersion, error) {
	pragmaRe := regexp.MustCompile(`pragma\s+solidity\s+([^;]+);`)
	allMatches := pragmaRe.FindAllStringSubmatch(source, -1)

	if len(allMatches) == 0 {
		return nil, fmt.Errorf("no pragma solidity found")
	}

	var results []*DetectedVersion
	constraintRe := regexp.MustCompile(`^(\^|~|>=|<=|>|<|=)?(\d+\.\d+(\.\d+)?)`)

	for _, matches := range allMatches {
		raw := strings.TrimSpace(matches[1])
		constraintMatches := constraintRe.FindStringSubmatch(raw)

		if constraintMatches == nil {
			continue
		}

		version, err := Parse(constraintMatches[2])
		if err != nil {
			continue
		}

		results = append(results, &DetectedVersion{
			Raw:        raw,
			Constraint: constraintMatches[1],
			Version:    version,
		})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no valid pragma solidity found")
	}

	return results, nil
}
