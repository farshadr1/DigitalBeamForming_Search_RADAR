package modules

import (
	"math"
	"math/cmplx"
)

// Signal represents a sampled complex baseband signal.
type Signal struct {
	Samples    []complex128
	SampleRate float64 // Hz
}

// LFMConfig contains parameters for an LFM chirp.
type LFMConfig struct {
	SampleRate float64 // Sampling frequency (Hz)
	PulseWidth float64 // Pulse duration (seconds)
	Bandwidth  float64 // Chirp bandwidth (Hz)
}

// GenerateLFM generates a complex baseband Linear Frequency Modulated (LFM) pulse.
//
// The generated signal is:
//
//	s(t) = exp(j * pi * k * t^2)
//
// where:
//
//	k = Bandwidth / PulseWidthcomplex128
//
// The instantaneous frequency sweeps approximately from
// -Bandwidth/2 to +Bandwidth/2.
func GenerateLFM(cfg LFMConfig) Signal {
	numSamples := int(math.Round(cfg.SampleRate * cfg.PulseWidth))

	samples := make([]complex128, numSamples)

	chirpRate := cfg.Bandwidth / cfg.PulseWidth

	for n := 0; n < numSamples; n++ {
		t := float64(n) / cfg.SampleRate

		// Center time around zero.
		t -= cfg.PulseWidth / 2

		phase := math.Pi * chirpRate * t * t

		samples[n] = cmplx.Exp(complex(0, phase))
	}

	return Signal{
		Samples:    samples,
		SampleRate: cfg.SampleRate,
	}
}
