package figures

import (
	"math/cmplx"

	"github.com/farshadr1/DigitalBeamForming_Search_RADAR/modules"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func Output_inSample(a modules.Signal) {
	if len(a.Samples) == 0 {
		return
	}

	points := make(plotter.XYs, len(a.Samples))
	minY := cmplx.Abs(a.Samples[0])
	maxY := cmplx.Abs(a.Samples[0])

	for i, sample := range a.Samples {
		points[i].X = float64(i)
		points[i].Y = cmplx.Abs(sample)

		if points[i].Y < minY {
			minY = points[i].Y
		}
		if points[i].Y > maxY {
			maxY = points[i].Y
		}
	}

	plt := plot.New()

	plt.Title.Text = "In-sample signal"
	plt.X.Label.Text = "Sample index"
	plt.Y.Label.Text = "Amplitude"
	plt.X.Min = 0
	plt.X.Max = float64(len(a.Samples) - 1)
	plt.Y.Min = minY - 1
	plt.Y.Max = maxY + 1

	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	plt.Add(line)

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "draw.png"); err != nil {
		panic(err)
	}
}

func signalTimeSpan(a modules.Signal) float64 {
	if len(a.Samples) == 0 || a.SampleRate <= 0 {
		return 0
	}
	return float64(len(a.Samples)-1) / a.SampleRate
}

func Output_inTime(a modules.Signal) {
	if len(a.Samples) == 0 || a.SampleRate <= 0 {
		return
	}

	points := make(plotter.XYs, len(a.Samples))
	minY := cmplx.Abs(a.Samples[0])
	maxY := cmplx.Abs(a.Samples[0])
	for i, sample := range a.Samples {
		points[i].X = float64(i) / a.SampleRate
		points[i].Y = cmplx.Abs(sample)

		if points[i].Y < minY {
			minY = points[i].Y
		}
		if points[i].Y > maxY {
			maxY = points[i].Y
		}
	}

	plt := plot.New()

	plt.Title.Text = "In-time signal"
	plt.X.Label.Text = "Time (s)"
	plt.Y.Label.Text = "Amplitude"
	plt.X.Min = 0
	plt.X.Max = signalTimeSpan(a)
	plt.Y.Min = minY - 1
	plt.Y.Max = maxY + 1

	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	plt.Add(line)

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "draw.png"); err != nil {
		panic(err)
	}
}
