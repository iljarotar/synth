# Synth

This is a simple modular-like command line synthesizer for Linux written in
[golang](https://go.dev/).

## Installation

To run the synthesizer you might also need to install
[portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html).

### Install with Go

1. Install and setup [Go](https://go.dev/doc/install)
2. Clone the repository

```
git clone git@github.com:iljarotar/synth.git
```

3. Install

```
cd synth
go install
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

```
git clone git@github.com:iljarotar/synth.git
cd synth
synth -f examples/a-major.yaml
```

More examples can be found in the [examples](https://github.com/iljarotar/synth/tree/main/examples) directory.

> Note: If you want to record the output, you must specify a non-negative
> duration. Otherwise you will get am empty .wav file.

## Writing a patch file

### Data types

| Synth          |                     |                                                                                                           |
| -------------- | ------------------- | --------------------------------------------------------------------------------------------------------- |
| **Field**      | **Type**            | **Description**                                                                                           |
| vol            | Float               | main volume in range [0,1]                                                                                |
| out            | String [0..*]       | names of all oscillators, noise generators and custom signals, whose outputs will be sent to the speakers |
| time           | Float               | initial time shift in seconds [0,7200]                                                                    |
| oscillators    | Oscillator [0..*]   | all oscillators                                                                                           |
| noises         | Noise [0..*]        | all noise generators                                                                                      |
| custom-signals | CustomSignal [0..*] | all custom signals                                                                                        |
| envelopes      | Envelope [0..*]     | all envelopes                                                                                             |

| Oscillator |                |                                           |
| ---------- | -------------- | ----------------------------------------- |
| **Field**  | **Type**       | **Description**                           |
| name       | String         | should be unique in the scope of the file |
| type       | OscillatorType | wave form                                 |
| freq       | Param          | frequency in range [0,20000]              |
| amp        | Param          | amplitude in range [0,1]                  |
| phase      | Float          | phase in range [-1,1]                     |
| pan        | Param          | stereo balance in range [-1,1]            |

| OscillatorType  |
| --------------- |
| Sine            |
| Triangle        |
| Square          |
| Sawtooth        |
| ReverseSawtooth |

| Noise     |          |                                           |
| --------- | -------- | ----------------------------------------- |
| **Field** | **Type** | **Description**                           |
| name      | String   | should be unique in the scope of the file |
| amp       | Param    | amplitude in range [0,1]                  |
| pan       | Param    | stereo balance in range [-1,1]            |
| filter    | Filter   | a lowpass, highpass or bandpass filter    |

| Filter      |          |                                                            |
| ----------- | -------- | ---------------------------------------------------------- |
| **Field**   | **Type** | **Description**                                            |
| order       | Integer  | order of the FIR filter in range [0,1000]                  |
| low-cutoff  | Float    | cutoff frequency of the highpass filter in range [0,20000] |
| high-cutoff | Float    | cutoff frequency of the lowpass filter in range [0,20000]  |

If both `low-cutoff` and `high-cutoff` are 0, the filter is disabled. If
`low-cutoff` is 0, the filter is a lowpass filter transitioning at the
`high-cutoff` frequency. If `high-cutoff` is 0, the filter is a highpass filter
transitioning at the `low-cutoff` frequency.

A higher `order` improves the filter's precision, but it also makes it more
expensive in terms of computation. If the sound becomes glitchy, decreasing the
filter `order` might be necessary. Sometimes a lower value for `order` also sounds better than a high value.

| CustomSignal |              |                                           |
| ------------ | ------------ | ----------------------------------------- |
| **Field**    | **Type**     | **Description**                           |
| name         | String       | should be unique in the scope of the file |
| amp          | Param        | amplitude in range [0,1]                  |
| pan          | Param        | stereo balance in range [-1,1]            |
| freq         | Param        | periods per second [0,20000]              |
| data         | Float [0..*] | custom values                             |

| Envelope      |               |                                           |
| ------------- | ------------- | ----------------------------------------- |
| **Field**     | **Type**      | **Description**                           |
| name          | String        | should be unique in the scope of the file |
| attack        | Param         | attack time in seconds [0,10000]          |
| decay         | Param         | decay time in seconds [0,10000]           |
| sustain       | Param         | sustain time in seconds [0,10000]         |
| release       | Param         | release time in seconds [0,10000]         |
| peak          | Param         | peak amplitude [0,1]                      |
| sustain-level | Param         | sustain amplitude [0,1]                   |
| threshold     | Param         | trigger treshold [0,1]                    |
| triggers      | String [0..*] | names of triggering modules               |
| negative      | Boolean       | if true, the envelope's sign is inverted  |

| Param     |               |                                            |
| --------- | ------------- | ------------------------------------------ |
| **Field** | **Type**      | **Description**                            |
| val       | Float         | initial value of the respective parameter  |
| mod       | String [0..*] | names of modulating modules                |
| mod-amp   | Float         | amplitude of the modulation in range [0,1] |

### Structure of a patch file

```yaml
vol:
out:
noises:
  - name:
    amp:
      val:
      mod:
      mod-amp:
    pan:
      val:
      mod:
      mod-amp:
    filter:
      order:
      low-cutoff:
      high-cutoff:

oscillators:
  - name:
    type:
    freq:
      val:
      mod:
      mod-amp:
    amp:
      val:
      mod:
      mod-amp:
    pan:
      val:
      mod:
      mod-amp:
    phase:

custom-signals:
  - name:
    freq:
      val:
      mod:
      mod-amp:
    amp:
      val:
      mod:
      mod-amp:
    pan:
      val:
      mod:
      mod-amp:
    data: []

envelopes:
  - name:
    attack:
      val:
      mod:
      mod-amp:
    decay:
      val:
      mod:
      mod-amp:
    sustain:
      val:
      mod:
      mod-amp:
    release:
      val:
      mod:
      mod-amp:
    peak:
      val:
      mod:
      mod-amp:
    sustain-level:
      val:
      mod:
      mod-amp:
    threshold:
      val:
      mod:
      mod-amp:
    triggers: []
    negative:
```

Most of the fields are optional. A simple 440hz sine wave looks like this:

```yaml
vol: 1
out: [osc]
oscillators:
  - name: osc
    type: Sine
    freq: { val: 440 }
    amp: { val: 1 }
```
