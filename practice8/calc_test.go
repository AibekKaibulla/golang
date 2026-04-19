package main

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		res  int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := Add(tt.a, tt.b)
			if ans != tt.res {
				t.Errorf("Add(%d, %d) = %d; res %d", tt.a, tt.b, ans, tt.res)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		res      int
		hasError bool
	}{
		{"prime + notprime", 998244353, 5, 199648870, false},
		{"multiple + divisor", 2622729, 123, 21323, false},
		{"negative + positive", -100, 10, -10, false},
		{"both negative", -6, -3, 2, false},
		{"zero + any", 0, 10, 0, false},
		{"any + zero", 4, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := Divide(tt.a, tt.b)
			if tt.hasError {
				if err == nil {
					t.Errorf("Divide(%d, %d) expected error but got none", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				}
				if ans != tt.res {
					t.Errorf("Divide(%d, %d) = %d; res %d", tt.a, tt.b, ans, tt.res)
				}
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		res  int
	}{
		{"both positive", 1001, 1, 1000},
		{"positive + zero", 123321321312, 0, 123321321312},
		{"negative + positive", -328137120, 382019312, -710156432},
		{"both negative", -12312, -312393231, 312380919},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := Subtract(tt.a, tt.b)
			if ans != tt.res {
				t.Errorf("Subtract(%d, %d) = %d; res %d", tt.a, tt.b, ans, tt.res)
			}
		})
	}
}
