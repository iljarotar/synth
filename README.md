# Synth

A modular synthesizer for the command line written in [golang](https://go.dev/).

## Installation

1. Install and setup [Go](https://go.dev/doc/install)
2. Clone the repository

```bash
git clone git@github.com:iljarotar/synth.git
```

3. Install

```bash
cd synth
make build
cp bin/synth /usr/local/bin # or somewhere else in your PATH
```

## Usage

> Run `synth examples/sequencer.yaml` to listen to an example.

Patches for the modular synthesizer are provided in [yaml](https://yaml.org/) format.
A patch contains configurations for all modules that you want the synthesizer to play.
When you modify and save a patch during playback the synthesizer will reload the file.
Since changing a parameter like volume or frequency by a large amount too quickly results in a clipping noise, most modules allow configuring a `fade` parameter that controls how long it takes for the module's parameters to transition from the previous value to the new one.
Such a fade-over is not only useful to avoid clipping sounds but can also be utilised to create slow transitions in the music.
Say, for example, you want to slowly fade in one module while slowly fading out another.
You can add a `fade` parameter to the mixer that controls both modules' volumes—e.g. `fade: 5` for 5 seconds—and change the new module's volume to a positive value and the other one's to `0`.
Then save the file and the transition will start.

### Patch Files

This section explains all available modules and provides example configurations.
For more examples see the [examples](examples/) directory.

Each module must have a unique name across all modules.
This name is used as a reference in other modules, e.g. when a module is used as a CV or modulator.
Each module outputs values in the interval `[-1, 1]`.
Additionally, all parameters of a module are limited not to extend the reasonable ranges for each specific parameter, e.g. an oscillator's frequency will never exceed 20,000Hz.
Such limitations make the outcome of a configuration more predictable.
Those modules whose main purpose is to provide CV values for other modules only output values between `0` and `1`.

#### CV

If a module is provided with a CV its static value is ignored.
For example, if you pass a CV to an oscillator this CV will provide the oscillator's frequency and the statically assigned frequency will be ignored.
When mapping a CV-provider's output to a parameter, only positive values are considered.
So a value in the interval `[0, 1]` is mapped to the respective parameter's range.
For example, a sequencer will output values in the range `[0, 1]` and those values will be mapped to an oscillator's frequency range `[0, 20000]` if the sequencer is used as a CV for that oscillator.

#### Modulation

If a module is provided with a modulator it will modulate a parameter around its static or CV-provided value.
For example, a mixer with a gain of `0.5` and a sine wave as a modulator will output a tremolo around the gain value of `0.5`.
For example, an oscillator's frequency is in the range `[0, 20000]`.
A modulator that outputs values in the entire possible range of `[-1, 1]` will modulate the oscillator's frequency in the entire range `[0, 20000]`.
To control the amount of modulation you must send the modulator through a mixer and attenuate its gain.

#### Module Reference

The following yaml file provides examples and explanations for all configuration parameters.

```yaml
# my-patch.yaml

# main volume control
# range [0, 1]
vol: 1

# name of the module to output
out: name-of-main-module

# delay effects
delays:
  # the unique module name to be used as a reference in other modules
  delay:
    # delay time in milliseconds in range [0, 5000]
    time: 100

    # gain of delayed signal in range [0, 1]
    gain: 0.25

    # name of the module to run through the delay
    in: name-of-input-module

    # cv for mix
    cv: name-of-cv-module

    # modulator for mix
    mod: name-of-modulator

    # fade controls the transition length in seconds
    # affected parameter is gain
    fade: 2

# adsr envelopes
# output values in range [0, 1]
envelopes:
  # the unique module name to be used as a reference in other modules
  envelope:
    # attack length in seconds
    # range [1e-15, 3600]
    attack: 0.1

    # decay length in seconds
    # range [1e-15, 3600]
    decay: 0.05

    # release length in seconds
    # range [1e-15, 3600]
    release: 2

    # peak level targeted during the attack phase
    # range [0, 1]
    peak: 1

    # sustain level
    # range [0, 1]
    level: 0.75

    # name of the module to use as a gate signal
    # when gate output changes from negative or zero to positive the envelope is triggered
    # while gate is positive the attack, decay or sustain phases are active
    # when gate output changes from positive to negative or zero the envelope is released
    gate: name-of-gate-module

    # fade controls the transition length in seconds
    # affected parameters are attack, decay, release, peak and level
    fade: 2

# filters of type low pass, high pass or band pass
filters:
  # the unique module name to be used as a reference in other modules
  filter:
    # one of LowPass, HighPass, BandPass
    type: BandPass

    # critical frequency
    # range [0, 20000]
    freq: 500

    # band width in case of type BandPass
    # ignored for other types
    width: 50

    # cv for freq
    cv: name-of-cv

    # modulator for freq
    mod: name-of-modulator

    # name of the module whose output will be filtered
    in: name-of-input-module

    # fade controls the transition length in seconds
    # affected parameters are freq and width
    fade: 2

# gates can be used as gates for envelopes or sequencers or as triggers for samplers.
gates:
  # the unique module name to be used as a reference in other modules
  gate:
    # beats per minute controls the tempo of the gate signal
    bpm: 260

    # cv for bpm
    cv: name-of-cv

    # modulator for bpm
    mod: name-of-modulator

    # binary signal
    # each negative or zero value will be mapped to -1, each positive to 1
    signal: [1, 0, 0, 1, 0, 1, 1, 0, 1, 0]

    # provides an initial offset to the signal
    # count starts at 0
    index: 4

    # fade controls the transition length in seconds
    # affected parameter is bpm
    fade: 2

# mixers combine outputs of multiple modules and control their output levels
mixers:
  # the unique module name to be used as a reference in other modules
  mixer:
    # gain in range [0, 1]
    gain: 0.5

    # cv for gain
    cv: name-of-cv

    # modulator for gain
    mod: name-of-modulator

    # mapping of module names to their corresponding gain levels
    # these gain levels can be in the range [0, 1000]
    # the output will be limited not to exceed the range [-1, 1]
    in:
      name-of-first-module: 0.5
      name-of-second-module: 0.25

    # fade controls the transition length in seconds
    # affected parameters are gain as well as all input modules' gain levels
    fade: 2

# noise modules simple output random values
noises:
  # the unique module name to be used as a reference in other modules
  # a noise module doesn't have any parameters to configure, so pass an empty object {}
  noise: {}

# oscillators output basic wave forms like sine waves, triangles, etc.
oscillators:
  # the unique module name to be used as a reference in other modules
  oscillator:
    # one of Sine, Square, Triangle, Sawtooth, ReverseSawtooth
    type: Sine

    # frequency in range [0, 20000]
    freq: 440

    # cv for freq
    cv: name-of-cv

    # modulator for freq
    # maximum amount of modulation is one octave up and down
    mod: name-of-mod

    # static phase shift in percent of one period
    # range [-1, 1]
    phase: 0.75

    # fade controls the transition length in seconds
    # affected parameters are freq and phase
    fade: 2

# pan modules are used to add stereo balance
pans:
  # the unique module name to be used as a reference in other modules
  pan:
    # specifies how big a portion of the signal is output through the left and right channels
    # range [1-, 1]
    # a value of -1 places the signal completely to the left, 1 places it to the right
    pan: -0.5

    # name of the module whose output should be stereo balanced
    in: name-of-input-module

    # modulator for pan
    mod: name-of-mod

    # fade controls the transition length in seconds
    # affected parameter is pan
    fade: 2

# sample and hold modules
samplers:
  # the unique module name to be used as a reference in other modules
  sampler:
    # name of the module that is being sampled
    in: name-of-input-module

    # name of the trigger module
    # when the trigger's output value changes from negative of zero to positive a new sample is taken from the input modules output
    trigger: name-of-trigger-module

# sequencers can be combined with oscillators or wavetables to create melodic sequences
# output values in range [0, 1]
sequencers:
  # the unique module name to be used as a reference in other modules
  sequencer:
    # a sequence of notes in scientific pitch notation
    # flats are denoted by 'b', sharps by '#'
    # a note is separated from its octave by an underscore
    # minimum octave is 0, maximum is 10
    sequence: ["a_4", "eb_3", "c#_5"]

    # when the trigger's value changes from negative or zero to positive the next note in the sequence is triggered
    trigger: name-of-trigger-module

    # base pitch from which to calculate all other frequencies
    pitch: 440

    # transpose the whole sequence by any number of semitones
    # range [-24, 24]
    transpose: -4

    # if true the notes in the sequence will be played in random order
    randomize: true

    # provides an initial offset to the sequence
    # count starts at 0
    index: 2

# pass any values to a wavetable to create arbitrary signals
wavetables:
  # the unique module name to be used as a reference in other modules
  wavetable:
    # frequency in range [0, 20000]
    # specifies how many times per second the entire signal given in the signal field will be played
    freq: 440

    # cv for freq
    cv: name-of-cv

    # modulator for freq
    mod: name-of-mod

    # an arbitrary signal
    # the signal can have any length
    signal: [-1, 0, 0.25, -0.3, 0.8, 1]

    # fade controls the transition length in seconds
    # affected parameter is freq
    fade: 2
```

### Configuration

On first run, synth will create a `synth/config.yaml` file in your default config directory.
Modify this file to adjust the configuration.
Run `synth -h` to see where this file was placed.
You can also override single parameters via command line flags.
