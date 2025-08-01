package tess

import (
	"math"
	"testing"
)

// BenchmarkSimpleTriangle benchmarks a simple triangle tessellation
func BenchmarkSimpleTriangle(b *testing.B) {
	vertices := []float32{
		0, 0, 0,
		1, 0, 0,
		0.5, 1, 0,
	}

	b.ResetTimer()
	for range b.N {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, vertices)
		_, _, _ = tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)

		tessellator.Delete()
	}
}

// BenchmarkComplexPolygon benchmarks a square with hole tessellation
func BenchmarkComplexPolygon(b *testing.B) {
	outerContour := []float32{
		0, 0, 0,
		4, 0, 0,
		4, 4, 0,
		0, 4, 0,
	}

	innerContour := []float32{
		1, 1, 0,
		3, 1, 0,
		2, 3, 0,
	}

	b.ResetTimer()
	for range b.N {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, outerContour)
		tessellator.AddContour(2, innerContour)
		tessellator.SetOption(OptionConstrainedDelaunay, true)
		_, _, _ = tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)

		tessellator.Delete()
	}
}

// BenchmarkLargePolygon benchmarks a larger polygon with many vertices
func BenchmarkLargePolygon(b *testing.B) {
	// Create a circle approximation with many vertices
	vertices := make([]float32, 100*2)
	radius := 5.0
	centerX, centerY := 0.0, 0.0

	for i := range 100 {
		angle := 2 * math.Pi * float64(i) / 100
		vertices[i*2] = float32(centerX + radius*math.Cos(angle))
		vertices[i*2+1] = float32(centerY + radius*math.Sin(angle))
	}

	b.ResetTimer()
	for range b.N {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, vertices)
		_, _, _ = tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)

		tessellator.Delete()
	}
}

// BenchmarkAddContour benchmarks just the AddContour operation
func BenchmarkAddContour(b *testing.B) {
	vertices := make([]float32, 1000)
	for i := 0; i < 1000; i++ {
		vertices[i*2] = float32(i)
		vertices[i*2+1] = float32(i * 2)
	}

	b.ResetTimer()
	for range b.N {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, vertices)
		tessellator.Delete()
	}
}

// BenchmarkDataRetrieval benchmarks the data retrieval operations
func BenchmarkDataRetrieval(b *testing.B) {
	vertices := []float32{
		0, 0, 0,
		1, 0, 0,
		0.5, 1, 0,
	}

	tessellator := NewTessellator()
	if tessellator == nil {
		b.Fatal("Failed to create tessellator")
	}
	defer tessellator.Delete()

	tessellator.AddContour(2, vertices)
	resultVertices, resultIndices, _ := tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)

	b.ResetTimer()
	for range b.N {
		// Access the results to ensure they're computed
		_ = len(resultVertices)
		_ = len(resultIndices)
	}
}

// BenchmarkWindingRules benchmarks different winding rules
func BenchmarkWindingRules(b *testing.B) {
	outerContour := []float32{
		0, 0, 0,
		4, 0, 0,
		4, 4, 0,
		0, 4, 0,
	}

	innerContour := []float32{
		1, 1, 0,
		3, 1, 0,
		2, 3, 0,
	}

	windingRules := []WindingRule{
		WindingOdd,
		WindingNonZero,
		WindingPositive,
		WindingNegative,
		WindingAbsGeqTwo,
	}

	b.ResetTimer()
	for i := range b.N {
		rule := windingRules[i%len(windingRules)]

		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, outerContour)
		tessellator.AddContour(2, innerContour)
		_, _, _ = tessellator.Tessellate(rule, ElementPolygons, 3, 2, nil)

		tessellator.Delete()
	}
}
