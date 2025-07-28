# Makefile for go-libtess2
# Builds libtess2 C library and provides Go build targets

# Compiler settings
CC = gcc
CFLAGS = -Wall -Wextra -O2 -fPIC -I./libtess2/Include
LDFLAGS = 

# Source files for libtess2
LIBTESS2_SOURCES = libtess2/Source/tess.c \
                   libtess2/Source/mesh.c \
                   libtess2/Source/sweep.c \
                   libtess2/Source/geom.c \
                   libtess2/Source/dict.c \
                   libtess2/Source/bucketalloc.c \
                   libtess2/Source/priorityq.c

# Object files
LIBTESS2_OBJECTS = $(LIBTESS2_SOURCES:.c=.o)

# Library name
LIBTESS2_LIB = libtess2/libtess2.a

# Default target
all: $(LIBTESS2_LIB) test

# Download and extract libtess2 source
download-libtess2:
	@echo "Downloading libtess2 source code..."
	curl -L "https://github.com/memononen/libtess2/archive/refs/heads/master.zip" -o libtess2.zip
	@echo "Extracting libtess2 source code..."
	unzip -o libtess2.zip
	rm -rf libtess2
	mv libtess2-master libtess2
	rm libtess2.zip
	@echo "libtess2 source code updated successfully"

# Update libtess2 to latest version
update-libtess2: download-libtess2
	@echo "libtess2 updated to latest version"

# Build libtess2 static library
$(LIBTESS2_LIB): $(LIBTESS2_OBJECTS)
	@echo "Building libtess2 static library..."
	ar rcs $@ $^
	@echo "Static library built: $@"

# Compile C source files
%.o: %.c
	@echo "Compiling $<..."
	$(CC) $(CFLAGS) -c $< -o $@

# Clean build artifacts (but keep the static library)
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(LIBTESS2_OBJECTS)
	rm -f *.o *.so *.dylib
	go clean -cache

# Clean everything including the static library
clean-all: clean
	@echo "Cleaning everything including static library..."
	rm -f $(LIBTESS2_LIB)

# Clean everything including downloaded source
clean-source: clean-all
	@echo "Cleaning everything including downloaded source..."
	rm -rf libtess2
	rm -f libtess2.zip

# Run Go tests
test: $(LIBTESS2_LIB)
	@echo "Running Go tests..."
	go test -v ./...

# Run Go tests with race detection
test-race: $(LIBTESS2_LIB)
	@echo "Running Go tests with race detection..."
	go test -race -v ./...

# Build Go examples
examples: $(LIBTESS2_LIB)
	@echo "Building examples..."
	go build -o examples/simple_triangle/simple_triangle examples/simple_triangle/main.go
	go build -o examples/complex_polygon/complex_polygon examples/complex_polygon/main.go

# Install the package
install: $(LIBTESS2_LIB)
	@echo "Installing package..."
	go install ./...

# Run benchmarks
bench: $(LIBTESS2_LIB)
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  all              - Build library and run tests"
	@echo "  download-libtess2 - Download latest libtess2 source"
	@echo "  update-libtess2  - Update libtess2 to latest version"
	@echo "  clean            - Clean build artifacts (keeps libtess2.a)"
	@echo "  clean-all        - Clean everything including libtess2.a"
	@echo "  clean-source     - Clean everything including source"
	@echo "  test             - Run Go tests"
	@echo "  test-race        - Run Go tests with race detection"
	@echo "  examples         - Build example programs"
	@echo "  install          - Install the package"
	@echo "  bench            - Run benchmarks"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run linter"
	@echo "  help             - Show this help"

.PHONY: all download-libtess2 update-libtess2 clean clean-all clean-source test test-race examples install bench fmt lint help 