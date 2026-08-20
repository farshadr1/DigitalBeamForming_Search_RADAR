package modules

import "math"

type Target struct {
	ID int

	X float64 // meters
	Y float64 // meters
	Z float64 // meters

	VelocityX float64 // m/s
	VelocityY float64 // m/s
	VelocityZ float64 // m/s

	RCS float64 // m²
}

// Range
func (t Target) Range() float64 {
	return math.Sqrt(t.X*t.X + t.Y*t.Y + t.Z*t.Z)
}

// Azimuth in Degree
func (t Target) Azimuth() float64 {
	return math.Atan2(t.Y, t.X) * 180 / math.Pi
}

// Elevation in Degree
func (t Target) Elevation() float64 {
	r := t.Range()
	return math.Asin(t.Z/r) * 180 / math.Pi
}

// Radial Speed
func (t Target) RadialVelocity() float64 {
	r := t.Range()

	return (t.X*t.VelocityX +
		t.Y*t.VelocityY +
		t.Z*t.VelocityZ) / r
}

// Update moves the target forward by dt seconds.
func (t *Target) Update(dt float64) {
	t.X += t.VelocityX * dt
	t.Y += t.VelocityY * dt
	t.Z += t.VelocityZ * dt
}
