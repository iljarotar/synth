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

If you load a file or change the current file, while the synth is playing, it will automatically adapt to the changes.

## Creating a patch

A patch is a yaml file, that specifies all of the synth's parameters.

The basic structure looks like this:

```yaml
volume: #decimal value
wavetables: #array of wavetables
  - amplitude:
      value: #decimal value
      modulation: #wavetable

    filters: #array of filters
      - type: #type of filter
        cutoff:
          value: #decimal value
          modulation: #wavetable
        ramp: #decimal value

        #another filter
      - type:
        cutoff:
        ramp:

    oscillators: #array of oscillators
      - type: #wave form
        freq: #decimal value
        amplitude:
          value: #decimal value
          modulation: #wavetable
        phase:
          value: #decimal value
          modulation: #wavetable

    #another wavetable
  - amplitude:
    filters:
    oscillators:
```

Most of the fields are optional. A simple sine wave with a frequency of 440 Hz can be generated like this:

```yaml
wavetables:
  - oscillators:
      - type: Sine
        freq: 440
```

You can find more examples in the `examples` directory.

---

### Synth

The topmost level of the synth has two fields:

`volume`  
is the synth's main volume.

`wavetables`  
is an array of wavetables.

---

### Wavetable

`amplitude`  
is a parameter, that has a `value` and an optional `modulation` field.

`filters`  
is an array of filters.

`oscillators`  
is an array of oscillators.

---

### Amplitude

`value`  
The value of a wavetable's amplitude affects the entire wavetable, whereas an oscillator's amplitude only affects that particular oscillator.

`modulation`  
is an optional wavetable, that modulates the amplitude.

### Oscillators

Types of oscillators:  
`Sine`  
`Square`  
`Triangle`  
`Sawtooth`  
`InvertedSawtooth`  
`Noise`

`amplitude`  
takes an initial value and an optional modulator wavetable. Modulating the amplitude results in a tremolo effect.

`phase`  
takes an initial value and an optional modulator wavetable. Modulating the phase results in a vibrato effect. Phase doesn't have any effect on an oscillator of type `Noise`.

### Filters

Types of filters  
`Lowpass`  
`Highpass`

`cutoff`  
takes an initial value and an optional modulator wavetable. The value specifies the frequency at which the filter sets in. Modulating the cutoff results in kind of a wah wah effect.

`ramp`  
specifies the distance between the cutoff frequency and the point, where a frequency isn't audible anymore.

Note: filters don't affect oscillators of type `Noise`.
