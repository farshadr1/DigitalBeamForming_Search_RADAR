package modules

import (
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
)

// MatchedFilter with iFFT implementation
type MatchedFilter struct {
	fftLen int
	fft    *fourier.CmplxFFT
	H      []complex128 // precomputed FFT of the flipped-conjugated reference pulse
	outLen int
	hLen   int
}

func nextPow2(n int) int {
	p := 1
	for p < n {
		p *= 2
	}
	return p
}

// Constructure of MatchedFilter
func NewMatchedFilter(signalLen int, refPulse Signal) *MatchedFilter {
	hLen := len(refPulse.Samples)
	outLen := signalLen + hLen - 1
	fftLen := nextPow2(outLen)

	// h[n] = conj(x[M-1-n])  (flip + conjugate the reference pulse)
	h := make([]complex128, hLen)
	for i, v := range refPulse.Samples {
		h[hLen-1-i] = cmplx.Conj(v)
	}

	hPadded := make([]complex128, fftLen)
	copy(hPadded, h)

	fft := fourier.NewCmplxFFT(fftLen)
	H := fft.Coefficients(nil, hPadded)

	return &MatchedFilter{
		fftLen: fftLen,
		fft:    fft,
		H:      H,
		outLen: outLen,
		hLen:   hLen,
	}
}

// Apply Match Filter on signal
func (mf *MatchedFilter) Apply(signal Signal) Signal {
	sigPadded := make([]complex128, mf.fftLen)
	copy(sigPadded, signal.Samples)

	X := mf.fft.Coefficients(nil, sigPadded)

	Y := make([]complex128, mf.fftLen)
	for i := range Y {
		Y[i] = X[i] * mf.H[i]
	}

	// inverse FFT
	y := mf.fft.Sequence(nil, Y)

	// Normalize
	norm := complex(1.0/float64(mf.fftLen), 0)
	for i := range y {
		y[i] *= norm
	}

	return Signal{
		// trim output
		Samples:    y[:mf.outLen],
		SampleRate: signal.SampleRate,
	}
}
