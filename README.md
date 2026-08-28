# DigitalBeamForming_Search_RADAR
An typycal Search RADAR simulation that use DigitalBeamForming in receive


## Signal Generator
Generation Chirp LFM base this formula : $exp(j * pi * \frac{Bandwidth}{PulseWidth} * t^2)$

the instance bandwidth sweep between $[\frac{-Bandwidth}{2} :\frac{+Bandwidth}{2}]$

SampleRate, PulseWidth and Bandwidth read from config.yaml file.

## Target Model
Target model is contain initial lication and velocity in cartesian and RCS(Radar cross section). the read from config.yaml file. one or more targets can describe.

## Echo Model
base on each target situation and envirment, the echoed pulse delayed and change amp and frequncy.

$delaySec = \frac{2*r}{SpeedOfLight}$ ,where r is the range.

$amp = 

$fd = 2 * V_r / \lambda$ , where $V_r$ is radial velocity and $\lambda$ is the radiated wave length base radar frequency that raed from config.yaml file.
