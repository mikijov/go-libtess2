# Go Bindings for libtess2

[![Go Report Card](https://goreportcard.com/badge/github.com/mikijov/go-libtess2)](https://goreportcard.com/report/github.com/mikijov/go-libtess2)
[![GoDoc](https://godoc.org/github.com/mikijov/go-libtess2?status.svg)](https://godoc.org/github.com/mikijov/go-libtess2)
[![License](https://img.shields.io/badge/License-SGI%20Free%20Software%20License%20B%20v2.0-blue.svg)](https://directory.fsf.org/wiki/License:SGIFreeBv2)

Go bindings for the [libtess2](https://github.com/memononen/libtess2) polygon tessellation library. libtess2 is a high-quality polygon tessellator and triangulator library that can handle complex polygons with holes, self-intersections, and other challenging geometric cases.

## Features

- **Complete API Coverage**: Full Go bindings for all libtess2 functionality
- **Idiomatic Go Interface**: Clean, Go-style API with proper error handling
- **Memory Management**: Automatic cleanup with Go's garbage collector
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Comprehensive Testing**: Full test coverage with benchmarks
- **Documentation**: Complete API documentation and examples

## Installation

### Prerequisites

- Go 1.18 or later
- GCC compiler (for building the C library)
- Make (for build automation)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/mikijov/go-libtess2.git
cd go-libtess2

# Build the library and run tests
make

# Install the package
make install
```

### Manual Installation

```bash
# Build the C library
make libtess2/libtess2.a

# Run tests
go test -v ./...

# Install
go install ./...
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/mikijov/go-libtess2"
)

func main() {
    // Create a new tessellator
    tess := tess.NewTessellator()
    if tess == nil {
        log.Fatal("Failed to create tessellator")
    }
    defer tess.Delete()

    // Define a simple triangle
    vertices := []tess.Vertex{
        {X: 0, Y: 0, Z: 0},   // Bottom-left
        {X: 1, Y: 0, Z: 0},   // Bottom-right
        {X: 0.5, Y: 1, Z: 0}, // Top
    }

    // Add the contour
    err := tess.AddContour(2, vertices)
    if err != nil {
        log.Fatalf("Failed to add contour: %v", err)
    }

    // Perform tessellation
    err = tess.Tessellate(tess.WindingOdd, tess.ElementPolygons, 3, 2, nil)
    if err != nil {
        log.Fatalf("Tessellation failed: %v", err)
    }

    // Get results
    outputVertices := tess.GetVertices()
    elements := tess.GetElements()

    fmt.Printf("Vertices: %v\n", outputVertices)
    fmt.Printf("Triangles: %v\n", elements)
}
```

### Complex Polygon with Holes

```go
// Define a square with a triangular hole
outerContour := []tess.Vertex{
    {X: 0, Y: 0, Z: 0},   // Bottom-left
    {X: 4, Y: 0, Z: 0},   // Bottom-right
    {X: 4, Y: 4, Z: 0},   // Top-right
    {X: 0, Y: 4, Z: 0},   // Top-left
}

innerContour := []tess.Vertex{
    {X: 1, Y: 1, Z: 0},   // Bottom-left of hole
    {X: 3, Y: 1, Z: 0},   // Bottom-right of hole
    {X: 2, Y: 3, Z: 0},   // Top of hole
}

// Add contours
tess.AddContour(2, outerContour)
tess.AddContour(2, innerContour)

// Enable constrained Delaunay triangulation
tess.SetOption(tess.OptionConstrainedDelaunay, true)

// Tessellate
tess.Tessellate(tess.WindingOdd, tess.ElementPolygons, 3, 2, nil)
```

## API Reference

### Types

#### Vertex

```go
type Vertex struct {
    X, Y, Z float32
}
```

Represents a 2D or 3D vertex.

#### Tessellator

```go
type Tessellator struct {
    // Private fields
}
```

Main tessellation context.

### Winding Rules

- `WindingOdd`: Standard odd-even rule
- `WindingNonZero`: Non-zero winding rule
- `WindingPositive`: Positive winding rule
- `WindingNegative`: Negative winding rule
- `WindingAbsGeqTwo`: Absolute value >= 2 rule

### Element Types

- `ElementPolygons`: Output as individual polygons
- `ElementConnectedPolygons`: Output as connected polygons with neighbor information
- `ElementBoundaryContours`: Output as boundary contours

### Options

- `OptionConstrainedDelaunay`: Enable constrained Delaunay triangulation
- `OptionReverseContours`: Reverse contour winding

### Methods

#### NewTessellator()

Creates a new tessellator instance.

#### Delete()

Destroys the tessellator and frees memory.

#### AddContour(size int, vertices []Vertex) error

Adds a contour to be tessellated. `size` must be 2 or 3.

#### SetOption(option Option, enabled bool) error

Enables or disables tessellation options.

#### Tessellate(windingRule WindingRule, elementType ElementType, polySize, vertexSize int, normal \*Vertex) error

Performs the tessellation operation.

#### GetVertexCount() int

Returns the number of vertices in the output.

#### GetVertices() []Vertex

Returns the tessellated vertices.

#### GetElementCount() int

Returns the number of elements in the output.

#### GetElements() []int

Returns the tessellated elements (triangle indices).

#### GetVertexIndices() []int

Returns vertex indices mapping to original vertices.

#### GetStatus() Status

Returns the current tessellation status.

## Examples

The repository includes several example programs:

### Simple Triangle

```bash
cd examples/simple_triangle
go run main.go
```

### Complex Polygon

```bash
cd examples/complex_polygon
go run main.go
```

## Building

### Build Targets

- `make` - Build library and run tests
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make test-race` - Run tests with race detection
- `make examples` - Build example programs
- `make install` - Install the package
- `make bench` - Run benchmarks
- `make fmt` - Format code
- `make lint` - Run linter

### Cross-Platform Building

The library supports cross-platform compilation. Set the appropriate environment variables:

```bash
# For Windows
GOOS=windows GOARCH=amd64 go build

# For macOS
GOOS=darwin GOARCH=amd64 go build

# For Linux
GOOS=linux GOARCH=amd64 go build
```

## Testing

Run the test suite:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem ./...
```

Run tests with race detection:

```bash
go test -race -v ./...
```

## Performance

The Go bindings provide near-native performance with minimal overhead. Benchmarks show:

- Simple polygon tessellation: ~2.8μs per operation
- Memory overhead: ~56 bytes per operation
- Allocation count: 2 allocations per operation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the SGI Free Software License B Version 2.0, the same license as libtess2.

## Acknowledgments

- [Mikko Mononen](https://github.com/memononen) for the original libtess2 library
- The OpenGL community for the original tessellation algorithms

## Links

- [libtess2 Repository](https://github.com/memononen/libtess2)
- [Go Documentation](https://godoc.org/github.com/mikijov/go-libtess2)
- [OpenGL Programming Guide (Red Book)](https://www.glprogramming.com/red/chapter11.html)
