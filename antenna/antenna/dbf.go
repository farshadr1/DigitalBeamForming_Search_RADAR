package antenna

import (
	"math"
	"math/cmplx"
	"sync"
)

// -----------------------------------------------------------------------------
// Configuration
// -----------------------------------------------------------------------------

type ArrayConfig struct {
	NumElements    int
	ElementSpacing float64
	Frequency      float64

	ElementAzHPBW float64
	ElementElHPBW float64

	NumBeams int
}

// -----------------------------------------------------------------------------
// Usefull functions
// -----------------------------------------------------------------------------

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func Wavelength(freq float64, speedOfLight float64) float64 {
	return speedOfLight / freq
}

func wrapAngle180(angle float64) float64 {
	if angle > 180 {
		angle -= 360
	}
	if angle < -180 {
		angle += 360
	}
	return angle
}

func powerToDB(x float64) float64 {
	if x <= 1e-12 {
		return -120
	}
	return 10 * math.Log10(x)
}

// -----------------------------------------------------------------------------
// Element pattern
// -----------------------------------------------------------------------------
//
// Returns normalized element POWER gain.
// We approximate the element pattern with a Gaussian-like pattern:
// G = exp(-4ln(2) * (angle / HPBW)^2)

func elementPattern(
	azimuthDeg float64,
	elevationDeg float64,
	cfg ArrayConfig,
) float64 {

	az := wrapAngle180(azimuthDeg)
	el := elevationDeg

	azTerm := math.Exp(
		-4.0 * math.Ln2 *
			math.Pow(az/cfg.ElementAzHPBW, 2),
	)

	elTerm := math.Exp(
		-4.0 * math.Ln2 *
			math.Pow(el/cfg.ElementElHPBW, 2),
	)

	return azTerm * elTerm
}

// -----------------------------------------------------------------------------
// Array geometry
// -----------------------------------------------------------------------------
//
// z_n = element position relative to array center.

func elementPosition(index int, cfg ArrayConfig) float64 {

	center := float64(cfg.NumElements-1) / 2.0

	return (float64(index) - center) * cfg.ElementSpacing
}

// -----------------------------------------------------------------------------
// Digital Beamforming weights
// -----------------------------------------------------------------------------
//
// Generate phase weights for a beam pointed at beamElevationDeg.
// wn = exp(-j * k * z * sin(phiBeam))
// These weights compensate the incoming wave phase and steer the receive beam.

func DBFWeights(
	beamElevationDeg float64,
	cfg ArrayConfig,
) []complex128 {

	lambda := Wavelength(cfg.Frequency, 3e8)
	k := 2.0 * math.Pi / lambda

	phi := deg2rad(beamElevationDeg)

	weights := make([]complex128, cfg.NumElements)

	for n := 0; n < cfg.NumElements; n++ {

		z := elementPosition(n, cfg)

		phase := -k * z * math.Sin(phi)

		weights[n] = cmplx.Exp(
			complex(0, phase),
		)
	}

	return weights
}

// -----------------------------------------------------------------------------
// Array Factor
// -----------------------------------------------------------------------------
//
// Computes normalized ARRAY POWER GAIN for a particular DBF beam.
// The incoming signal from elevation phi produces:
// exp(j*k*z*sin(phi))
// The DBF weights provide: w[n]
// Their combination produces the beam response.

func arrayFactor(
	azimuthDeg float64,
	elevationDeg float64,
	beamElevationDeg float64,
	cfg ArrayConfig,
	weights []complex128,
) float64 {

	_ = azimuthDeg
	lambda := Wavelength(cfg.Frequency, 3e8)
	k := 2.0 * math.Pi / lambda
	phi := deg2rad(elevationDeg)

	var sum complex128

	for n := 0; n < cfg.NumElements; n++ {

		z := elementPosition(n, cfg)

		// Incoming wave phase
		phase := k * z * math.Sin(phi)

		signalPhase := cmplx.Exp(
			complex(0, phase),
		)

		sum += weights[n] * signalPhase
	}

	// Normalize amplitude by number of elements.
	amplitude := cmplx.Abs(sum) / float64(cfg.NumElements)

	// Convert amplitude to power.
	power := amplitude * amplitude

	return power
}

// -----------------------------------------------------------------------------
// Total antenna gain
// -----------------------------------------------------------------------------
//
// gain = element pattern * array factor

func gain(
	azimuthDeg float64,
	elevationDeg float64,
	beamElevationDeg float64,
	cfg ArrayConfig,
	weights []complex128,
) float64 {

	elementGain := elementPattern(
		azimuthDeg,
		elevationDeg,
		cfg,
	)

	afGain := arrayFactor(
		azimuthDeg,
		elevationDeg,
		beamElevationDeg,
		cfg,
		weights,
	)

	return elementGain * afGain
}

// -----------------------------------------------------------------------------
// Gain in dB
// -----------------------------------------------------------------------------

func GainDB(
	azimuthDeg float64,
	elevationDeg float64,
	beamElevationDeg float64,
	cfg ArrayConfig,
	weights []complex128,
) float64 {

	power := gain(azimuthDeg,
		elevationDeg,
		beamElevationDeg,
		cfg,
		weights,
	)

	return powerToDB(power)
}

// -----------------------------------------------------------------------------
// Generate DBF beams
// -----------------------------------------------------------------------------
//
// Generates equally spaced elevation beams.

func GenerateBeamAngles(
	numBeams int,
	centerBeam float64,
	coverage float64,
) []float64 {

	if numBeams == 1 {
		return []float64{centerBeam}
	}

	beams := make([]float64, numBeams)

	var centerIndex int

	if numBeams%2 == 1 {
		//odd
		centerIndex = (numBeams - 1) / 2
	} else {
		centerIndex = numBeams / 2
	}

	step := coverage / float64(numBeams-1)

	for i := range beams {
		beams[i] = centerBeam +
			float64(i-centerIndex)*step
	}

	return beams
}

// -----------------------------------------------------------------------------
// Generate one 2-D antenna pattern
// -----------------------------------------------------------------------------

type Pattern struct {
	Azimuths   []float64
	Elevations []float64

	// Gain[beam][elevation][azimuth]
	Gain [][][]float64
}

func GeneratePattern(
	cfg ArrayConfig,
	beamAngles []float64,
	azimuths []float64,
	elevations []float64,
) Pattern {

	pattern := Pattern{
		Azimuths:   azimuths,
		Elevations: elevations,
		Gain:       make([][][]float64, len(beamAngles)),
	}

	// Precompute DBF weights.
	weights := make([][]complex128, len(beamAngles))

	for b, beamAngle := range beamAngles {
		weights[b] = DBFWeights(
			beamAngle,
			cfg,
		)
	}

	// Each beam can be calculated independently.
	var wg sync.WaitGroup

	for b := range beamAngles {

		wg.Add(1)

		go func(beamIndex int) {

			defer wg.Done()

			beamGain := make([][]float64, len(elevations))

			for ei, el := range elevations {

				beamGain[ei] = make(
					[]float64,
					len(azimuths),
				)

				for ai, az := range azimuths {

					beamGain[ei][ai] = GainDB(
						az,
						el,
						beamAngles[beamIndex],
						cfg,
						weights[beamIndex],
					)
				}
			}

			pattern.Gain[beamIndex] = beamGain

		}(b)
	}

	wg.Wait()

	return pattern
}
