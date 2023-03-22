# Synth

This is a simple modular-like command line synthesizer written in
[golang](https://go.dev/).

## Installation

Note: I have tested the synth only on a Fedora x86_64 machine. If you encounter
any problems during installation or any unexpected behavior in runtime, please
let me know.

To run the synthesizer you will need to install
[portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html). On
Fedora it is simply

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

The synth is not meant to be played, but to be programmed by providing a patch
file. A patch is a `yaml` file, that tells the synth, which oscillators,
filters, etc. should be created and how they should be connected. When you tell
the synth to load a file, if will start playing immediately. During playback a
hot reload is possible, so if you change and save the patch file, it will be
applied instantly. But the transition will be audible, that's why it isn't meant
to be "played".

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

Note: If you want to record the output, you must specify a non-negative
duration. Otherwise you will get am empty .wav file.

## Writing a patch file

### Data types

| Synth       |                   |                                                                                           |
| ----------- | ----------------- | ----------------------------------------------------------------------------------------- |
| **Field**   | **Type**          | **Description**                                                                           |
| vol         | Float             | main volume in range [0,1]                                                                |
| out         | String [0..*]     | names of all oscillators and noise generators, whose outputs will be sent to the speakers |
| oscillators | Oscillator [0..*] | all oscillators                                                                           |
| filters     | Filter [0..*]     | all filters                                                                               |
| noise       | Noise [0..*]      | all noise generators                                                                      |

| Oscillator |                |                                           |
| ---------- | -------------- | ----------------------------------------- |
| **Field**  | **Type**       | **Description**                           |
| name       | String         | should be unique in the scope of the file |
| type       | OscillatorType | wave form                                 |
| freq       | Param          | frequency in range [0,20000]              |
| amp        | Param          | amplitude in range [0,1]                  |
| phase      | Float          | phase in range [-1,1]                     |
| filters    | String [0..*]  | names of filters to be applied            |
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

| Filter    |          |                                            |
| --------- | -------- | ------------------------------------------ |
| **Field** | **Type** | **Description**                            |
| name      | String   | should be unique in the scope of the file  |
| low       | Param    | low frequency cutoff                       |
| high      | Param    | high frequency cutoff                      |
| vol       | Param    | volume of unfiltered signal                |
| ramp      | Float    | length of the linear ramp from cutoff to 0 |

Note: All filters are bandpass filters. To create a highpass or a lowpass filter
just place one of the cutoffs outside of the audible range.

| Param     |               |                                            |
| --------- | ------------- | ------------------------------------------ |
| **Field** | **Type**      | **Description**                            |
| val       | Float         | initial value of the respective parameter  |
| mod       | String [0..*] | names of modulating oscillators            |
| modamp    | Float         | amplitude of the modulation in range [0,1] |

### Structure of a patch file

```yaml
vol:
out:
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
    filters:

filters:
  - name:
    low:
      val:
      mod:
      modamp:
    high:
      val:
      mod:
      modamp:
    vol:
      val:
      mod:
      modamp:
    ramp:
```

Most of the fields are optional. A simple 440hz sine wave looks like this:

```yaml
vol: 1
out: [osc]
oscillators:
  - name: osc
    type: Sine
    freq: 
      val: 440
    amp:
      val: 1
```
