# Synth

This is a simple modular-like command line synthesizer for Linux written in
[golang](https://go.dev/).

## Installation

To run the synthesizer you might also need to install
[portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html).

### Install with Go

1. Install and setup [Go](https://go.dev/doc/install)
2. Clone the repository

```bash
git clone git@github.com:iljarotar/synth.git
```

3. Install

```bash
cd synth
make build
cp bin/synth <SOMEWHERE_IN_YOUR_PATH>
```

## Usage

### How it works

The synth is programmed by providing a patch file. A patch is a `yaml` file,
that tells the synth, which oscillators, noise generators, etc. should be
created and how they should be connected. When you tell the synth to load a
file, if will start playing immediately. During playback a hot reload is
possible, so if you change and save the patch file, it will be applied
instantly.

### Command line interface

To listen to an example run

```bash
synth examples/a-major.yaml
```

Run `synth -h` to see all configuration options.

## Writing a patch file

### Data types

| Synth       |                  |                                                    |
| ----------- | ---------------- | -------------------------------------------------- |
| **Field**   | **Type**         | **Description**                                    |
| vol         | Float            | main volume in range [0,2]                         |
| out         | String[0..*]     | names of all modules whose outputs will be audible |
| time        | Float            | initial time shift in seconds [0,7200]             |
| oscillators | Oscillator[0..*] | all oscillators                                    |
| noises      | Noise[0..*]      | all noise generators                               |
| wavetables  | Wavetables[0..*] | all wavetables                                     |
| samplers    | Sampler[0..*]    | all samplers                                       |

| Oscillator |                |                                                                 |
| ---------- | -------------- | --------------------------------------------------------------- |
| **Field**  | **Type**       | **Description**                                                 |
| name       | String         | should be unique in the scope of the file                       |
| type       | OscillatorType | wave form                                                       |
| freq       | Input          | frequency in range [0,20000]                                    |
| amp        | Input          | amplitude in range [0,2]                                        |
| phase      | Float          | phase in range [-1,1]                                           |
| pan        | Input          | stereo balance in range [-1,1]                                  |
| filters    | String[0..*]   | names of the filters to apply                                   |
| envelope   | Envelope       | envelope to apply; if omitted, oscillator will constantly sound |

| OscillatorType  |
| --------------- |
| Sine            |
| Triangle        |
| Square          |
| Sawtooth        |
| ReverseSawtooth |

| Noise     |              |                                                            |
| --------- | ------------ | ---------------------------------------------------------- |
| **Field** | **Type**     | **Description**                                            |
| name      | String       | should be unique in the scope of the file                  |
| amp       | Input        | amplitude in range [0,2]                                   |
| pan       | Input        | stereo balance in range [-1,1]                             |
| filters   | String[0..*] | names of the filters to apply                              |
| envelope  | Envelope     | envelope to apply; if omitted, noise will constantly sound |

| Wavetable |              |                                                                |
| --------- | ------------ | -------------------------------------------------------------- |
| **Field** | **Type**     | **Description**                                                |
| name      | String       | should be unique in the scope of the file                      |
| amp       | Input        | amplitude in range [0,2]                                       |
| pan       | Input        | stereo balance in range [-1,1]                                 |
| freq      | Input        | periods per second [0,20000]                                   |
| table     | Float[0..*]  | output values                                                  |
| filters   | String[0..*] | names of the filters to apply                                  |
| envelope  | Envelope     | envelope to apply; if omitted, wavetable will constantly sound |

| Sampler   |              |                                                              |
| --------- | ------------ | ------------------------------------------------------------ |
| **Field** | **Type**     | **Description**                                              |
| name      | String       | should be unique in the scope of the file                    |
| amp       | Input        | amplitude in range [0,2]                                     |
| pan       | Input        | stereo balance in range [-1,1]                               |
| freq      | Input        | frequency in range [0,SAMPLE_RATE (default 44100)]           |
| filters   | String[0..*] | names of the filters to apply                                |
| inputs    | String[0..*] | names of the modules that will be sampled                    |
| envelope  | Envelope     | envelope to apply; if omitted, sampler will constantly sound |

A sampler periodically samples the output values of the given inputs and outputs their sum.

| Sequence  |                |                                                                                                                     |
| --------- | -------------- | ------------------------------------------------------------------------------------------------------------------- |
| **Field** | **Type**       | **Description**                                                                                                     |
| name      | String         | should be unique in the scope of the file                                                                           |
| amp       | Input          | amplitude in range [0,2]                                                                                            |
| pan       | Input          | stereo balance in range [-1,1]                                                                                      |
| type      | OscillatorType | wave form                                                                                                           |
| sequence  | String[0..*]   | a sequence of notes written in [scientific pitch notation](https://en.wikipedia.org/wiki/Scientific_pitch_notation) |
| randomize | Boolean        | if true, the notes of the sequence will be played in random order                                                   |
| pitch     | Float          | standard pitch in hz [400,500]                                                                                      |
| transpose | Input          | transposition in semitones [-24,24]                                                                                 |
| filters   | String[0..*]   | names of the filters to apply                                                                                       |
| envelope  | Envelope       | envelope to apply; if omitted, first note of sequence will constantly sound                                         |

| Filter      |          |                                                            |
| ----------- | -------- | ---------------------------------------------------------- |
| **Field**   | **Type** | **Description**                                            |
| low-cutoff  | Float    | cutoff frequency of the highpass filter in range [1,20000] |
| high-cutoff | Float    | cutoff frequency of the lowpass filter in range [1,20000]  |

If both `low-cutoff` and `high-cutoff` are omitted, the filter is disabled. If
`low-cutoff` is omitted, the filter is a lowpass filter transitioning at the
`high-cutoff` frequency. If `high-cutoff` is omitted, the filter is a highpass filter
transitioning at the `low-cutoff` frequency. If both cutoff frequencies are defined, it becomes a bandpass filter.

| Envelope      |          |                                   |
| ------------- | -------- | --------------------------------- |
| **Field**     | **Type** | **Description**                   |
| attack        | Input    | attack time in seconds [0,10000]  |
| decay         | Input    | decay time in seconds [0,10000]   |
| sustain       | Input    | sustain time in seconds [0,10000] |
| release       | Input    | release time in seconds [0,10000] |
| peak          | Input    | peak amplitude [0,1]              |
| sustain-level | Input    | sustain amplitude [0,1]           |
| bpm           | Input    | triggers per minute [0,600000]    |
| time-shift    | Float    | initial time shift                |

| Input     |              |                                                                         |
| --------- | ------------ | ----------------------------------------------------------------------- |
| **Field** | **Type**     | **Description**                                                         |
| val       | Float        | initial value of the respective parameter                               |
| mod       | String[0..*] | names of modulating modules (oscillators, samplers, wavetables, noises) |
| mod-amp   | Float        | amplitude of the modulation in range [0,1]                              |

## Example patch file

Playing this file outputs a 440 Hz sine wave.

```yaml
vol: 1
out: [osc]
oscillators:
  - name: osc
    type: Sine
    freq: { val: 440 }
    amp: { val: 1 }
```
