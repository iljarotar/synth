# Command line synth in golang

## Installation

1. Get the binary [here](https://github.com/iljarotar/synth/releases)
2. Make it executable with `chmod +x synth`
3. Put it into some directory inside your `$PATH`
4. Install [portaudio](http://portaudio.com/)

To listen to an example clone this repository and run

```bash
synth -f examples/a-major.yaml
```

## Usage

Starting the synth with

```bash
synth
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

starts playing the file, if the format is correct. If you change and save the file, while the synth is running, it will instantly reload the patch.

Optionally you can pass a sample rate

```bash
synth -f <PATH_TO_YAML_FILE> -s 44100
```

44100 Hz is the default sample rate. 1000 Hz is the minimum.

---

## Creating a patch

The basic structure of a patch looks like this

```yaml
vol: # should not exceed 1
out: # list of oscillators
oscillators:
  - name: # choose any name
    type: # oscillator type
    freq: # list of frequencies
    pan:
      val: # initial pan in the interval [-1;1]
      mod: # list of oscillators
    amp:
      val: # amplitude value (should not exceed 1)
      mod: # list of oscillators
    phase:
      val: # initial phase shift
      mod: # list of oscillators
    filters: # list of filters

  # add as many oscillators as you need here

filters:
  - name: # choose any name
    ramp: # length of linear ramp
    low:
      val: # lower frequency limit
      mod: # list of oscillators
    high:
      val: # higher frequency limit
      mod: #list of oscillators
    vol:
      val: # volume of frequencies between low and high
      mod: # list of oscillators


  # add as many filters as you need here
```

### Synth

The `out` parameter is a list of oscillators, that will be sent to the speaker.

### Oscillators

Possible oscillator types are  
`Sine`  
`Square`  
`Triangle`  
`Sawtooth`  
`InvertedSawtooth`  
`Noise`

Stereo panning  
`-1` signal will be on left channel only  
`1` signal will be on right channel only

Everything in between will place the signal somewhere between the left and the right channel according to the ratio.

### Filters

All Filters are band pass filters with a frequency range between `low` and `high`.

### Examples

Most of the parameters are optional. A very simple patch may look like this

```yaml
volume: 1
out: [my-oscillator]
oscillators:
  - name: my-oscillator
    type: Sine
    freq: [440]
    amp:
      val: 1
```

### Modulation

The `mod` field of the `amp`, `phase`, `pan`, `low`, `high` or `vol` parameters is a list of oscillators, that will modulate that respective paramter. Here is an example of a tremolo effect.

```yaml
volume: 1
out: [my-oscillator]
oscillators:
  - name: my-oscillator
    type: Sine
    freq: [440]
    amp:
      val: 0.5
      mod: [modulator]

  - name: modulator
    type: Sine
    freq: [4]
    amp:
      val: 0.1
```
