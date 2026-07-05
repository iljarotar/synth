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
# My-patch.yaml.

# Main volume control.
# Range [0, 1].
vol: 1

# Name of the module to output.
out: name-of-main-module

# Global beats per minute used as the base tempo for all gate modules.
# Range [0, 400]
bpm: 120

# Delay effects.
delays:
  # The unique module name to be used as a reference in other modules.
  delay:
    # Delay time in milliseconds in range [0, 5000].
    time: 100

    # Gain of delayed signal in range [0, 1].
    gain: 0.25

    # Name of the module to run through the delay.
    in: name-of-input-module

    # CV for mix.
    cv: name-of-cv-module

    # Modulator for mix.
    mod: name-of-modulator

    # Fade controls the transition length in seconds.
    # Affected parameter is gain.
    fade: 2

# Adsr envelopes.
# Output values in range [0, 1].
envelopes:
  # The unique module name to be used as a reference in other modules.
  envelope:
    # Attack length in seconds.
    # Range [1e-15, 3600].
    attack: 0.1

    # Decay length in seconds.
    # Range [1e-15, 3600].
    decay: 0.05

    # Release length in seconds.
    # Range [1e-15, 3600].
    release: 2

    # Peak level targeted during the attack phase.
    # Range [0, 1].
    peak: 1

    # Sustain level.
    # Range [0, 1].
    level: 0.75

    # Name of the module to use as a gate signal.
    # When gate output changes from negative or zero to positive the envelope is triggered.
    # While gate is positive the attack, decay or sustain phases are active.
    # When gate output changes from positive to negative or zero the envelope is released.
    gate: name-of-gate-module

    # Fade controls the transition length in seconds.
    # Affected parameters are attack, decay, release, peak and level.
    fade: 2

# Filters of type low pass, high pass or band pass.
filters:
  # The unique module name to be used as a reference in other modules.
  filter:
    # One of LowPass, HighPass, BandPass.
    type: BandPass

    # Critical frequency.
    # Range [0, 20000].
    freq: 500

    # Band width in case of type BandPass.
    # Ignored for other types.
    width: 50

    # CV for freq.
    cv: name-of-cv

    # Modulator for freq.
    mod: name-of-modulator

    # Name of the module whose output will be filtered.
    in: name-of-input-module

    # Fade controls the transition length in seconds.
    # Affected parameters are freq and width.
    fade: 2

# Gates can be used as gates for envelopes or sequencers or as triggers for samplers.
gates:
  # The unique module name to be used as a reference in other modules.
  gate:
    # Number of sub divisions of a single beat in range [0, 128].
    # For example, a global bpm of 120 with 4 sub divisions results in 4*120 advancements of the signal per minute.
    # Since the signal has to change from 0 to one 1 or vice versa to actually trigger some other module, this would result in 2*120 possible beats.
    sub: 4

    # Binary signal.
    # Each negative or zero value will be mapped to -1, each positive to 1.
    signal: [1, 0, 0, 1, 0, 1, 1, 0, 1, 0]

    # Provides an initial offset to the signal.
    # Count starts at 0.
    # Index will not be updated while playing.
    index: 4

# Mixers combine outputs of multiple modules and control their output levels.
mixers:
  # The unique module name to be used as a reference in other modules.
  mixer:
    # Gain in range [0, 1].
    gain: 0.5

    # CV for gain.
    cv: name-of-cv

    # Modulator for gain.
    mod: name-of-modulator

    # Mapping of module names to their corresponding gain levels.
    # These gain levels can be in the range [0, 1000].
    # The output will be limited not to exceed the range [-1, 1].
    in:
      name-of-first-module: 0.5
      name-of-second-module: 0.25

    # Fade controls the transition length in seconds.
    # Affected parameters are gain as well as all input modules' gain levels.
    fade: 2

# Noise modules simple output random values.
noises:
  # The unique module name to be used as a reference in other modules.
  # A noise module doesn't have any parameters to configure, so pass an empty object {}.
  noise: {}

# Oscillators output basic wave forms like sine waves, triangles, etc..
oscillators:
  # The unique module name to be used as a reference in other modules.
  oscillator:
    # One of Sine, Square, Triangle, Sawtooth, ReverseSawtooth.
    type: Sine

    # Frequency in range [0, 20000].
    freq: 440

    # CV for freq.
    cv: name-of-cv

    # Modulator for freq.
    # Maximum amount of modulation is one octave up and down.
    mod: name-of-mod

    # Static phase shift in percent of one period.
    # Range [-1, 1].
    phase: 0.75

    # Fade controls the transition length in seconds.
    # Affected parameters are freq and phase.
    fade: 2

# Pan modules are used to add stereo balance.
pans:
  # The unique module name to be used as a reference in other modules.
  pan:
    # Specifies how big a portion of the signal is output through the left and right channels.
    # Range [1-, 1].
    # A value of -1 places the signal completely to the left, 1 places it to the right.
    pan: -0.5

    # Name of the module whose output should be stereo balanced.
    in: name-of-input-module

    # Modulator for pan.
    mod: name-of-mod

    # Fade controls the transition length in seconds.
    # Affected parameter is pan.
    fade: 2

# Sample and hold modules.
samplers:
  # The unique module name to be used as a reference in other modules.
  sampler:
    # Name of the module that is being sampled.
    in: name-of-input-module

    # Name of the trigger module.
    # When the trigger's output value changes from negative of zero to positive a new sample is taken from the input modules output.
    trigger: name-of-trigger-module

# Sequencers can be combined with oscillators or wavetables to create melodic sequences.
# Output values in range [0, 1].
sequencers:
  # The unique module name to be used as a reference in other modules.
  sequencer:
    # A sequence of notes in scientific pitch notation.
    # Flats are denoted by 'b', sharps by '#'.
    # A note is separated from its octave by an underscore.
    # Minimum octave is 0, maximum is 10.
    sequence: ["a_4", "eb_3", "c#_5"]

    # When the trigger's value changes from negative or zero to positive the next note in the sequence is triggered.
    trigger: name-of-trigger-module

    # Base pitch from which to calculate all other frequencies.
    pitch: 440

    # Transpose the whole sequence by any number of semitones.
    # Range [-24, 24].
    transpose: -4

    # If true the notes in the sequence will be played in random order.
    randomize: true

    # Provides an initial offset to the sequence.
    # Count starts at 0.
    index: 2

# Pass any values to a wavetable to create arbitrary signals.
wavetables:
  # The unique module name to be used as a reference in other modules.
  wavetable:
    # Frequency in range [0, 20000].
    # Specifies how many times per second the entire signal given in the signal field will be played.
    freq: 440

    # CV for freq.
    cv: name-of-cv

    # Modulator for freq.
    mod: name-of-mod

    # An arbitrary signal.
    # The signal can have any length.
    signal: [-1, 0, 0.25, -0.3, 0.8, 1]

    # Fade controls the transition length in seconds.
    # Affected parameter is freq.
    fade: 2
```

### Configuration

On first run, synth will create a `synth/config.yaml` file in your default config directory.
Modify this file to adjust the configuration.
Run `synth -h` to see where this file was placed.
You can also override single parameters via command line flags.
