
# Digital Beamforming Search Radar

A compact Go-based simulation of a search radar that uses digital beamforming (DBF) on receive. The project simulates LFM chirp transmission, multiple targets, echo generation (including range delay, amplitude scaling and Doppler), and a 1-D antenna array with DBF producing multiple elevation beams.

## Features

- LFM (chirp) signal generator
- Configurable sample rate, pulse width and bandwidth via `config.yaml`
- Multiple target model with position, velocity and RCS
- Echo model with range delay, amplitude scaling and Doppler shift
- 1-D planar antenna array and digital beamforming to form multiple beams

## Configuration

Simulation parameters (SampleRate, PulseWidth, Bandwidth, antenna geometry, radar carrier frequency, targets, etc.) are read from `config.yaml`.

## Signal generator

The transmitted LFM chirp follows the complex baseband expression:

$$s(t)=\exp\left(j\pi\frac{\mathrm{Bandwidth}}{\mathrm{PulseWidth}}t^2\right)$$

The instantaneous frequency sweeps across $[ -\mathrm{Bandwidth}/2, +\mathrm{Bandwidth}/2 ]$.

## Target model

Each target is described by Cartesian position, velocity and an RCS value. One or more targets can be defined in `config.yaml`.

## Echo model

For a target at range $r$ and radial velocity $V_r$ the main effects on the received pulse are:

- Range delay: $\mathrm{delaySec}=\dfrac{2r}{c}$, where $c$ is the speed of light.
- Doppler frequency shift: $f_d=\dfrac{2V_r}{\lambda}$, where $\lambda$ is the radar wavelength (from carrier frequency).
- Amplitude scaling according to range and RCS (implementation-specific in `echo.go`).

## Antenna array and digital beamforming

This project assumes a 1-D uniform linear array (elements stacked along the y-axis) with element pattern approximated by a Gaussian 3-dB beamwidth in azimuth and elevation (`ElemAzBWDeg`, `ElemElBWDeg`). The array and element parameters are read from `config.yaml`.

The overall array gain is modeled as:

$$G(\theta,\phi)=G_e(\theta,\phi)\,|AF(\theta,\phi)|^2$$

The phase at element $n$ for a plane wave arriving from angle $\phi$ uses the steering term

$$\exp\big(jk z_n\sin\phi\big)\quad\text{with }k=\dfrac{2\pi}{\lambda},\;z_n=\big(n-\tfrac{N-1}{2}\big)d$$

To steer a beam to angle $\phi_0$ the array weights are

$$w_n=\exp\big(-jk z_n\sin\phi_0\big)\qquad y=\sum_{n=0}^{N-1} w_n x_n$$

An approximate half-power beamwidth (HPBW) for a uniform linear array is

$$\mathrm{HPBW}\approx\dfrac{0.886\,\lambda}{N d \cos\phi_0}$$

To form $M$ simultaneous beams a complex beamforming matrix $W_{N\times M}$ is used. The number of DBF beams required to cover an angular sector $S$ can be estimated as

$$M\approx 1+\dfrac{S}{K\,\mathrm{HPBW}}$$

where $K$ is an overlap factor (e.g. $K=1$ no overlap, $K=0.5$ for 50% overlap).

The element power pattern is approximated with a Gaussian:

$$G(\theta)=\exp\left(-4\ln 2\,\left(\dfrac{\theta}{\theta_{3\mathrm{dB}}}\right)^2\right)$$

## Files of interest

- `main.go` — simulation entry point
- `config.yaml` — simulation configuration
- `signal.go` — signal generation
- `target.go` — target definitions
- `echo.go` — echo (channel) model
- `fig.go` / `figures/` — plotting helpers

## Running

Build and run the simulation with the standard Go toolchain:

```bash
go build ./...
./DigitalBeamForming_Search_RADAR
```

Adjust `config.yaml` to modify simulation parameters.

---
If you want, I can further tailor this README (examples, diagrams, or usage samples). 


