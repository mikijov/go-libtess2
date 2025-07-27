package main

import (
	"fmt"
	"log"

	tess "github.com/miki/go-libtess2"
)

func main() {
	fmt.Println("Simple Triangle Tessellation Example")
	fmt.Println("===================================")

	// Create a new tessellator
	tessellator := tess.NewTessellator()
	if tessellator == nil {
		log.Fatal("Failed to create tessellator")
	}
	defer tessellator.Delete()

	// Define a simple triangle
	vertices := []tess.Vertex{
		{X: 0, Y: 0, Z: 0},   // Bottom-left
		{X: 1, Y: 0, Z: 0},   // Bottom-right
		{X: 0.5, Y: 1, Z: 0}, // Top
	}

	fmt.Printf("Input vertices: %v\n", vertices)

	// Add the contour
	err := tessellator.AddContour(2, vertices)
	if err != nil {
		log.Fatalf("Failed to add contour: %v", err)
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
} 