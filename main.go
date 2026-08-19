package main

import (
	"fmt"

	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/modules"
)

func main() {
	cfg := modules.LFMConfig{
		SampleRate: 10e6,  // 10 MHz
		PulseWidth: 10e-6, // 10 µs
		Bandwidth:  5e6,   // 5 MHz
	}

	signal := modules.GenerateLFM(cfg)

	fmt.Println("Number of samples:", len(signal.Samples))
	fmt.Println("Sample rate:", signal.SampleRate)

}
