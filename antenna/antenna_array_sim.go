package main

import (
	"fmt"
	"math"
	"modules/antenna"
)

const (
	NumElements = 48
	NumBeams    = 32

	// Physical array parameters
	ElementSpacing = 0.015 // 15 mm

	// Element pattern
	ElementAzHPBW = 1.5  // degrees
	ElementElHPBW = 50.0 // degrees

	// Example RF frequency.
	// Make this configurable for your real radar.
	Frequency = 10e9 // 10 GHz

	// Speed of light
	C = 299792458.0
)

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

func main() {

	cfg := antenna.ArrayConfig{
		NumElements:    NumElements,
		ElementSpacing: ElementSpacing,
		Frequency:      Frequency,

		ElementAzHPBW: ElementAzHPBW,
		ElementElHPBW: ElementElHPBW,

		NumBeams: NumBeams,
	}

	lambda := antenna.Wavelength(cfg.Frequency, 3e8)

	fmt.Println("Radar Column Array Simulation")
	fmt.Println("--------------------------------")

	fmt.Printf("Elements       : %d\n", cfg.NumElements)
	fmt.Printf("Spacing        : %.3f mm\n", cfg.ElementSpacing*1000)
	fmt.Printf("Frequency      : %.2f GHz\n", cfg.Frequency/1e9)
	fmt.Printf("Wavelength     : %.3f mm\n", lambda*1000)
	fmt.Printf("Spacing/lambda : %.3f\n",
		cfg.ElementSpacing/lambda)

	fmt.Printf("Element Az HPBW: %.2f deg\n",
		cfg.ElementAzHPBW)

	fmt.Printf("Element El HPBW: %.2f deg\n",
		cfg.ElementElHPBW)

	fmt.Printf("DBF Beams      : %d\n",
		cfg.NumBeams)

	// -------------------------------------------------------------------------
	// Generate 32 elevation beams.
	//
	// Here we cover -30° to +30°.

	beamAngles := antenna.GenerateBeamAngles(
		cfg.NumBeams,
		0,
		60.0,
	)

	fmt.Println("\nDBF beam pointing angles:")

	for i, angle := range beamAngles {
		fmt.Printf(
			"Beam %02d : %+6.2f deg\n",
			i,
			angle,
		)
	}

	// -------------------------------------------------------------------------
	// Example: calculate gain at one point for Beam 16.
	// -------------------------------------------------------------------------

	beamIndex := 16
	beamElevation := beamAngles[beamIndex]

	weights := antenna.DBFWeights(
		beamElevation,
		cfg,
	)

	az := 0.0
	el := beamElevation

	g := antenna.GainDB(
		az,
		el,
		beamElevation,
		cfg,
		weights,
	)

	fmt.Printf(
		"\nBeam %d gain at Az=%.1f°, El=%.1f°: %.2f dB\n",
		beamIndex,
		az,
		el,
		g,
	)

	// -------------------------------------------------------------------------
	// Generate complete pattern for beamIndex = 16
	//
	// 1° azimuth resolution in [-40 40]
	// 0.1° elevation resolution in [-10 10]
	// -------------------------------------------------------------------------

	var azimuths []float64

	for az := -10.0; az <= 10.0; az += 1.0 {
		azimuths = append(azimuths, az)
	}

	var elevations []float64

	for el := -40.0; el <= 40.0; el += 0.1 {
		elevations = append(elevations, el)
	}

	fmt.Println("\nGenerating 32-beam antenna pattern...")

	pattern := antenna.GeneratePattern(
		cfg,
		beamAngles,
		azimuths,
		elevations,
	)

	fmt.Println("Pattern generation complete.")

	// Example output.
	fmt.Printf(
		"Pattern dimensions: %d beams × %d elevations × %d azimuths\n",
		len(pattern.Gain),
		len(pattern.Elevations),
		len(pattern.Azimuths),
	)

	// Print gain around beam 16.
	fmt.Println("\nBeam 16 elevation cut:")

	for i, el := range elevations {

		if math.Abs(el-beamElevation) < 3.0 {

			// Azimuth = 0°
			azIndex := len(azimuths) / 2

			fmt.Printf(
				"El=%+6.2f°  Gain=%7.2f dB\n",
				el,
				pattern.Gain[beamIndex][i][azIndex],
			)
		}
	}

	// -------------------------------------------------------------------------
	// Plot Stacked Beams
	// -------------------------------------------------------------------------

}
