# DigitalBeamForming_Search_RADAR
An typycal Search RADAR simulation that use DigitalBeamForming in receive


## Signal Generator
Generation Chirp LFM base this formula : $exp(j * \pi * \frac{Bandwidth}{PulseWidth} * t^2)$

the instance bandwidth sweep between $[\frac{-Bandwidth}{2} :\frac{+Bandwidth}{2}]$

SampleRate, PulseWidth and Bandwidth read from config.yaml file.

## Target Model
Target model is contain initial lication and velocity in cartesian and RCS(Radar cross section). the read from config.yaml file. one or more targets can describe.

## Echo Model
base on each target situation and envirment, the echoed pulse delayed and change amp and frequncy.

$delaySec = \frac{2*r}{SpeedOfLight}$ ,where r is the range.

$amp = 

$fd = 2 * V_r / \lambda$ , where $V_r$ is radial velocity and $\lambda$ is the radiated wave length base radar frequency that raed from config.yaml file.

## Antenna and DBF

Antann Stacked N element in y-axis that each element conncetd to T/R module. with Digtal Weightening Matrix we make M beam in Elevation.

we assume 1-D planar array with Gaussain Pattern elements with ElemAzBWDeg and ElemElBWDeg 3dB Beamwidth.  
Antenna element and array configuration read from config.yaml file.


### Mathematics:
the Complete Antenna is:  
$G(\theta,\phi)=G_e(\theta,\phi)|AF(\theta,\phi)|^2$  

The received signal at element (n) has a phase term:  
$exp(jkz_nsin(\phi))$ where $k=2\frac{\pi}{\lambda}$ and $z_n=(n-\frac{N-1}{2})d$ where d is element spacing.

For pointing beam to $\phi_0$ :  
$y=\sum_{n=0}^{N-1} w_n x_n$ where $w_n=exp(-jkz_nsin(\phi_0))$

For a uniform linear array, the approximate half-power beamwidth is:  
$HPBW=\frac{0.886\lambda}{Ndcos(\phi_0)}$ in radian

for Build M simlutantly Beam we need $W_{N*M}$ complex beamforming matrix.

how many DBF beams need to cover an angular sector S:  
$M\approx 1+\frac{S}{K HPBW}$ where k is overlap factor(for example 1:no overlap ,0.5: 0.5 HPBW)



