# Command line synth in golang

## Installation

To install the synth [Golang](https://go.dev/doc/install) is required. Once you have Go installed and set up, run

```bash
git clone git@github.com:iljarotar/synth.git
cd synth
go install
```

Listen to an example

```bash
synth -f examples/a-major.yaml
```

## Usage

Starting the synth with

```bash
synth -h
```

outputs

```
command line synthesizer

documentation and usage: https://github.com/iljarotar/synth

Usage:
  synth [flags]

Flags:
  -f, --file string          specify which file to load
  -h, --help                 print help
  -s, --sample-rate string   specify sample rate
```

Loading a file with

```bash
synth -f <PATH_TO_YAML_FILE>
```

starts playing the file, if the format is correct.

Optionally you can pass a sample rate

```bash
synth -f <PATH_TO_YAML_FILE> -s 44100
```

44100 Hz is the default sample rate. 1000 Hz is the minimum.

---

## Creating a patch

The basic structure of a patch looks like this

```yaml
volume: # should not exceed 1
out: [# array of oscillators]
oscillators:
  - name: # choose any name
    type: # oscillator type (see below)
    freq: # frequency
    amp:
      val: # amplitude value (should not exceed 1)
      mod: [# array of oscillators]
    phase:
      val: # initial phase shift
      mod: [# array of oscillators]
    filters: [# array of filters]

  # add as many oscillators as you need here

filters:
  - name: # choose any name
    type: # type of filter
    ramp: # length of linear ramp
    cutoff:
      val: # initial cutoff frequency
      mod: [# array of oscillators]

  # add as many filters as you need here
```

The `out` parameter is an array of oscillators, that will be sent to the speaker.

Possible oscillator types are  
`Sine`  
`Square`  
`Triangle`  
`Sawtooth`  
`InvertedSawtooth`  
`Noise`

Possible filter types are  
`Lowpass`  
`Highpass`

Most of the parameters are optional. The most simple patch may look like this

```yaml
volume: 1
out: [my-oscillator]
oscillators:
  - name: my-oscillator
    type: Sine
    freq: 440
    amp:
      val: 1
```

The `mod` field of the `amp`, `phase` or `cutoff` parameters is a list of oscillators, that will modulate that respective paramter. Here is an example of a tremolo effect.

```yaml
volume: 1
out: [my-oscillator]
oscillators:
  - name: my-oscillator
    type: Sine
    freq: 440
    amp:
      val: 0.5
      mod: [modulator]

  - name: modulator
    type: Sine
    freq: 4
    amp:
      val: 0.1
```
