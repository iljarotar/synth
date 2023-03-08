# Synth

This is a simple modular-like command line synthesizer written in [golang](https://go.dev/).

## Installation

<small>Note: I have tested the synth only on a Fedora x86_64 machine. If you encounter any problems during installation or any unexpected behavior in runtime, please let me know.</small>

To run the synthesizer you will need to install [portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html). On Fedora it is simply

```bash
sudo dnf -y install portaudio
```

To install the synthesizer you have two options.

### Install with Go

1. Install and setup [Go](https://go.dev/doc/install)
2. Clone the repository

```bash
git clone git@github.com:iljarotar/synth.git
```

3. Install

```bash
cd synth
go install
```

### Download the binary

<small>Note: Currently binaries are only available for linux amd64.</small>

1. Download the binary from the [releases](https://github.com/iljarotar/synth/releases) page
2. Make it executable

```bash
chmod +x synth_linux_amd64
mv synth_linux_amd64 <SOMEWHERE_INSIDE_YOUR_PATH>/synth
```

## Usage

## Writing a patch file
