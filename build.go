//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the directory where this script is located
	scriptDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Check if libtess2.a exists
	libPath := filepath.Join(scriptDir, "libtess2", "libtess2.a")
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		fmt.Println("libtess2.a not found. Building libtess2...")
		
		// Run make to build the library
		cmd := exec.Command("make", "libtess2/libtess2.a")
		cmd.Dir = scriptDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error building libtess2: %v\n", err)
			fmt.Fprintf(os.Stderr, "Please ensure you have gcc and make installed.\n")
			os.Exit(1)
		}
		
		fmt.Println("libtess2.a built successfully.")
	} else {
		fmt.Println("libtess2.a already exists.")
	}
} 