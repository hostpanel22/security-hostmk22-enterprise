package main

import (
	"fmt"
	"time"
)

const (
	Version = "0.1.0"
	Name    = "HostMK Security Framework"
)

func main() {
	fmt.Printf(
		"%s v%s started at %s\n",
		Name,
		Version,
		time.Now().UTC().Format(time.RFC3339),
	)
}
