package main

import (
	"fmt"
	"log"

	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/modules"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	signal := modules.GenerateLFM(modules.LFMConfig{
		SampleRate: cfg.Signal.SampleRate,
		PulseWidth: cfg.Signal.PulseWidth,
		Bandwidth:  cfg.Signal.Bandwidth,
	})

	fmt.Println("Sample rate:", signal.SampleRate)
	fmt.Println("Number of samples:", len(signal.Samples))
	fmt.Println("Radar carrierFreq:", cfg.Radar.CarrierFreq)

}
