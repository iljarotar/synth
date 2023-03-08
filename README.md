# Synth

This is a simple modular-like command line synthesizer written in [golang](https://go.dev/).

## Installation

Note: I have tested the synth only on a Fedora x86_64 machine. If you encounter any problems during installation or any unexpected behavior in runtime, please let me know.

To run the synthesizer you will need to install [portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html). On Fedora it is simply

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

The synth is not meant to be played, but to be programmed by providing a patch file. A patch is a `yaml` file, that tells the synth, which oscillators, filters, etc. should be created and how they should be connected. When you tell the synth to load a file, if will start playing immediately. During playback a hot reload is possible, so if you change and save the patch file, it will be applied instantly. But the transition will be audible, that's why it isn't meant to be "played".

### Command line interface

Type `synth` and you should get an output like this

```
command line synthesizer

documentation and usage: https://github.com/iljarotar/synth

Usage:
  synth [flags]

Flags:
  -d, --duration string      duration in seconds excluding fade-in and fade-out (default "0")
      --fade-in string       fade-in in seconds (default "1")
      --fade-out string      fade-out in seconds (default "1")
  -f, --file string          path to your patch file
  -h, --help                 print help
  -o, --out string           if provided, a recording will be written to the given file
  -s, --sample-rate string   sample rate (default "44100")
```

To listen to an example run

```
git clone git@github.com:iljarotar/synth.git
cd synth
synth -f examples/a-major.yaml
```

More examples can be found [here](https://github.com/iljarotar/synth-patches).

## Writing a patch file

### Data types

| Synth       |                   |                                                            |
| ----------- | ----------------- | ---------------------------------------------------------- |
| Field       | Type              | Description                                                |
| vol         | Float             | main volume in range [0,1]                                 |
| out         | String [0..*]     | names of the oscillators that will be sent to the speakers |
| oscillators | Oscillator [0..*] | all oscillators                                            |
| filters     | Filter [0..*]     | all filters                                                |

| Oscillator |                |                                           |
| ---------- | -------------- | ----------------------------------------- |
| Field      | Type           | Description                               |
| name       | String         | should be unique in the scope of the file |
| type       | OscillatorType | wave form or noise                        |
| freq       | Float [0..*]   | frequencies in range [0,20000]            |
| amp        | Param          | amplitude in range [0,1]                  |
| phase      | Param          | phase in range [-1,1]                     |
| filters    | String [0..*]  | names of filters to be applied            |
| pan        | Param          | stereo balance in range [-1,1]            |

| OscillatorType   |
| ---------------- |
| Sine             |
| Triangle         |
| Square           |
| Sawtooth         |
| InvertedSawtooth |
| Noise            |

Note: Noise will not be affected by filters or frequency

| Filter |        |                                            |
| ------ | ------ | ------------------------------------------ |
| Field  | Type   | Description                                |
| name   | String | should be unique in the scope of the file  |
| low    | Param  | low frequency cutoff                       |
| high   | Param  | high frequency cutoff                      |
| vol    | Param  | volume of unfiltered signal                |
| ramp   | Float  | length of the linear ramp from cutoff to 0 |

Note: All filters are bandpass filters. To create a highpass or lowpass filter just place one of the cutoffs outside of the audible range.

| Param  |               |                                            |
| ------ | ------------- | ------------------------------------------ |
| Field  | Type          | Description                                |
| val    | Float         | initial value of the respective param      |
| mod    | String [0..*] | names of modulator oscillators             |
| modamp | Float         | amplitude of the modulation in range [0,1] |

### Structure of a patch file

```yaml
vol: 1
out:
oscillators:
  - name:
    type:
    freq:
    amp:
      val:
      mod:
      modamp:
    pan:
      val:
      mod:
      modamp:
    phase:
      val:
      mod:
      modamp:
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

Most of the fields are optional. A simple 440hz sine wave would look like this:

```yaml
vol: 1
out: [osc]
oscillators:
  - name: osc
    type: Sine
    freq: [440]
    amp:
      val: 1
```
