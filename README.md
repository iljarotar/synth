# Synth

This is a simple modular-like command line synthesizer written in
[golang](https://go.dev/).

## Installation

> Note: I have tested the synth only on Fedora. If you encounter any problems
> during installation or any unexpected behavior in runtime, please let me know.

To run the synthesizer you will need to install
[portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html). On
Fedora it is

```
sudo dnf -y install portaudio
```

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

Type `synth` and you should get an output like this

```
command line synthesizer

documentation and usage: https://github.com/iljarotar/synth

Usage:
  synth [flags]

Flags:
  -d, --duration string      duration in seconds excluding fade-in and fade-out. a negative duration will cause the synth to play until stopped manually (default "-1")
      --fade-in string       fade-in in seconds (default "1")
      --fade-out string      fade-out in seconds (default "1")
  -f, --file string          path to your patch file
  -h, --help                 print help
  -o, --out string           if provided, a .wav file with the given name will be recorded
  -s, --sample-rate string   sample rate (default "44100")
```

To listen to an example run

```
git clone git@github.com:iljarotar/synth.git
cd synth
synth -f examples/a-major.yaml
```

More examples can be found [here](https://github.com/iljarotar/synth-patches).

> Note: If you want to record the output, you must specify a non-negative
> duration. Otherwise you will get am empty .wav file.

## Writing a patch file

### Data types

| Synth       |                   |                                                                                                           |
| ----------- | ----------------- | --------------------------------------------------------------------------------------------------------- |
| **Field**   | **Type**          | **Description**                                                                                           |
| vol         | Float             | main volume in range [0,1]                                                                                |
| out         | String [0..*]     | names of all oscillators, noise generators and custom signals, whose outputs will be sent to the speakers |
| oscillators | Oscillator [0..*] | all oscillators                                                                                           |
| noise       | Noise [0..*]      | all noise generators                                                                                      |

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

| Filter     |          |                                                            |
| ---------- | -------- | ---------------------------------------------------------- |
| **Field**  | **Type** | **Description**                                            |
| order      | Integer  | order of the FIR filter                                    |
| lowcutoff  | Float    | cutoff frequency of the highpass filter in range [0,20000] |
| highcutoff | Float    | cutoff frequency of the lowpass filter in range [0,20000]  |

If both lowcutoff and highcutoff are 0, the filter is disabled. If lowcutoff is
0, the filter is a lowpass filter transitioning at the highcutoff frequency. If
highcutoff is 0, the filter is a highpass filter transitioning at the lowcutoff
frequency.

A higher order improves the filter's precision, but it also makes it more
expensive in terms of computation. If the sound becomes glitchy, decreasing the
filter order might be necessary.

| Custom    |              |                                           |
| --------- | ------------ | ----------------------------------------- |
| **Field** | **Type**     | **Description**                           |
| name      | String       | should be unique in the scope of the file |
| amp       | Param        | amplitude in range [0,1]                  |
| pan       | Param        | stereo balance in range [-1,1]            |
| freq      | Param        | periods per second [0,20000]              |
| data      | Float [0..*] | custom values                             |

| Param     |               |                                                    |
| --------- | ------------- | -------------------------------------------------- |
| **Field** | **Type**      | **Description**                                    |
| val       | Float         | initial value of the respective parameter          |
| mod       | String [0..*] | names of modulating oscillators and custom signals |
| modamp    | Float         | amplitude of the modulation in range [0,1]         |

### Structure of a patch file

```yaml
vol:
out:
noise:
  - name:
    amp:
      val:
      mod:
      modamp:
    pan:
      val:
      mod:
      modamp:
      
oscillators:
  - name:
    type:
    freq:
      val:
      mod:
      modamp:
    amp:
      val:
      mod:
      modamp:
    pan:
      val:
      mod:
      modamp:
    phase:

custom:
  - name:
    freq:
      val:
      mod:
      modamp:
    amp:
      val:
      mod:
      modamp:
    pan:
      val:
      mod:
      modamp:
    data: []
```

Most of the fields are optional. A simple 440hz sine wave looks like this:

```yaml
vol: 1
out: [osc]
oscillators:
  - name: osc
    type: Sine
    freq: {val: 440}
    amp: {val: 1}
```
