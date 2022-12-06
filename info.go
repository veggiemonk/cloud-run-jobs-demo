package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

func printBuildInfo() {
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Git commit: %s\n", gitSha)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		_, err := fmt.Fprintln(os.Stderr, "no build info available")
		if err != nil {
			return
		}
	}
	fmt.Printf("Build info: %s\n", info)
}
