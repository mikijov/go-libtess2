// Package tess provides Go bindings for the libtess2 polygon tessellation library.
// libtess2 is a high-quality polygon tessellator and triangulator library.
//
// This package provides an idiomatic Go interface to the C library, allowing
// you to tessellate complex polygons into triangles, contours, and other
// geometric primitives.
package tess

/*
#cgo CFLAGS: -I./libtess2/Include
#cgo LDFLAGS: -L./libtess2 -ltess2
#include "tesselator.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// WindingRule defines the winding rule used for tessellation.
type WindingRule int

const (
	WindingOdd        WindingRule = C.TESS_WINDING_ODD
	WindingNonZero    WindingRule = C.TESS_WINDING_NONZERO
	WindingPositive   WindingRule = C.TESS_WINDING_POSITIVE
	WindingNegative   WindingRule = C.TESS_WINDING_NEGATIVE
	WindingAbsGeqTwo  WindingRule = C.TESS_WINDING_ABS_GEQ_TWO
)

// ElementType defines the type of output elements from tessellation.
type ElementType int

const (
	ElementPolygons           ElementType = C.TESS_POLYGONS
	ElementConnectedPolygons  ElementType = C.TESS_CONNECTED_POLYGONS
	ElementBoundaryContours   ElementType = C.TESS_BOUNDARY_CONTOURS
)

// Option defines tessellation options.
type Option int

const (
	OptionConstrainedDelaunay Option = C.TESS_CONSTRAINED_DELAUNAY_TRIANGULATION
	OptionReverseContours     Option = C.TESS_REVERSE_CONTOURS
)

// Status represents the tessellation status.
type Status int

const (
	StatusOK             Status = C.TESS_STATUS_OK
	StatusOutOfMemory    Status = C.TESS_STATUS_OUT_OF_MEMORY
	StatusInvalidInput   Status = C.TESS_STATUS_INVALID_INPUT
)

// Vertex represents a 2D or 3D vertex.
type Vertex struct {
	X, Y, Z float32
}

// Tessellator represents a tessellation context.
type Tessellator struct {
	tess *C.TESStesselator
}

// NewTessellator creates a new tessellator instance.
// Returns nil if allocation fails.
func NewTessellator() *Tessellator {
	tess := C.tessNewTess(nil)
	if tess == nil {
		return nil
	}
	
	t := &Tessellator{tess: tess}
	runtime.SetFinalizer(t, (*Tessellator).Delete)
	return t
}

// Delete destroys the tessellator and frees associated memory.
func (t *Tessellator) Delete() {
	if t.tess != nil {
		C.tessDeleteTess(t.tess)
		t.tess = nil
	}
}

// AddContour adds a contour to be tessellated.
// size must be 2 or 3 (for 2D or 3D vertices).
// vertices is a slice of vertices forming the contour.
func (t *Tessellator) AddContour(size int, vertices []Vertex) error {
	if t == nil || t.tess == nil {
		return fmt.Errorf("tessellator is nil or deleted")
	}
	
	if size != 2 && size != 3 {
		return fmt.Errorf("size must be 2 or 3, got %d", size)
	}
	
	if len(vertices) == 0 {
		return fmt.Errorf("vertices slice cannot be empty")
	}
	
	// Convert vertices to C array using type-safe approach
	cVerts := make([]C.float, len(vertices)*size)
	for i, v := range vertices {
		cVerts[i*size] = C.float(v.X)
		cVerts[i*size+1] = C.float(v.Y)
		if size == 3 {
			cVerts[i*size+2] = C.float(v.Z)
		}
	}
	
	C.tessAddContour(t.tess, C.int(size), unsafe.Pointer(&cVerts[0]), C.int(size*4), C.int(len(vertices)))
	return nil
}

// SetOption enables or disables a tessellation option.
func (t *Tessellator) SetOption(option Option, enabled bool) error {
	if t == nil || t.tess == nil {
		return fmt.Errorf("tessellator is nil or deleted")
	}
	
	value := 0
	if enabled {
		value = 1
	}
	
	C.tessSetOption(t.tess, C.int(option), C.int(value))
	return nil
}

// Tessellate performs the tessellation operation.
// windingRule: the winding rule to use
// elementType: the type of output elements
// polySize: maximum vertices per polygon (for polygon output)
// vertexSize: number of coordinates per vertex (2 or 3)
// normal: normal vector (can be nil for auto-calculation)
func (t *Tessellator) Tessellate(windingRule WindingRule, elementType ElementType, polySize, vertexSize int, normal *Vertex) error {
	if t == nil || t.tess == nil {
		return fmt.Errorf("tessellator is nil or deleted")
	}
	
	if vertexSize != 2 && vertexSize != 3 {
		return fmt.Errorf("vertexSize must be 2 or 3, got %d", vertexSize)
	}
	
	var normalPtr *C.TESSreal
	if normal != nil {
		normalArray := [3]C.TESSreal{C.TESSreal(normal.X), C.TESSreal(normal.Y), C.TESSreal(normal.Z)}
		normalPtr = &normalArray[0]
	}
	
	result := C.tessTesselate(t.tess, C.int(windingRule), C.int(elementType), C.int(polySize), C.int(vertexSize), normalPtr)
	
	if result == 0 {
		status := t.GetStatus()
		return fmt.Errorf("tessellation failed with status: %v", status)
	}
	
	return nil
}

// GetVertexCount returns the number of vertices in the tessellated output.
func (t *Tessellator) GetVertexCount() int {
	if t == nil || t.tess == nil {
		return 0
	}
	return int(C.tessGetVertexCount(t.tess))
}

// GetVertices returns the tessellated vertices.
func (t *Tessellator) GetVertices() []Vertex {
	if t == nil || t.tess == nil {
		return nil
	}
	
	count := t.GetVertexCount()
	if count == 0 {
		return nil
	}
	
	// Get pointer to vertex data
	vertexPtr := C.tessGetVertices(t.tess)
	if vertexPtr == nil {
		return nil
	}
	
	// Convert C array to Go slice using unsafe.Slice
	ptr := unsafe.Pointer(vertexPtr)
	verts := unsafe.Slice((*C.TESSreal)(ptr), count*3) // 3 coordinates per vertex
	
	// Convert to Vertex slice
	vertices := make([]Vertex, count)
	for i := 0; i < count; i++ {
		offset := i * 3
		vertices[i] = Vertex{
			X: float32(verts[offset]),
			Y: float32(verts[offset+1]),
			Z: float32(verts[offset+2]),
		}
	}
	
	return vertices
}

// GetVertexIndices returns the vertex indices mapping generated vertices to original vertices.
func (t *Tessellator) GetVertexIndices() []int {
	if t == nil || t.tess == nil {
		return nil
	}
	
	count := t.GetVertexCount()
	if count == 0 {
		return nil
	}
	
	// Get pointer to index data
	indexPtr := C.tessGetVertexIndices(t.tess)
	if indexPtr == nil {
		return nil
	}
	
	// Convert C array to Go slice using unsafe.Slice
	ptr := unsafe.Pointer(indexPtr)
	indices := unsafe.Slice((*C.TESSindex)(ptr), count)
	
	// Convert to int slice
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = int(indices[i])
	}
	
	return result
}

// GetElementCount returns the number of elements in the tessellated output.
func (t *Tessellator) GetElementCount() int {
	if t == nil || t.tess == nil {
		return 0
	}
	return int(C.tessGetElementCount(t.tess))
}

// GetElements returns the tessellated elements.
func (t *Tessellator) GetElements() []int {
	if t == nil || t.tess == nil {
		return nil
	}
	
	count := t.GetElementCount()
	if count == 0 {
		return nil
	}
	
	// Get pointer to element data
	elementPtr := C.tessGetElements(t.tess)
	if elementPtr == nil {
		return nil
	}
	
	// Convert C array to Go slice using unsafe.Slice
	// We'll use a reasonable default size based on typical usage
	// The actual size depends on element type and polySize, but this should work for most cases
	ptr := unsafe.Pointer(elementPtr)
	elements := unsafe.Slice((*C.TESSindex)(ptr), count*3) // Assume 3 vertices per element (triangles)
	
	// Convert to int slice
	result := make([]int, count*3)
	for i := 0; i < count*3; i++ {
		result[i] = int(elements[i])
	}
	
	return result
}

// GetElementsWithSize returns the tessellated elements with the specified element type and poly size.
// This is more accurate than GetElements() as it uses the correct array size.
func (t *Tessellator) GetElementsWithSize(elementType ElementType, polySize int) []int {
	if t == nil || t.tess == nil {
		return nil
	}
	
	count := t.GetElementCount()
	if count == 0 {
		return nil
	}
	
	// Get pointer to element data
	elementPtr := C.tessGetElements(t.tess)
	if elementPtr == nil {
		return nil
	}
	
	// Calculate the correct array size based on element type
	var arraySize int
	switch elementType {
	case ElementPolygons:
		arraySize = count * polySize
	case ElementConnectedPolygons:
		arraySize = count * polySize * 2
	case ElementBoundaryContours:
		arraySize = count * 2
	default:
		arraySize = count * polySize // Default to polygons
	}
	
	// Convert C array to Go slice using unsafe.Slice
	ptr := unsafe.Pointer(elementPtr)
	elements := unsafe.Slice((*C.TESSindex)(ptr), arraySize)
	
	// Convert to int slice
	result := make([]int, arraySize)
	for i := 0; i < arraySize; i++ {
		result[i] = int(elements[i])
	}
	
	return result
}

// GetTriangles returns the tessellated triangles as a slice of triangle indices.
// Each triangle is represented by 3 vertex indices.
func (t *Tessellator) GetTriangles() [][]int {
	elements := t.GetElementsWithSize(ElementPolygons, 3)
	if elements == nil {
		return nil
	}
	
	elementCount := t.GetElementCount()
	triangles := make([][]int, elementCount)
	
	for i := 0; i < elementCount; i++ {
		base := i * 3
		triangles[i] = []int{elements[base], elements[base+1], elements[base+2]}
	}
	
	return triangles
}

// GetStatus returns the current tessellation status.
func (t *Tessellator) GetStatus() Status {
	if t == nil || t.tess == nil {
		return StatusInvalidInput
	}
	return Status(C.tessGetStatus(t.tess))
}

// String returns a string representation of the winding rule.
func (w WindingRule) String() string {
	switch w {
	case WindingOdd:
		return "Odd"
	case WindingNonZero:
		return "NonZero"
	case WindingPositive:
		return "Positive"
	case WindingNegative:
		return "Negative"
	case WindingAbsGeqTwo:
		return "AbsGeqTwo"
	default:
		return "Unknown"
	}
}

// String returns a string representation of the element type.
func (e ElementType) String() string {
	switch e {
	case ElementPolygons:
		return "Polygons"
	case ElementConnectedPolygons:
		return "ConnectedPolygons"
	case ElementBoundaryContours:
		return "BoundaryContours"
	default:
		return "Unknown"
	}
}

// String returns a string representation of the status.
func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusOutOfMemory:
		return "OutOfMemory"
	case StatusInvalidInput:
		return "InvalidInput"
	default:
		return "Unknown"
	}
} 