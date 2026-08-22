package main

import (
	"log"

	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/figures"
	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/modules"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	tx_signal := modules.GenerateLFM(modules.LFMConfig{
		SampleRate: cfg.Signal.SampleRate,
		PulseWidth: cfg.Signal.PulseWidth,
		Bandwidth:  cfg.Signal.Bandwidth,
	})

	tgt1 := modules.Target{
		ID:        1,
		X:         5e3,
		Y:         12e3,
		Z:         900,
		VelocityX: 100,
		VelocityY: 100,
		VelocityZ: 0,
		RCS:       10,
	}

	pri := float64(200e-6)
	echoed_signal := modules.GenerateEcho(tx_signal, tgt1, cfg.Radar.CarrierFreq, pri)

	mf := modules.NewMatchedFilter(len(echoed_signal.Samples), tx_signal)
	mf_signal := mf.Apply(echoed_signal)
	figures.Output_inSample(mf_signal)
}
