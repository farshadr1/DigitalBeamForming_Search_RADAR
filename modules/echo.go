package modules

import (
	"math"
	"math/cmplx"
)

const (
	SpeedOfLight = 300000000.0 // m/s
)

// r(t) = amp * s(t - τ) * exp(j2πfdt)
func GenerateEcho(tx Signal,
	target Target,
	carrierFreq float64,
) Signal {

	r := target.Range()

	// Round-trip delay
	delaySec := 2.0 * r / SpeedOfLight

	// Delay in samples
	delaySamples := int(math.Round(delaySec * tx.SampleRate))

	// Doppler frequency
	lambda := SpeedOfLight / carrierFreq
	fd := 2.0 * target.RadialVelocity() / lambda

	// Very simple amplitude model
	amp := math.Sqrt(target.RCS)

	echo := make(
		[]complex128,
		len(tx.Samples)+delaySamples,
	)

	for n := range tx.Samples {

		t := float64(n) / tx.SampleRate

		doppler :=
			cmplx.Exp(
				complex(0, (2.0 * math.Pi * fd * t)),
			)

		echo[n+delaySamples] =
			complex(amp, 0) *
				tx.Samples[n] *
				doppler
	}

	return Signal{
		Samples:    echo,
		SampleRate: tx.SampleRate,
	}
}
