package tess

import (
	"math"
	"testing"
)

// BenchmarkSimpleTriangle benchmarks a simple triangle tessellation
func BenchmarkSimpleTriangle(b *testing.B) {
	vertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 0.5, Y: 1, Z: 0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, vertices)
		tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
		
		// Access results to ensure they're computed
		_ = tessellator.GetVertexCount()
		_ = tessellator.GetElementCount()
		_ = tessellator.GetVertices()
		_ = tessellator.GetElements()
		
		tessellator.Delete()
	}
}

// BenchmarkComplexPolygon benchmarks a square with hole tessellation
func BenchmarkComplexPolygon(b *testing.B) {
	outerContour := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 4, Y: 0, Z: 0},
		{X: 4, Y: 4, Z: 0},
		{X: 0, Y: 4, Z: 0},
	}

	innerContour := []Vertex{
		{X: 1, Y: 1, Z: 0},
		{X: 3, Y: 1, Z: 0},
		{X: 2, Y: 3, Z: 0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, outerContour)
		tessellator.AddContour(2, innerContour)
		tessellator.SetOption(OptionConstrainedDelaunay, true)
		tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
		
		// Access results to ensure they're computed
		_ = tessellator.GetVertexCount()
		_ = tessellator.GetElementCount()
		_ = tessellator.GetVertices()
		_ = tessellator.GetElements()
		
		tessellator.Delete()
	}
}

// BenchmarkLargePolygon benchmarks a larger polygon with many vertices
func BenchmarkLargePolygon(b *testing.B) {
	// Create a circle approximation with many vertices
	vertices := make([]Vertex, 100)
	radius := 5.0
	centerX, centerY := 0.0, 0.0
	
	for i := 0; i < 100; i++ {
		angle := 2 * math.Pi * float64(i) / 100
		vertices[i] = Vertex{
			X: float32(centerX + radius*math.Cos(angle)),
			Y: float32(centerY + radius*math.Sin(angle)),
			Z: 0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, vertices)
		tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
		
		// Access results to ensure they're computed
		_ = tessellator.GetVertexCount()
		_ = tessellator.GetElementCount()
		_ = tessellator.GetVertices()
		_ = tessellator.GetElements()
		
		tessellator.Delete()
	}
}

// BenchmarkAddContour benchmarks just the AddContour operation
func BenchmarkAddContour(b *testing.B) {
	vertices := make([]Vertex, 1000)
	for i := 0; i < 1000; i++ {
		vertices[i] = Vertex{
			X: float32(i),
			Y: float32(i * 2),
			Z: 0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	vertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 0.5, Y: 1, Z: 0},
	}

	tessellator := NewTessellator()
	if tessellator == nil {
		b.Fatal("Failed to create tessellator")
	}
	defer tessellator.Delete()

	tessellator.AddContour(2, vertices)
	tessellator.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tessellator.GetVertexCount()
		_ = tessellator.GetElementCount()
		_ = tessellator.GetVertices()
		_ = tessellator.GetElements()
		_ = tessellator.GetVertexIndices()
	}
}

// BenchmarkWindingRules benchmarks different winding rules
func BenchmarkWindingRules(b *testing.B) {
	outerContour := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 4, Y: 0, Z: 0},
		{X: 4, Y: 4, Z: 0},
		{X: 0, Y: 4, Z: 0},
	}

	innerContour := []Vertex{
		{X: 1, Y: 1, Z: 0},
		{X: 3, Y: 1, Z: 0},
		{X: 2, Y: 3, Z: 0},
	}

	windingRules := []WindingRule{
		WindingOdd,
		WindingNonZero,
		WindingPositive,
		WindingNegative,
		WindingAbsGeqTwo,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule := windingRules[i%len(windingRules)]
		
		tessellator := NewTessellator()
		if tessellator == nil {
			b.Fatal("Failed to create tessellator")
		}

		tessellator.AddContour(2, outerContour)
		tessellator.AddContour(2, innerContour)
		tessellator.Tessellate(rule, ElementPolygons, 3, 2, nil)
		
		_ = tessellator.GetElementCount()
		tessellator.Delete()
	}
} 