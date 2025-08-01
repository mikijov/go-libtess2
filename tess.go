// Package tess provides Go bindings for the libtess2 polygon tessellation library.
// libtess2 is a high-quality polygon tessellator and triangulator library.
//
// This package provides an idiomatic Go interface to the C library, allowing
// you to tessellate complex polygons into triangles, contours, and other
// geometric primitives.
package tess

/*
#cgo CFLAGS: -I${SRCDIR}/libtess2/Include -I${SRCDIR}/libtess2/Source
#cgo LDFLAGS: -L${SRCDIR}/libtess2 -l:libtess2.a -lm
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
	WindingOdd       WindingRule = C.TESS_WINDING_ODD
	WindingNonZero   WindingRule = C.TESS_WINDING_NONZERO
	WindingPositive  WindingRule = C.TESS_WINDING_POSITIVE
	WindingNegative  WindingRule = C.TESS_WINDING_NEGATIVE
	WindingAbsGeqTwo WindingRule = C.TESS_WINDING_ABS_GEQ_TWO
)

// ElementType defines the type of output elements from tessellation.
type ElementType int

const (
	ElementPolygons          ElementType = C.TESS_POLYGONS
	ElementConnectedPolygons ElementType = C.TESS_CONNECTED_POLYGONS
	ElementBoundaryContours  ElementType = C.TESS_BOUNDARY_CONTOURS
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
	StatusOK           Status = C.TESS_STATUS_OK
	StatusOutOfMemory  Status = C.TESS_STATUS_OUT_OF_MEMORY
	StatusInvalidInput Status = C.TESS_STATUS_INVALID_INPUT
)

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
func (t *Tessellator) AddContour(size int, vertices []float32) error {
	if t == nil || t.tess == nil {
		return fmt.Errorf("tessellator is nil or deleted")
	}

	if size != 2 && size != 3 {
		return fmt.Errorf("size must be 2 or 3, got %d", size)
	}
	if len(vertices) < size {
		return fmt.Errorf("vertices slice must contain at least one vertex")
	}
	if len(vertices)%size != 0 {
		return fmt.Errorf("len(vertices)(%d) must be multiple of size (%d)", len(vertices), size)
	}

	// stride := uintptr(unsafe.Pointer(&vertices[size])) - uintptr(unsafe.Pointer(&vertices[0]))
	stride := 4 * size
	// fmt.Printf("size:%d len:%d stride:%d\n", size, len(vertices)/2, stride)

	C.tessAddContour(
		t.tess,
		C.int(size),
		unsafe.Pointer(&vertices[0]),
		C.int(stride),
		C.int(len(vertices)/size),
	)

	status := t.getStatus()
	if status != StatusOK {
		return fmt.Errorf("error adding contour: %s", status)
	}
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

// Tessellate performs the tessellation operation and returns vertices and indices in a single call.
// windingRule: the winding rule to use
// elementType: the type of output elements
// polySize: maximum vertices per polygon (for polygon output)
// vertexSize: number of coordinates per vertex (2 or 3)
// normal: normal vector as []float32 (can be nil for auto-calculation)
// Returns:
//   - vertices: flat slice of vertex coordinates (size * vertexCount)
//   - indices: slice of vertex indices for elements
//   - err: error if tessellation fails
func (t *Tessellator) Tessellate(windingRule WindingRule, elementType ElementType, polySize, vertexSize int, normal []float32) (vertices []float32, indices []int, err error) {
	if t == nil || t.tess == nil {
		return nil, nil, fmt.Errorf("tessellator is nil or deleted")
	}

	// Perform tessellation
	err = t.internalTessellate(windingRule, elementType, polySize, vertexSize, normal)
	if err != nil {
		return nil, nil, err
	}

	// indexPtr := C.tessGetVertexIndices(t.tess)
	// indexCount := C.tessGetElementCount(t.tess)
	// if indexPtr != nil {
	// 	// Convert C array to Go slice using unsafe.Slice
	// 	ptr := unsafe.Pointer(indexPtr)
	// 	indices := unsafe.Slice((*C.TESSindex)(ptr), indexCount)
	// 	fmt.Printf("indices %d\n", indexCount)
	// 	for i, idx := range indices {
	// 		fmt.Printf("%2d %8x\n", i, idx)
	// 	}
	// }

	// Get vertices
	vertices = t.getVertices(vertexSize)
	if vertices == nil {
		status := t.getStatus()
		if status != StatusOK {
			return nil, nil, fmt.Errorf("failed to get vertices: %s", status)
		}
		vertices = []float32{}
	}

	// Get indices based on element type
	switch elementType {
	case ElementPolygons:
		indices = t.getElementsWithSize(elementType, polySize)
	case ElementConnectedPolygons:
		indices = t.getElementsWithSize(elementType, polySize)
	case ElementBoundaryContours:
		indices = t.getElementsWithSize(elementType, polySize)
	default:
		return nil, nil, fmt.Errorf("unsupported element type: %v", elementType)
	}

	if indices == nil {
		status := t.getStatus()
		if status != StatusOK {
			return nil, nil, fmt.Errorf("failed to get vertices: %s", status)
		}
		indices = []int{}
	}

	return vertices, indices, nil
}

// internalTessellate performs the tessellation operation.
// windingRule: the winding rule to use
// elementType: the type of output elements
// polySize: maximum vertices per polygon (for polygon output)
// vertexSize: number of coordinates per vertex (2 or 3)
// normal: normal vector as []float32 (can be nil for auto-calculation)
func (t *Tessellator) internalTessellate(windingRule WindingRule, elementType ElementType, polySize, vertexSize int, normal []float32) error {
	if t == nil || t.tess == nil {
		return fmt.Errorf("tessellator is nil or deleted")
	}

	if vertexSize != 2 && vertexSize != 3 {
		return fmt.Errorf("vertexSize must be 2 or 3, got %d", vertexSize)
	}

	var normalPtr *C.TESSreal
	if normal != nil {
		if len(normal) < 3 {
			return fmt.Errorf("normal vector must have at least 3 components, got %d", len(normal))
		}
		normalArray := [3]C.TESSreal{C.TESSreal(normal[0]), C.TESSreal(normal[1]), C.TESSreal(normal[2])}
		normalPtr = &normalArray[0]
	}

	result := C.tessTesselate(t.tess, C.int(windingRule), C.int(elementType), C.int(polySize), C.int(vertexSize), normalPtr)

	if result == 0 {
		status := t.getStatus()
		return fmt.Errorf("tessellation failed with status: %v", status)
	}

	return nil
}

// getVertexCount returns the number of vertices in the tessellated output.
func (t *Tessellator) getVertexCount() int {
	if t == nil || t.tess == nil {
		return 0
	}
	return int(C.tessGetVertexCount(t.tess))
}

// getVertices returns the tessellated vertices as a flat slice of float32 coordinates.
// The size parameter indicates the number of coordinates per vertex (2 or 3).
// Returns nil if no vertices are available.
func (t *Tessellator) getVertices(size int) []float32 {
	if t == nil || t.tess == nil {
		return nil
	}

	count := t.getVertexCount()
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
	verts := unsafe.Slice((*C.TESSreal)(ptr), count*size)

	// Convert to float32 slice
	vertices := make([]float32, count*size)
	for i := 0; i < count*size; i++ {
		vertices[i] = float32(verts[i])
	}

	return vertices
}

// getVertexIndices returns the vertex indices mapping generated vertices to original vertices.
func (t *Tessellator) getVertexIndices() []int {
	if t == nil || t.tess == nil {
		return nil
	}

	count := t.getVertexCount()
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

// getElementCount returns the number of elements in the tessellated output.
func (t *Tessellator) getElementCount() int {
	if t == nil || t.tess == nil {
		return 0
	}
	return int(C.tessGetElementCount(t.tess))
}

// getElements returns the tessellated elements.
func (t *Tessellator) getElements() []int {
	if t == nil || t.tess == nil {
		return nil
	}

	count := t.getElementCount()
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

// getElementsWithSize returns the tessellated elements with the specified element type and poly size.
// This is more accurate than getElements() as it uses the correct array size.
func (t *Tessellator) getElementsWithSize(elementType ElementType, polySize int) []int {
	if t == nil || t.tess == nil {
		return nil
	}

	count := t.getElementCount()
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

// getTriangles returns the tessellated triangles as a slice of triangle indices.
// Each triangle is represented by 3 vertex indices.
func (t *Tessellator) getTriangles() [][]int {
	elements := t.getElementsWithSize(ElementPolygons, 3)
	if elements == nil {
		return nil
	}

	elementCount := t.getElementCount()
	triangles := make([][]int, elementCount)

	for i := 0; i < elementCount; i++ {
		base := i * 3
		triangles[i] = []int{elements[base], elements[base+1], elements[base+2]}
	}

	return triangles
}

// getStatus returns the current tessellation status.
func (t *Tessellator) getStatus() Status {
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
