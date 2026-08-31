package antenna

import (
	"fmt"
	"math"
	"math/cmplx"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

	// Using Go Routine
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

// -----------------------------------------------------------------------------
// Plot Elevation-cut
// -----------------------------------------------------------------------------

func PlotElevationBeamPatterns(
	pattern Pattern,
	beamAngles []float64,
	beamIndices []int,
	filename string,
) error {

	p := plot.New()

	p.Title.Text = "DBF Elevation Beam Patterns at Azimuth = 0°"
	p.X.Label.Text = "Elevation (deg)"
	p.Y.Label.Text = "Gain (dB)"

	p.Y.Tick.Marker = plot.ConstantTicks(
		[]plot.Tick{
			{Value: -40, Label: "-40"},
			{Value: -35, Label: "-35"},
			{Value: -30, Label: "-30"},
			{Value: -25, Label: "-25"},
			{Value: -20, Label: "-20"},
			{Value: -15, Label: "-15"},
			{Value: -10, Label: "-10"},
			{Value: -5, Label: "-5"},
			{Value: 0, Label: "0"},
		},
	)

	grid := plotter.NewGrid()
	p.Add(grid)

	// Place the legend in the upper-right corner.
	p.Legend.Top = true
	p.Legend.Left = false

	// Only one azimuth point: Az = 0°
	azIndex := 0

	for _, beam := range beamIndices {

		// Safety check
		if beam < 0 || beam >= len(pattern.Gain) {
			continue
		}

		pts := make(
			plotter.XYs,
			len(pattern.Elevations),
		)

		for ei, elevation := range pattern.Elevations {

			pts[ei].X = elevation
			pts[ei].Y = pattern.Gain[beam][ei][azIndex]
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			return err
		}

		line.Width = vg.Points(1.2)

		p.Add(line)

		p.Legend.Add(
			fmt.Sprintf(
				"Beam %d (%.2f°)",
				beam,
				beamAngles[beam],
			),
			line,
		)
	}

	p.X.Min = pattern.Elevations[0]
	p.X.Max = pattern.Elevations[len(pattern.Elevations)-1]

	p.Y.Min = -40
	p.Y.Max = 5

	return p.Save(
		10*vg.Inch,
		6*vg.Inch,
		filename,
	)
}

//-------------------------------------------------------------------------
// Plot Conture chart
//-------------------------------------------------------------------------

type BeamGrid struct {
	pattern Pattern
	beam    int
}

func (g BeamGrid) Dims() (c, r int) {
	return len(g.pattern.Azimuths),
		len(g.pattern.Elevations)
}

func (g BeamGrid) X(c int) float64 {
	return g.pattern.Azimuths[c]
}

func (g BeamGrid) Y(r int) float64 {
	return g.pattern.Elevations[r]
}

func (g BeamGrid) Z(c, r int) float64 {
	return g.pattern.Gain[g.beam][r][c]
}

func PlotBeamContour(
	pattern Pattern,
	beam int,
	filename string,
) error {

	grid := BeamGrid{
		pattern: pattern,
		beam:    beam,
	}

	p := plot.New()

	p.Title.Text = fmt.Sprintf(
		"Beam %d Antenna Pattern",
		beam,
	)

	p.X.Label.Text = "Azimuth (deg)"
	p.Y.Label.Text = "Elevation (deg)"

	// Contour levels in dB.
	levels := []float64{
		-40,
		-35,
		-30,
		-25,
		-20,
		-15,
		-10,
		-5,
		0,
	}

	// Create contour plot.
	contour := plotter.NewContour(
		grid,
		levels,
		palette.Rainbow(
			len(levels),
			palette.Blue,
			palette.Red,
			1,
			1,
			1,
		),
	)

	p.Add(contour)

	// Add normal X/Y grid.
	p.Add(plotter.NewGrid())

	// Axis limits.
	p.X.Min = pattern.Azimuths[0]
	p.X.Max = pattern.Azimuths[len(pattern.Azimuths)-1]

	p.Y.Min = pattern.Elevations[0]
	p.Y.Max = pattern.Elevations[len(pattern.Elevations)-1]

	return p.Save(
		10*vg.Inch,
		8*vg.Inch,
		filename,
	)
}
