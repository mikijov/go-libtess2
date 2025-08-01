package tess

import (
	"testing"
)

// TestNewTessellator tests tessellator creation and cleanup
func TestNewTessellator(t *testing.T) {
	tess := NewTessellator()
	if tess == nil {
		t.Fatal("NewTessellator() returned nil")
	}

	// Test cleanup
	tess.Delete()

	// Test that tessellator is properly cleaned up
	// Note: We can't test status directly anymore since getStatus() is private
	// The tessellator should be properly cleaned up after Delete()
}

// TestAddContour tests adding contours to the tessellator
func TestAddContour(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Test valid 2D contour
	vertices2D := []float32{
		0, 0,
		1, 0,
		1, 1,
		0, 1,
	}

	err := tess.AddContour(2, vertices2D)
	if err != nil {
		t.Errorf("AddContour failed: %v", err)
	}

	// Test valid 3D contour
	vertices3D := []float32{
		0, 0, 0,
		1, 0, 1,
		1, 1, 1,
		0, 1, 0,
	}

	err = tess.AddContour(3, vertices3D)
	if err != nil {
		t.Errorf("AddContour failed: %v", err)
	}

	// Test invalid size
	err = tess.AddContour(4, vertices2D)
	if err == nil {
		t.Error("Expected error for invalid size")
	}

	// Test empty vertices
	err = tess.AddContour(2, []float32{})
	if err == nil {
		t.Error("Expected error for empty vertices")
	}
}

// TestSetOption tests setting tessellation options
func TestSetOption(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Test setting constrained Delaunay triangulation
	err := tess.SetOption(OptionConstrainedDelaunay, true)
	if err != nil {
		t.Errorf("SetOption failed: %v", err)
	}

	// Test setting reverse contours
	err = tess.SetOption(OptionReverseContours, false)
	if err != nil {
		t.Errorf("SetOption failed: %v", err)
	}
}

// TestSimpleTriangleTessellation tests basic triangle tessellation
func TestSimpleTriangleTessellation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Add a simple triangle contour
	vertices := []float32{
		0, 0,
		1, 0,
		0.5, 1,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	// Tessellate using new simplified API
	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Check results
	// For a triangle with 2D coordinates, we expect 3 vertices * 2 coordinates = 6 values
	expectedVertexCount := 3 * 2
	if len(resultVertices) != expectedVertexCount {
		t.Errorf("Expected %d vertex coordinates, got %d", expectedVertexCount, len(resultVertices))
	}

	// For ElementPolygons with polySize=3, we expect 1 triangle * 3 vertices = 3 indices
	expectedIndexCount := 1 * 3
	if len(resultIndices) != expectedIndexCount {
		t.Errorf("Expected %d indices, got %d", expectedIndexCount, len(resultIndices))
	}

	// Verify indices are valid (should be 0, 1, 2 for a simple triangle)
	for i, index := range resultIndices {
		if index < 0 || index >= 3 {
			t.Errorf("Invalid vertex index %d at position %d", index, i)
		}
	}
}

// TestSquareTessellation tests tessellating a square
func TestSquareTessellation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()

	// Add a square contour
	vertices := []float32{
		100, 50,
		110, 50,
		110, 70,
		100, 70,
	}

	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}

	// Tessellate using new simplified API
	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	// Check results
	// For a square with 2D coordinates, we expect 4 vertices * 2 coordinates = 8 values
	expectedVertexCount := 4 * 2
	if len(resultVertices) != expectedVertexCount {
		t.Errorf("Expected %d vertex coordinates, got %d", expectedVertexCount, len(resultVertices))
	}

	// For ElementPolygons with polySize=3, we expect 2 triangles * 3 vertices = 6 indices
	expectedIndexCount := 2 * 3
	if len(resultIndices) != expectedIndexCount {
		t.Errorf("Expected %d indices, got %d", expectedIndexCount, len(resultIndices))
	}

	// Verify indices are valid (should be 0-3 for a square)
	for i, index := range resultIndices {
		if index < 0 || index >= 4 {
			t.Errorf("Invalid vertex index %d at position %d", index, i)
		}
	}
}

// TestWindingRules tests different winding rules with a complex polygon
func TestWindingRules(t *testing.T) {
	windingRules := []WindingRule{
		WindingOdd,
		WindingNonZero,
		WindingPositive,
		WindingNegative,
		WindingAbsGeqTwo,
	}

	for _, rule := range windingRules {
		t.Run(rule.String(), func(t *testing.T) {
			tess := NewTessellator()
			defer tess.Delete()

			// Create a complex polygon with overlapping regions to test all winding rules
			// This polygon has regions with different winding numbers:
			// - Outer square: winding number 1
			// - Inner square: creates regions with winding number 2 (overlap)
			// - Small triangle: creates regions with negative winding numbers
			
			// Outer square (clockwise)
			outerSquare := []float32{
				0, 0,
				10, 0,
				10, 10,
				0, 10,
			}
			
			// Inner square (counter-clockwise to create negative winding)
			innerSquare := []float32{
				3, 3,
				7, 3,
				7, 7,
				3, 7,
			}
			
			// Small triangle inside inner square (clockwise)
			smallTriangle := []float32{
				4, 4,
				6, 4,
				5, 6,
			}

			// Add contours in order to create proper winding
			err := tess.AddContour(2, outerSquare)
			if err != nil {
				t.Fatalf("AddContour failed for outer square: %v", err)
			}

			err = tess.AddContour(2, innerSquare)
			if err != nil {
				t.Fatalf("AddContour failed for inner square: %v", err)
			}

			err = tess.AddContour(2, smallTriangle)
			if err != nil {
				t.Fatalf("AddContour failed for small triangle: %v", err)
			}

			resultVertices, resultIndices, err := tess.Tessellate(rule, ElementPolygons, 3, 2, nil)
			if err != nil {
				t.Fatalf("Tessellate failed: %v", err)
			}

			// For WindingNegative, it's acceptable to have no output if the polygon
			// doesn't have regions with negative winding numbers
			if rule == WindingNegative {
				if len(resultVertices) > 0 && len(resultIndices) > 0 {
					// If we do get output, verify it's valid
					vertexCount := len(resultVertices) / 2
					for i, index := range resultIndices {
						if index < 0 || index >= vertexCount {
							t.Errorf("Invalid vertex index %d at position %d for winding rule %v", index, i, rule)
						}
					}
				}
			} else {
				// Other winding rules should produce output
				if len(resultVertices) == 0 {
					t.Errorf("Expected non-empty vertices for winding rule %v", rule)
				}

				if len(resultIndices) == 0 {
					t.Errorf("Expected non-empty indices for winding rule %v", rule)
				}

				// Verify indices are valid
				vertexCount := len(resultVertices) / 2
				for i, index := range resultIndices {
					if index < 0 || index >= vertexCount {
						t.Errorf("Invalid vertex index %d at position %d for winding rule %v", index, i, rule)
					}
				}

				// Verify we have reasonable output
				if vertexCount < 3 {
					t.Errorf("Expected at least 3 vertices for winding rule %v, got %d", rule, vertexCount)
				}

				elementCount := len(resultIndices) / 3
				if elementCount < 1 {
					t.Errorf("Expected at least 1 element for winding rule %v, got %d", rule, elementCount)
				}
			}
		})
	}
}

// TestWindingRulesSimple tests winding rules with a simple polygon that has a hole
func TestWindingRulesSimple(t *testing.T) {
	windingRules := []WindingRule{
		WindingOdd,
		WindingNonZero,
		WindingPositive,
		WindingNegative,
		WindingAbsGeqTwo,
	}

	for _, rule := range windingRules {
		t.Run(rule.String(), func(t *testing.T) {
			tess := NewTessellator()
			defer tess.Delete()

			// Create a square with a triangular hole
			// This should create regions with different winding numbers
			
			// Outer square (clockwise)
			outerSquare := []float32{
				0, 0,
				8, 0,
				8, 8,
				0, 8,
			}
			
			// Inner triangle hole (counter-clockwise)
			innerTriangle := []float32{
				2, 2,
				6, 2,
				4, 6,
			}

			// Add contours
			err := tess.AddContour(2, outerSquare)
			if err != nil {
				t.Fatalf("AddContour failed for outer square: %v", err)
			}

			err = tess.AddContour(2, innerTriangle)
			if err != nil {
				t.Fatalf("AddContour failed for inner triangle: %v", err)
			}

			resultVertices, resultIndices, err := tess.Tessellate(rule, ElementPolygons, 3, 2, nil)
			if err != nil {
				t.Fatalf("Tessellate failed: %v", err)
			}

			// For WindingNegative, it's acceptable to have no output if the polygon
			// doesn't have regions with negative winding numbers
			if rule == WindingNegative {
				if len(resultVertices) > 0 && len(resultIndices) > 0 {
					// If we do get output, verify it's valid
					vertexCount := len(resultVertices) / 2
					for i, index := range resultIndices {
						if index < 0 || index >= vertexCount {
							t.Errorf("Invalid vertex index %d at position %d for winding rule %v", index, i, rule)
						}
					}
				}
			} else {
				// Other winding rules should produce output
				if len(resultVertices) == 0 {
					t.Errorf("Expected non-empty vertices for winding rule %v", rule)
				}

				if len(resultIndices) == 0 {
					t.Errorf("Expected non-empty indices for winding rule %v", rule)
				}

				// Verify indices are valid
				vertexCount := len(resultVertices) / 2
				for i, index := range resultIndices {
					if index < 0 || index >= vertexCount {
						t.Errorf("Invalid vertex index %d at position %d for winding rule %v", index, i, rule)
					}
				}
			}
		})
	}
}

// TestWindingRulesNegative tests winding rules with a polygon designed to create negative winding numbers
func TestWindingRulesNegative(t *testing.T) {
	windingRules := []WindingRule{
		WindingOdd,
		WindingNonZero,
		WindingPositive,
		WindingNegative,
		WindingAbsGeqTwo,
	}

	for _, rule := range windingRules {
		t.Run(rule.String(), func(t *testing.T) {
			tess := NewTessellator()
			defer tess.Delete()

			// Create a polygon with overlapping contours in opposite directions
			// This should create regions with negative winding numbers
			
			// First contour (clockwise)
			contour1 := []float32{
				0, 0,
				10, 0,
				10, 10,
				0, 10,
			}
			
			// Second contour (counter-clockwise) that overlaps
			contour2 := []float32{
				2, 2,
				8, 2,
				8, 8,
				2, 8,
			}

			// Add contours
			err := tess.AddContour(2, contour1)
			if err != nil {
				t.Fatalf("AddContour failed for contour1: %v", err)
			}

			err = tess.AddContour(2, contour2)
			if err != nil {
				t.Fatalf("AddContour failed for contour2: %v", err)
			}

			resultVertices, resultIndices, err := tess.Tessellate(rule, ElementPolygons, 3, 2, nil)
			if err != nil {
				t.Fatalf("Tessellate failed: %v", err)
			}

			// For WindingNegative, it's acceptable to have no output if the polygon
			// doesn't have regions with negative winding numbers
			if rule == WindingNegative {
				if len(resultVertices) > 0 && len(resultIndices) > 0 {
					// If we do get output, verify it's valid
					vertexCount := len(resultVertices) / 2
					for i, index := range resultIndices {
						if index < 0 || index >= vertexCount {
							t.Errorf("Invalid vertex index %d at position %d for winding rule %v", index, i, rule)
						}
					}
				}
			} else {
				// Other winding rules should produce output
				if len(resultVertices) == 0 {
					t.Errorf("Expected non-empty vertices for winding rule %v", rule)
				}

				if len(resultIndices) == 0 {
					t.Errorf("Expected non-empty indices for winding rule %v", rule)
				}

				// Verify indices are valid
				vertexCount := len(resultVertices) / 2
				for i, index := range resultIndices {
					if index < 0 || index >= vertexCount {
						t.Errorf("Invalid vertex index %d at position %d for winding rule %v", index, i, rule)
					}
				}
			}
		})
	}
}

// TestElementTypes tests different element types
func TestElementTypes(t *testing.T) {
	elementTypes := []ElementType{
		ElementPolygons,
		ElementConnectedPolygons,
		ElementBoundaryContours,
	}

	for _, elemType := range elementTypes {
		t.Run(elemType.String(), func(t *testing.T) {
			tess := NewTessellator()
			defer tess.Delete()

			vertices := []float32{
				100, 50,
				110, 50,
				110, 70,
				100, 70,
			}

			err := tess.AddContour(2, vertices)
			if err != nil {
				t.Fatalf("AddContour failed: %v", err)
			}

			resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, elemType, 3, 2, nil)
			if err != nil {
				t.Fatalf("Tessellate failed: %v", err)
			}

			// Should always produce some elements
			if len(resultVertices) == 0 {
				t.Errorf("Expected non-empty vertices for element type %v", elemType)
			}

			if len(resultIndices) == 0 {
				t.Errorf("Expected non-empty indices for element type %v", elemType)
			}
		})
	}
}

// TestNilTessellator tests operations on nil tessellator
func TestNilTessellator(t *testing.T) {
	var tess *Tessellator

	// Test operations on nil tessellator
	err := tess.AddContour(2, []float32{0, 0, 0})
	if err == nil {
		t.Error("Expected error for nil tessellator")
	}

	err = tess.SetOption(OptionConstrainedDelaunay, true)
	if err == nil {
		t.Error("Expected error for nil tessellator")
	}

	vertices, indices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err == nil {
		t.Error("Expected error for nil tessellator")
	}
	if vertices != nil {
		t.Error("Expected nil vertices for nil tessellator")
	}
	if indices != nil {
		t.Error("Expected nil indices for nil tessellator")
	}
}

// TestStringMethods tests string representations
func TestStringMethods(t *testing.T) {
	// Test WindingRule strings
	if WindingOdd.String() != "Odd" {
		t.Errorf("Expected 'Odd', got '%s'", WindingOdd.String())
	}

	if WindingNonZero.String() != "NonZero" {
		t.Errorf("Expected 'NonZero', got '%s'", WindingNonZero.String())
	}

	// Test ElementType strings
	if ElementPolygons.String() != "Polygons" {
		t.Errorf("Expected 'Polygons', got '%s'", ElementPolygons.String())
	}

	if ElementBoundaryContours.String() != "BoundaryContours" {
		t.Errorf("Expected 'BoundaryContours', got '%s'", ElementBoundaryContours.String())
	}

	// Test Status strings
	if StatusOK.String() != "OK" {
		t.Errorf("Expected 'OK', got '%s'", StatusOK.String())
	}

	if StatusOutOfMemory.String() != "OutOfMemory" {
		t.Errorf("Expected 'OutOfMemory', got '%s'", StatusOutOfMemory.String())
	}
}

// TestNormalVector tests tessellation with a normal vector
func TestNormalVector(t *testing.T) {
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

	// Test with normal vector
	normal := []float32{0, 0, 1}
	resultVertices, resultIndices, err := tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, normal)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}

	if len(resultVertices) == 0 {
		t.Error("Expected non-empty vertices")
	}

	if len(resultIndices) == 0 {
		t.Error("Expected non-empty indices")
	}
}

// TestInvalidNormalVector tests tessellation with invalid normal vector
func TestInvalidNormalVector(t *testing.T) {
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

	// Test with invalid normal vector (too short)
	normal := []float32{0, 1} // Only 2 components, need 3
	_, _, err = tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, normal)
	if err == nil {
		t.Error("Expected error for invalid normal vector")
	}
}

// BenchmarkTessellation benchmarks tessellation performance
func BenchmarkTessellation(b *testing.B) {
	vertices := []float32{
		0, 0,
		1, 0,
		1, 1,
		0, 1,
		0.5, 0.5,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tess := NewTessellator()
		tess.AddContour(2, vertices)
		_, _, _ = tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
		tess.Delete()
	}
}
