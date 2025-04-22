// internal/metrics/util_test.go
package metrics

import (
	"testing"
)

func BenchmarkCalculateStats(b *testing.B) {
	// Create sample data
	data := make([]float64, 1000)
	for i := range data {
		data[i] = float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateStats(data)
	}
}

func BenchmarkFindClosestPointSize(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findClosestPointSize(4.2, standardPointSizes)
	}
}