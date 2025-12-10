package version

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Version
		wantErr bool
	}{
		{"0.8.20", Version{0, 8, 20}, false},
		{"0.8.0", Version{0, 8, 0}, false},
		{"0.8", Version{0, 8, 0}, false},
		{"1.0.0", Version{1, 0, 0}, false},
		{"invalid", Version{}, true},
		{"0.8.20.1", Version{}, true},
		{"a.b.c", Version{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		v1, v2 Version
		want   int
	}{
		{Version{0, 8, 0}, Version{0, 8, 0}, 0},
		{Version{0, 8, 0}, Version{0, 8, 1}, -1},
		{Version{0, 8, 1}, Version{0, 8, 0}, 1},
		{Version{0, 7, 0}, Version{0, 8, 0}, -1},
		{Version{0, 9, 0}, Version{0, 8, 0}, 1},
		{Version{0, 8, 0}, Version{1, 0, 0}, -1},
		{Version{1, 0, 0}, Version{0, 8, 0}, 1},
	}

	for _, tt := range tests {
		got := tt.v1.Compare(tt.v2)
		if got != tt.want {
			t.Errorf("(%v).Compare(%v) = %d, want %d", tt.v1, tt.v2, got, tt.want)
		}
	}
}

func TestVersionComparisons(t *testing.T) {
	v1 := Version{0, 8, 0}
	v2 := Version{0, 8, 20}
	v3 := Version{0, 8, 0}

	if !v1.LessThan(v2) {
		t.Errorf("%v should be less than %v", v1, v2)
	}
	if !v2.GreaterThan(v1) {
		t.Errorf("%v should be greater than %v", v2, v1)
	}
	if !v1.Equal(v3) {
		t.Errorf("%v should equal %v", v1, v3)
	}
	if !v1.LessThanOrEqual(v2) {
		t.Errorf("%v should be <= %v", v1, v2)
	}
	if !v1.LessThanOrEqual(v3) {
		t.Errorf("%v should be <= %v", v1, v3)
	}
	if !v2.GreaterThanOrEqual(v1) {
		t.Errorf("%v should be >= %v", v2, v1)
	}
	if !v1.GreaterThanOrEqual(v3) {
		t.Errorf("%v should be >= %v", v1, v3)
	}
}

func TestVersionString(t *testing.T) {
	v := Version{0, 8, 20}
	if v.String() != "0.8.20" {
		t.Errorf("String() = %q, want %q", v.String(), "0.8.20")
	}
}

func TestVersionIsZero(t *testing.T) {
	zero := Version{}
	nonZero := Version{0, 8, 0}

	if !zero.IsZero() {
		t.Error("zero version should be zero")
	}
	if nonZero.IsZero() {
		t.Error("non-zero version should not be zero")
	}
}

func TestNew(t *testing.T) {
	v := New(0, 8, 20)
	if v.Major != 0 || v.Minor != 8 || v.Patch != 20 {
		t.Errorf("New(0, 8, 20) = %v", v)
	}
}

func TestMustParse(t *testing.T) {
	v := MustParse("0.8.20")
	if v.String() != "0.8.20" {
		t.Errorf("MustParse failed")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParse should panic on invalid input")
		}
	}()
	MustParse("invalid")
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantRaw    string
		wantOp     string
		wantVer    Version
		wantErr    bool
	}{
		{
			name:    "caret constraint",
			source:  `pragma solidity ^0.8.0;`,
			wantRaw: "^0.8.0",
			wantOp:  "^",
			wantVer: Version{0, 8, 0},
		},
		{
			name:    "exact version",
			source:  `pragma solidity 0.8.20;`,
			wantRaw: "0.8.20",
			wantOp:  "",
			wantVer: Version{0, 8, 20},
		},
		{
			name:    "gte constraint",
			source:  `pragma solidity >=0.6.0;`,
			wantRaw: ">=0.6.0",
			wantOp:  ">=",
			wantVer: Version{0, 6, 0},
		},
		{
			name:    "with whitespace",
			source:  `  pragma   solidity   ^0.8.0 ;  `,
			wantRaw: "^0.8.0",
			wantOp:  "^",
			wantVer: Version{0, 8, 0},
		},
		{
			name:    "in contract",
			source:  `// SPDX\npragma solidity ^0.8.20;\ncontract Test {}`,
			wantRaw: "^0.8.20",
			wantOp:  "^",
			wantVer: Version{0, 8, 20},
		},
		{
			name:    "no pragma",
			source:  `contract Test {}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Detect(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Raw != tt.wantRaw {
				t.Errorf("Raw = %q, want %q", got.Raw, tt.wantRaw)
			}
			if got.Constraint != tt.wantOp {
				t.Errorf("Constraint = %q, want %q", got.Constraint, tt.wantOp)
			}
			if got.Version != tt.wantVer {
				t.Errorf("Version = %v, want %v", got.Version, tt.wantVer)
			}
		})
	}
}

func TestDetectAll(t *testing.T) {
	source := `
		pragma solidity ^0.8.0;
		pragma solidity >=0.6.0;
	`
	results, err := DetectAll(source)
	if err != nil {
		t.Fatalf("DetectAll() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("DetectAll() returned %d results, want 2", len(results))
	}
}
