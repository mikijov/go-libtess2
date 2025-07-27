package main

import (
	"fmt"
	"log"

	tess "github.com/miki/go-libtess2"
)

func main() {
	fmt.Println("Complex Polygon Tessellation Example")
	fmt.Println("====================================")

	// Create a new tessellator
	tessellator := tess.NewTessellator()
	if tessellator == nil {
		log.Fatal("Failed to create tessellator")
	}
	defer tessellator.Delete()

	// Define a square with a triangular hole
	// Outer contour (clockwise)
	outerContour := []tess.Vertex{
		{X: 0, Y: 0, Z: 0},   // Bottom-left
		{X: 4, Y: 0, Z: 0},   // Bottom-right
		{X: 4, Y: 4, Z: 0},   // Top-right
		{X: 0, Y: 4, Z: 0},   // Top-left
	}

	// Inner contour (hole, counter-clockwise)
	innerContour := []tess.Vertex{
		{X: 1, Y: 1, Z: 0},   // Bottom-left of hole
		{X: 3, Y: 1, Z: 0},   // Bottom-right of hole
		{X: 2, Y: 3, Z: 0},   // Top of hole
	}

	fmt.Printf("Outer contour: %v\n", outerContour)
	fmt.Printf("Inner contour (hole): %v\n", innerContour)

	// Add the outer contour
	err := tessellator.AddContour(2, outerContour)
	if err != nil {
		log.Fatalf("Failed to add outer contour: %v", err)
	}

	// Add the inner contour (hole)
	err = tessellator.AddContour(2, innerContour)
	if err != nil {
		log.Fatalf("Failed to add inner contour: %v", err)
	}

	// Enable constrained Delaunay triangulation for better quality
	err = tessellator.SetOption(tess.OptionConstrainedDelaunay, true)
	if err != nil {
		log.Fatalf("Failed to set option: %v", err)
	}

	// Perform tessellation
	err = tessellator.Tessellate(tess.WindingOdd, tess.ElementPolygons, 3, 2, nil)
	if err != nil {
		log.Fatalf("Tessellation failed: %v", err)
	}

	// Get results
	vertexCount := tessellator.GetVertexCount()
	elementCount := tessellator.GetElementCount()
	outputVertices := tessellator.GetVertices()
	elements := tessellator.GetElements()
	indices := tessellator.GetVertexIndices()

	fmt.Printf("\nTessellation Results:\n")
	fmt.Printf("Vertex count: %d\n", vertexCount)
	fmt.Printf("Element count: %d\n", elementCount)
	fmt.Printf("Output vertices: %v\n", outputVertices)
	fmt.Printf("Elements: %v\n", elements)
	fmt.Printf("Vertex indices: %v\n", indices)

	// Print triangle information
	fmt.Printf("\nTriangle Details:\n")
	for i := 0; i < elementCount; i++ {
		base := i * 3
		v1 := elements[base]
		v2 := elements[base+1]
		v3 := elements[base+2]
		fmt.Printf("Triangle %d: vertices [%d, %d, %d] -> positions [%v, %v, %v]\n",
			i, v1, v2, v3,
			outputVertices[v1], outputVertices[v2], outputVertices[v3])
	}

	// Demonstrate different winding rules
	fmt.Printf("\nWinding Rule Comparison:\n")
	windingRules := []tess.WindingRule{
		tess.WindingOdd,
		tess.WindingNonZero,
		tess.WindingPositive,
		tess.WindingNegative,
		tess.WindingAbsGeqTwo,
	}

	for _, rule := range windingRules {
		// Create a new tessellator for each test
		testTess := tess.NewTessellator()
		if testTess == nil {
			continue
		}

		testTess.AddContour(2, outerContour)
		testTess.AddContour(2, innerContour)
		testTess.SetOption(tess.OptionConstrainedDelaunay, true)

		err := testTess.Tessellate(rule, tess.ElementPolygons, 3, 2, nil)
		if err == nil {
			fmt.Printf("  %s: %d triangles\n", rule.String(), testTess.GetElementCount())
		} else {
			fmt.Printf("  %s: failed (%v)\n", rule.String(), err)
		}

		testTess.Delete()
	}

	// Demonstrate different element types
	fmt.Printf("\nElement Type Comparison:\n")
	elementTypes := []tess.ElementType{
		tess.ElementPolygons,
		tess.ElementConnectedPolygons,
		tess.ElementBoundaryContours,
	}

	for _, elemType := range elementTypes {
		// Create a new tessellator for each test
		testTess := tess.NewTessellator()
		if testTess == nil {
			continue
		}

		testTess.AddContour(2, outerContour)
		testTess.AddContour(2, innerContour)

		err := testTess.Tessellate(tess.WindingOdd, elemType, 3, 2, nil)
		if err == nil {
			fmt.Printf("  %s: %d elements\n", elemType.String(), testTess.GetElementCount())
		} else {
			fmt.Printf("  %s: failed (%v)\n", elemType.String(), err)
		}

		testTess.Delete()
	}
} 