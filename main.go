package main

import (
	"fmt"
	"log"

	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/antenna/antenna"
	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/modules"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// tx_signal := modules.GenerateLFM(modules.LFMConfig{
	// 	SampleRate: cfg.Signal.SampleRate,
	// 	PulseWidth: cfg.Signal.PulseWidth,
	// 	Bandwidth:  cfg.Signal.Bandwidth,
	// })

	// pri := float64(200e-6)
	// echoed_signal := modules.GenerateEcho(tx_signal, tgt1, cfg.Radar.CarrierFreq, pri)

	// mf := modules.NewMatchedFilter(len(echoed_signal.Samples), tx_signal)
	// mf_signal := mf.Apply(echoed_signal)
	// figures.Output_inTime(mf_signal)

	tgts := make([]modules.Target, len(cfg.Targets))
	for i, target := range cfg.Targets {
		fmt.Printf("Target %d\n", i+1)
		fmt.Printf("  Position: %v\n", target.Position)
		fmt.Printf("  Velocity: %v\n", target.Velocity)
		fmt.Printf("  RCS:      %.2f m²\n", target.RCS)

		tgts[i].ID = i
		tgts[i].X, tgts[i].Y, tgts[i].Z = target.Position[0], target.Position[1], target.Position[2]
		tgts[i].VelocityX, tgts[i].VelocityY, tgts[i].VelocityZ = target.Velocity[0], target.Velocity[1], target.Velocity[2]
		tgts[i].RCS = target.RCS
	}

	beamAngles := antenna.GenerateBeamAngles(
		32,
		0,
		60.0,
	)

	for i, angle := range beamAngles {
		fmt.Printf(
			"Beam %02d : %+6.2f deg\n",
			i,
			angle,
		)
	}

}
