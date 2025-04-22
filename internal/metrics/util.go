package metrics

import (
	"math"
	"sort"
)

// Standard story point sizes for grouping
var standardPointSizes = []float64{1, 2, 3, 5, 8, 13, 21}

// calculateStats calculates statistical values from a set of data points
func calculateStats(values []float64) (min, max, avg, median float64) {
	if len(values) == 0 {
		return 0, 0, 0, 0
	}
	
	// Sort for min, max, median
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	min = sorted[0]
	max = sorted[len(sorted)-1]
	
	// Calculate average
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	avg = sum / float64(len(values))
	
	// Calculate median
	if len(sorted)%2 == 0 {
		// Even number of values
		middle1 := sorted[len(sorted)/2-1]
		middle2 := sorted[len(sorted)/2]
		median = (middle1 + middle2) / 2
	} else {
		// Odd number of values
		median = sorted[len(sorted)/2]
	}
	
	return min, max, avg, median
}

// calculateCorrelation calculates the Pearson correlation coefficient between two sets of values
func calculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}
	
	n := float64(len(x))
	
	// Calculate sums
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	sumY2 := 0.0
	
	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}
	
	// Calculate correlation coefficient
	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))
	
	if denominator == 0 {
		return 0
	}
	
	return numerator / denominator
}

// findClosestPointSize finds the closest story point size from standard sizes
func findClosestPointSize(estimate float64, standardSizes []float64) float64 {
	if estimate == 0 {
		return 0
	}
	
	closest := standardSizes[0]
	minDiff := math.Abs(estimate - closest)
	
	for _, size := range standardSizes {
		diff := math.Abs(estimate - size)
		if diff < minDiff {
			minDiff = diff
			closest = size
		}
	}
	
	return closest
}