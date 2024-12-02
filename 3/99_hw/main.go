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

	start = time.Now()
	FastSearchV0_1(out)
	durationFast01 := time.Since(start)

	fmt.Println("Slow:", durationSlow)
	fmt.Println("Fast:", durationFast)
	fmt.Println("Fast01:", durationFast01)
}
