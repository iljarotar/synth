# Command line synth in golang

## Installation

To install the synth [Golang](https://go.dev/doc/install) is required. Once you have Go installed and set up, run:

```bash
git clone git@github.com:iljarotar/synth.git
cd synth
make install
```

## Configuration

After running `make install` you should have a file called `~/.config/synth/config.yaml`, which specifies the sample rate and the root path, where the synth will look for patch files. You can change the root path from within the synth (see below). The sample rate can only be changed manually and will be applied on startup.

## Usage

Start the synth:

```bash
synth
```

The output should look like this

```
Menu

apply, a        reload last file
clear, c        clear screen
exit, e         exit synth
help, h         print this menu
load, l         load file
play, p         start synth
root, r         specify root path
stop, s         stop synth

:
```

### Load patches

By default the synth will look for files in a directory called `~/synth`. If you want to use this default, create the directory and place all your patches there.

If you want to use a different directory, start the synth and type

```bash
root <PATH-TO-YOUR-PATCHES>
```

load a file

```bash
load <PATH-TO-YOUR-FILE-RELATIVE-TO-YOUR-ROOT-PATH>
```

load the same file again

```bash
apply
```

Note: if you load a file, while the synth is playing, it will instantaneously adapt.

## Creating a patch

A patch is a yaml-file and may look like this

```yaml
gain: 1
wavetables:
  - filters:
      - type: "Highpass"
        cutoff: 300
        ramp: 110
        cutoff-mod:
          oscillators:
            - type: "Sine"
              amplitude: 100
              freq: 0.25
      - type: "Lowpass"
        cutoff: 300
        ramp: 110
        cutoff-mod:
          oscillators:
            - type: "Sine"
              amplitude: 100
              freq: 0.25
    oscillators:
      - type: "Sine"
        amplitude: 0.5
        freq: 220
        am:
          oscillators:
            - type: "Sine"
              amplitude: 0.5
              freq: 1
            - type: "Sine"
              amplitude: 0.5
              freq: 2

      - type: "Sine"
        amplitude: 0.5
        freq: 275
        pm:
          oscillators:
            - type: "Sine"
              amplitude: 100
              freq: 1
            - type: "Sine"
              amplitude: 30
              freq: 1

      - type: "Sine"
        amplitude: 0.5
        freq: 330
        am:
          oscillators:
            - type: "Sine"
              amplitude: 0.5
              freq: 1

      - type: "Sine"
        amplitude: 0.5
        freq: 415
        pm:
          oscillators:
            - type: "Sine"
              amplitude: 1
              freq: 1
            - type: "Sine"
              amplitude: 1.2
              freq: 0.25
```

The basic structure is

```yaml
gain:
wavetables:
  - filters:
    oscillators:
  - filters:
    oscillators:
```

`gain` is the synthesizer's main volume, `wavetables` is an array of wavetables. Each wavetable may have multiple filters and oscillators.

### Oscillators

| type             | amplitude | freq        | phase       | am       | pm          |
| ---------------- | --------- | ----------- | ----------- | -------- | ----------- |
| Sine             | required  | required    | optional    | optional | optional    |
| Square           | required  | required    | optional    | optional | optional    |
| Sawtooth         | required  | required    | optional    | optional | optional    |
| InvertedSawtooth | required  | required    | optional    | optional | optional    |
| Triangle         | required  | required    | optional    | optional | optional    |
| Noise            | required  | ineffective | ineffective | optional | ineffective |

Examples

```yaml
oscillators:
  - type: Square
    amplitude: 0.75
    freq: 440
    phase: 0.5
  - type: Sine
    amplitude: 1
    freq: 220
  - type: Noise
    amplitude: 0.05
```

Optional amplitude and phase modulation can be added to every oscillator individually. A Modulator is itself a wavetable with filters and oscillators, so arbitrarily deep nesting is possible.

Examples

```yaml
oscillators:
  - type: Sine
    amplitude: 0.5
    freq: 220
    am:
      oscillators:
        - type: "Sine"
          amplitude: 0.5
          freq: 1

  - type: Square
    amplitude: 1
    freq: 275
    pm:
      oscillators:
        - type: Noise
          amplitude: 0.5
```

Note: Phase modulation effectively results in pitch modulation.

### Filters

A filter can be of type `Lowpass` or `Highpass`. The basic structure is

```yaml
filters:
  - type: "Highpass"
    cutoff: 300
    ramp: 110
    cutoff-mod:
      oscillators:
        - type: "Sine"
          amplitude: 100
          freq: 0.25
```

`cutoff` specifies the cutoff frequency and `ramp` specifies the length of the linear ramp from the cutoff frequency to the point, where a frequency is filtered out entirely. For example, if `cutoff` is `440` and `ramp` is `60`, then 480 Hz will be softer then 440 Hz and 500 Hz and higher will not be audible any more.

`cutoff-mod` is a wavetable, that modulates the `cutoff`.

Note: unfortunately filters don't affect oscillators of type `Noise`.
