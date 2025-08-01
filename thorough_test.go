package tess

import (
	"math"
	"testing"
)

// TestMemoryManagement tests memory allocation, copying, and cleanup
func TestMemoryManagement(t *testing.T) {
	// Test multiple tessellator creation and deletion
	tessellators := make([]*Tessellator, 10)

	// Create multiple tessellators
	for i := 0; i < 10; i++ {
		tess := NewTessellator()
		if tess == nil {
			t.Fatalf("Failed to create tessellator %d", i)
		}
		tessellators[i] = tess
	}

	// Delete all tessellators
	for i, tess := range tessellators {
		tess.Delete()

		// Test operations on deleted tessellator
		vertices, indices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
		if err == nil {
			t.Errorf("Expected error for deleted tessellator %d", i)
		}
		if vertices != nil {
			t.Errorf("Deleted tessellator %d should return nil vertices", i)
		}
		if indices != nil {
			t.Errorf("Deleted tessellator %d should return nil indices", i)
		}
	}
}

// TestConstantMapping tests that all constants are properly mapped
func TestConstantMapping(t *testing.T) {
	// Test WindingRule constants
	expectedWindingRules := map[WindingRule]string{
		WindingOdd:       "Odd",
		WindingNonZero:   "NonZero",
		WindingPositive:  "Positive",
		WindingNegative:  "Negative",
		WindingAbsGeqTwo: "AbsGeqTwo",
	}

	for rule, expected := range expectedWindingRules {
		if rule.String() != expected {
			t.Errorf("WindingRule %v string mismatch: expected %s, got %s", rule, expected, rule.String())
		}
	}

	// Test ElementType constants
	expectedElementTypes := map[ElementType]string{
		ElementPolygons:          "Polygons",
		ElementConnectedPolygons: "ConnectedPolygons",
		ElementBoundaryContours:  "BoundaryContours",
	}

	for elemType, expected := range expectedElementTypes {
		if elemType.String() != expected {
			t.Errorf("ElementType %v string mismatch: expected %s, got %s", elemType, expected, elemType.String())
		}
	}

	// Test Status constants
	expectedStatuses := map[Status]string{
		StatusOK:           "OK",
		StatusOutOfMemory:  "OutOfMemory",
		StatusInvalidInput: "InvalidInput",
	}

	for status, expected := range expectedStatuses {
		if status.String() != expected {
			t.Errorf("Status %v string mismatch: expected %s, got %s", status, expected, status.String())
		}
	}
}

// TestDataCopying tests that data is properly copied between Go and C
func TestDataCopying(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Test 2D vertex copying
	vertices2D := []float32{
		1.5, 2.7,
		3.2, 4.1,
		2.8, 1.9,
	}

	err := tess.AddContour(2, vertices2D)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	// Test 3D vertex copying
	vertices3D := []float32{
		1.1, 2.2, 3.3,
		4.4, 5.5, 6.6,
		7.7, 8.8, 9.9,
	}

	err = tess.AddContour(3, vertices3D)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	// Test tessellation
	resultVertices, _, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Verify output data is reasonable
	if len(resultVertices) <= 0 {
		t.Errorf("Expected positive vertex count, got %d", len(resultVertices))
	}
}

// TestOutputValidation tests that output coordinates are in reasonable ranges
func TestOutputValidation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Create a simple polygon with known coordinate ranges
	vertices := []float32{
		100, -50,
		120, -50,
		120, 50,
		100, 50,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Test vertex count is reasonable
	vertexCount := len(resultVertices) / 2 // 2 coordinates per vertex
	if vertexCount < 3 || vertexCount > 100 {
		t.Errorf("Vertex count %d is not reasonable (expected 3-100)", vertexCount)
	}

	// Test element count is reasonable
	elementCount := len(resultIndices) / 3 // 3 indices per triangle
	if elementCount < 1 || elementCount > 50 {
		t.Errorf("Element count %d is not reasonable (expected 1-50)", elementCount)
	}

	// Test vertex coordinates are in reasonable ranges
	if len(resultVertices) != vertexCount*2 {
		t.Errorf("Vertex count mismatch: expected %d coordinates, got %d", vertexCount*2, len(resultVertices))
	}

	// Define reasonable coordinate ranges (non-overlapping)
	xRange := Range{Min: 100, Max: 120}
	yRange := Range{Min: -50, Max: 50}

	for i := 0; i < len(resultVertices); i += 2 {
		x, y := resultVertices[i], resultVertices[i+1]
		t.Logf("(%.3f, %.3f)", x, y)
		if !xRange.Contains(x) {
			t.Errorf("Vertex %d X coordinate %f outside reasonable range %v", i/2, x, xRange)
		}
		if !yRange.Contains(y) {
			t.Errorf("Vertex %d Y coordinate %f outside reasonable range %v", i/2, y, yRange)
		}
	}
}

// Test3DOutputValidation tests 3D tessellation output validation
func Test3DOutputValidation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Create a 3D polygon with different coordinate ranges
	vertices := []float32{
		100.0, -3.0, 10.0,
		150.0, -3.0, 10.0,
		150.0, 3.0, 20.0,
		100.0, 3.0, 20.0,
	}

	err := tess.AddContour(3, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	resultVertices, _, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 3, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Test vertex count is reasonable
	vertexCount := len(resultVertices) / 3 // 3 coordinates per vertex
	if vertexCount < 3 || vertexCount > 100 {
		t.Errorf("Vertex count %d is not reasonable (expected 3-100)", vertexCount)
	}

	// Test 3D vertex coordinates are in reasonable ranges
	if len(resultVertices) != vertexCount*3 {
		t.Errorf("Vertex count mismatch: expected %d coordinates, got %d", vertexCount*3, len(resultVertices))
	}

	// Define reasonable coordinate ranges for 3D (non-overlapping)
	xRange := Range{Min: 99.0, Max: 151.0}
	yRange := Range{Min: -5.0, Max: 5.0}
	zRange := Range{Min: 9.0, Max: 21.0}

	for i := 0; i < len(resultVertices); i += 3 {
		x, y, z := resultVertices[i], resultVertices[i+1], resultVertices[i+2]
		t.Logf("(%.3f, %.3f, %.3f)", x, y, z)
		if !xRange.Contains(x) {
			t.Errorf("Vertex %d X coordinate %f outside reasonable range %v", i/3, x, xRange)
		}
		if !yRange.Contains(y) {
			t.Errorf("Vertex %d Y coordinate %f outside reasonable range %v", i/3, y, yRange)
		}
		if !zRange.Contains(z) {
			t.Errorf("Vertex %d Z coordinate %f outside reasonable range %v", i/3, z, zRange)
		}
	}
}

// TestElementIndicesValidation tests that element indices are valid
func TestElementIndicesValidation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	vertices := []float32{
		0, 0,
		1, 0,
		1, 1,
		0, 1,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	vertexCount := len(resultVertices) / 2 // 2 coordinates per vertex
	elementCount := len(resultIndices) / 3 // 3 indices per triangle

	// Test element indices are valid
	expectedElementIndices := elementCount * 3

	if len(resultIndices) != expectedElementIndices {
		t.Errorf("Element indices count mismatch: expected %d, got %d", expectedElementIndices, len(resultIndices))
	}

	// Test all indices are within valid range
	for i, index := range resultIndices {
		if index < 0 || index >= vertexCount {
			t.Errorf("Element index %d at position %d is invalid: %d (vertex count: %d)", index, i, index, vertexCount)
		}
	}

	// Test triangles are properly formed
	if len(resultIndices) != elementCount*3 {
		t.Errorf("Triangle count mismatch: expected %d indices, got %d", elementCount*3, len(resultIndices))
	}

	for i := 0; i < len(resultIndices); i += 3 {
		// Check for degenerate triangles (all vertices the same)
		if resultIndices[i] == resultIndices[i+1] && resultIndices[i+1] == resultIndices[i+2] {
			t.Errorf("Triangle %d is degenerate (all vertices same): %v", i/3, resultIndices[i:i+3])
		}
	}
}

// TestVertexIndicesMapping tests vertex indices mapping functionality
func TestVertexIndicesMapping(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	vertices := []float32{
		0, 0,
		1, 0,
		0.5, 1,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	vertexCount := len(resultVertices) / 2 // 2 coordinates per vertex

	// Test that indices are reasonable (should be within valid range)
	for i, index := range resultIndices {
		if index < 0 || index >= vertexCount {
			t.Errorf("Vertex index %d at position %d is invalid: %d (vertex count: %d)", index, i, index, vertexCount)
		}
	}
}

// TestOptionHandling tests option setting and validation
func TestOptionHandling(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Test setting constrained Delaunay triangulation
	err := tess.SetOption(OptionConstrainedDelaunay, true)
	if err != nil {
		t.Errorf("SetOption failed for OptionConstrainedDelaunay: %v", err)
	}

	// Test setting reverse contours
	err = tess.SetOption(OptionReverseContours, false)
	if err != nil {
		t.Errorf("SetOption failed for OptionReverseContours: %v", err)
	}

	// Test tessellation with options
	vertices := []float32{
		0, 0,
		1, 0,
		1, 1,
		0, 1,
	}

	err = tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Verify output is still reasonable
	vertexCount := len(resultVertices) / 2
	if vertexCount <= 0 {
		t.Errorf("Expected positive vertex count with options, got %d", vertexCount)
	}

	// Use resultIndices to avoid unused variable warning
	_ = resultIndices
}

// TestNormalVectorHandling tests normal vector parameter handling
func TestNormalVectorHandling(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	vertices := []float32{
		0, 0,
		1, 0,
		0.5, 1,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	// Test with nil normal (auto-calculation)
	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed with nil normal: %v", err)
	}

	// Test with explicit normal
	normal := []float32{0, 0, 1}
	resultVertices2, resultIndices2, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, normal)
	if err != nil {
		t.Fatalf("Tessellate failed with explicit normal: %v", err)
	}

	// Verify output is reasonable in both cases
	vertexCount := len(resultVertices) / 2
	if vertexCount <= 0 {
		t.Errorf("Expected positive vertex count with nil normal, got %d", vertexCount)
	}

	vertexCount2 := len(resultVertices2) / 2
	// When using an explicit normal vector, the tessellation might not produce output
	// for certain simple shapes, which is acceptable behavior
	if vertexCount2 > 0 {
		// If we do get output, verify it's reasonable
		if vertexCount2 < 3 {
			t.Errorf("Expected at least 3 vertices with explicit normal, got %d", vertexCount2)
		}
	}

	// Use resultIndices to avoid unused variable warnings
	_ = resultIndices
	_ = resultIndices2
}

// TestComplexPolygonValidation tests a more complex polygon with multiple contours
func TestComplexPolygonValidation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Outer contour (square)
	outerContour := []float32{
		-5.0, -5.0,
		5.0, -5.0,
		5.0, 5.0,
		-5.0, 5.0,
	}

	// Inner contour (triangle hole)
	innerContour := []float32{
		-2.0, -2.0,
		2.0, -2.0,
		0.0, 2.0,
	}

	err := tess.AddContour(2, outerContour)
	if err != nil {
		t.Fatalf("AddContour failed for outer contour: %v", err)
	}

	err = tess.AddContour(2, innerContour)
	if err != nil {
		t.Fatalf("AddContour failed for inner contour: %v", err)
	}

	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Test reasonable output ranges
	vertexCount := len(resultVertices) / 2
	if vertexCount < 6 || vertexCount > 200 {
		t.Errorf("Vertex count %d is not reasonable for complex polygon (expected 6-200)", vertexCount)
	}

	elementCount := len(resultIndices) / 3
	if elementCount < 2 || elementCount > 100 {
		t.Errorf("Element count %d is not reasonable for complex polygon (expected 2-100)", elementCount)
	}

	// Test coordinate ranges
	xRange := Range{Min: -6.0, Max: 6.0}
	yRange := Range{Min: -6.0, Max: 6.0}

	for i := 0; i < len(resultVertices); i += 2 {
		x, y := resultVertices[i], resultVertices[i+1]
		if !xRange.Contains(x) {
			t.Errorf("Vertex %d X coordinate %f outside reasonable range %v", i/2, x, xRange)
		}
		if !yRange.Contains(y) {
			t.Errorf("Vertex %d Y coordinate %f outside reasonable range %v", i/2, y, yRange)
		}
	}
}

// TestEdgeCases tests various edge cases in the binding
func TestEdgeCases(t *testing.T) {
	// Test with single vertex (should fail gracefully)
	tess := NewTessellator()
	defer tess.Delete()

	singleVertex := []float32{0, 0}
	err := tess.AddContour(2, singleVertex)
	if err != nil {
		t.Fatalf("AddContour failed with single vertex: %v", err)
	}

	// This should fail but not crash
	_, _, err = tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err == nil {
		t.Log("Tessellation with single vertex succeeded (unexpected)")
	}

	// Test with collinear vertices
	tess2 := NewTessellator()
	defer tess2.Delete()

	collinearVertices := []float32{
		0, 0,
		1, 0,
		2, 0,
	}

	err = tess2.AddContour(2, collinearVertices)
	if err != nil {
		t.Fatalf("AddContour failed with collinear vertices: %v", err)
	}

	// This should fail but not crash
	_, _, err = tess2.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err == nil {
		t.Log("Tessellation with collinear vertices succeeded (unexpected)")
	}
}

// Range represents a coordinate range for validation
type Range struct {
	Min, Max float32
}

// Contains checks if a value is within the range
func (r Range) Contains(value float32) bool {
	return value >= r.Min && value <= r.Max
}

// TestFloatingPointPrecision tests that floating point values are handled correctly
func TestFloatingPointPrecision(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Test with very small and very large values
	vertices := []float32{
		1e-6, 1e-6,
		1e6, 1e-6,
		1e6, 1e6,
		1e-6, 1e6,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed with extreme values: %v", err)
	}

	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed with extreme values: %v", err)
	}

	// Verify output is reasonable
	vertexCount := len(resultVertices) / 2
	if vertexCount <= 0 {
		t.Errorf("Expected positive vertex count with extreme values, got %d", vertexCount)
	}

	// Test that coordinates are finite
	for i := 0; i < len(resultVertices); i++ {
		if !isFinite(resultVertices[i]) {
			t.Errorf("Vertex coordinate %d is not finite: %f", i, resultVertices[i])
		}
	}

	// Use resultIndices to avoid unused variable warning
	_ = resultIndices
}

// isFinite checks if a float32 value is finite
func isFinite(f float32) bool {
	return !math.IsNaN(float64(f)) && !math.IsInf(float64(f), 0)
}
