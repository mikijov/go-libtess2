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
	
	// Test that we can get status
	status := tess.GetStatus()
	if status != StatusOK {
		t.Errorf("Expected status OK, got %v", status)
	}
	
	// Test cleanup
	tess.Delete()
	
	// Test that tessellator is properly cleaned up
	status = tess.GetStatus()
	if status != StatusInvalidInput {
		t.Errorf("Expected status InvalidInput after deletion, got %v", status)
	}
}

// TestAddContour tests adding contours to the tessellator
func TestAddContour(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()
	
	// Test valid 2D contour
	vertices2D := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}
	
	err := tess.AddContour(2, vertices2D)
	if err != nil {
		t.Errorf("AddContour failed: %v", err)
	}
	
	// Test valid 3D contour
	vertices3D := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 1},
		{X: 1, Y: 1, Z: 1},
		{X: 0, Y: 1, Z: 0},
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
	err = tess.AddContour(2, []Vertex{})
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
	vertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 0.5, Y: 1, Z: 0},
	}
	
	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}
	
	// Tessellate
	err = tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}
	
	// Check results
	vertexCount := tess.GetVertexCount()
	if vertexCount != 3 {
		t.Errorf("Expected 3 vertices, got %d", vertexCount)
	}
	
	elementCount := tess.GetElementCount()
	if elementCount != 1 {
		t.Errorf("Expected 1 element, got %d", elementCount)
	}
	
	vertices = tess.GetVertices()
	if len(vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(vertices))
	}
	
	// Test GetElementsWithSize
	elements := tess.GetElementsWithSize(ElementPolygons, 3)
	if len(elements) != 3 {
		t.Errorf("Expected 3 element indices, got %d", len(elements))
	}
	
	// Test GetTriangles
	triangles := tess.GetTriangles()
	if len(triangles) != 1 {
		t.Errorf("Expected 1 triangle, got %d", len(triangles))
	}
	if len(triangles[0]) != 3 {
		t.Errorf("Expected triangle to have 3 vertices, got %d", len(triangles[0]))
	}
	
	indices := tess.GetVertexIndices()
	if len(indices) != 3 {
		t.Errorf("Expected 3 vertex indices, got %d", len(indices))
	}
}

// TestSquareTessellation tests tessellating a square
func TestSquareTessellation(t *testing.T) {
	tess := NewTessellator()
	defer tess.Delete()
	
	// Add a square contour
	vertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}
	
	err := tess.AddContour(2, vertices)
	if err != nil {
		t.Fatalf("AddContour failed: %v", err)
	}
	
	// Tessellate
	err = tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err != nil {
		t.Fatalf("Tessellate failed: %v", err)
	}
	
	// Check results
	vertexCount := tess.GetVertexCount()
	if vertexCount != 4 {
		t.Errorf("Expected 4 vertices, got %d", vertexCount)
	}
	
	elementCount := tess.GetElementCount()
	if elementCount != 2 {
		t.Errorf("Expected 2 elements (triangles), got %d", elementCount)
	}
	
	// Test GetElementsWithSize
	elements := tess.GetElementsWithSize(ElementPolygons, 3)
	expectedElementIndices := elementCount * 3
	if len(elements) != expectedElementIndices {
		t.Errorf("Expected %d element indices (2 triangles * 3 vertices), got %d", expectedElementIndices, len(elements))
	}
	
	// Test GetTriangles
	triangles := tess.GetTriangles()
	if len(triangles) != 2 {
		t.Errorf("Expected 2 triangles, got %d", len(triangles))
	}
	for i, triangle := range triangles {
		if len(triangle) != 3 {
			t.Errorf("Expected triangle %d to have 3 vertices, got %d", i, len(triangle))
		}
	}
}

// TestWindingRules tests different winding rules
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
			
			vertices := []Vertex{
				{X: 0, Y: 0, Z: 0},
				{X: 1, Y: 0, Z: 0},
				{X: 0.5, Y: 1, Z: 0},
			}
			
			err := tess.AddContour(2, vertices)
			if err != nil {
				t.Fatalf("AddContour failed: %v", err)
			}
			
			err = tess.Tessellate(rule, ElementPolygons, 3, 2, nil)
			if err != nil {
				t.Fatalf("Tessellate failed: %v", err)
			}
			
			// Some winding rules may not produce elements for simple shapes
			// This is expected behavior
			elementCount := tess.GetElementCount()
			if elementCount < 0 {
				t.Errorf("Expected non-negative element count, got %d", elementCount)
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
			
			vertices := []Vertex{
				{X: 0, Y: 0, Z: 0},
				{X: 1, Y: 0, Z: 0},
				{X: 1, Y: 1, Z: 0},
				{X: 0, Y: 1, Z: 0},
			}
			
			err := tess.AddContour(2, vertices)
			if err != nil {
				t.Fatalf("AddContour failed: %v", err)
			}
			
			err = tess.Tessellate(WindingOdd, elemType, 3, 2, nil)
			if err != nil {
				t.Fatalf("Tessellate failed: %v", err)
			}
			
			// Should always produce some elements
			elementCount := tess.GetElementCount()
			if elementCount < 1 {
				t.Errorf("Expected at least 1 element, got %d", elementCount)
			}
		})
	}
}

// TestNilTessellator tests operations on nil tessellator
func TestNilTessellator(t *testing.T) {
	var tess *Tessellator
	
	// Test operations on nil tessellator
	err := tess.AddContour(2, []Vertex{{X: 0, Y: 0, Z: 0}})
	if err == nil {
		t.Error("Expected error for nil tessellator")
	}
	
	err = tess.SetOption(OptionConstrainedDelaunay, true)
	if err == nil {
		t.Error("Expected error for nil tessellator")
	}
	
	err = tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
	if err == nil {
		t.Error("Expected error for nil tessellator")
	}
	
	vertexCount := tess.GetVertexCount()
	if vertexCount != 0 {
		t.Errorf("Expected 0 vertices for nil tessellator, got %d", vertexCount)
	}
	
	elementCount := tess.GetElementCount()
	if elementCount != 0 {
		t.Errorf("Expected 0 elements for nil tessellator, got %d", elementCount)
	}
	
	vertices := tess.GetVertices()
	if vertices != nil {
		t.Error("Expected nil vertices for nil tessellator")
	}
	
	elements := tess.GetElements()
	if elements != nil {
		t.Error("Expected nil elements for nil tessellator")
	}
	
	indices := tess.GetVertexIndices()
	if indices != nil {
		t.Error("Expected nil indices for nil tessellator")
	}
	
	status := tess.GetStatus()
	if status != StatusInvalidInput {
		t.Errorf("Expected StatusInvalidInput for nil tessellator, got %v", status)
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

// BenchmarkTessellation benchmarks tessellation performance
func BenchmarkTessellation(b *testing.B) {
	vertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 0},
		{X: 0, Y: 1, Z: 0},
		{X: 0.5, Y: 0.5, Z: 0},
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		tess := NewTessellator()
		tess.AddContour(2, vertices)
		tess.Tessellate(WindingOdd, ElementPolygons, 3, 2, nil)
		tess.Delete()
	}
} 