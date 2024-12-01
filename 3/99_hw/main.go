package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	out := os.Stdout

	start := time.Now()
	SlowSearch(out)
	durationSlow := time.Since(start)

	start = time.Now()
	FastSearch(out)
	durationFast := time.Since(start)

	fmt.Println("Slow:", durationSlow)
	fmt.Println("Fast:", durationFast)
}
