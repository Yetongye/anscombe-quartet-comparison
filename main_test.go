package main

import (
	"math"
	"testing"
)

// test if the results of linear regression are correct
func TestLinearRegression(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10} // y = 2x

	slope, intercept, err := linearRegression(x, y)
	if err != nil {
		t.Fatalf("Linear regression failed: %v", err)
	}

	if math.Abs(slope-2.0) > 0.0001 {
		t.Errorf("Expected slope 2.0, got %.5f", slope)
	}

	if math.Abs(intercept-0.0) > 0.0001 {
		t.Errorf("Expected intercept 0.0, got %.5f", intercept)
	}
}

// test if regression metrics calculation reasonable
func TestRegressionMetrics(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10} // y = 2x

	slope, intercept, _ := linearRegression(x, y)
	r2, rse, fstat := regressionMetrics(x, y, slope, intercept)

	if math.Abs(r2-1.0) > 0.0001 {
		t.Errorf("Expected R2 = 1.0, got %.5f", r2)
	}

	if rse > 0.0001 {
		t.Errorf("Expected RSE ~ 0, got %.5f", rse)
	}

	if fstat < 1000 {
		t.Errorf("Expected very large F-statistic, got %.5f", fstat)
	}
}

// benchmark regression performance
func BenchmarkLinearRegression(b *testing.B) {
	x := []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5}
	y := []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82, 5.68}

	for i := 0; i < b.N; i++ {
		_, _, _ = linearRegression(x, y)
	}
}
