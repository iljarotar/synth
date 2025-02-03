package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/term"
)

type UI struct {
	logger               *Logger
	quit                 chan bool
	input                chan string
	autoStop             chan bool
	file                 string
	logs                 []string
	time                 string
	duration             float64
	showOverdriveWarning bool
	closing              *bool
	interrupt            chan bool
}

func NewUI(logger *Logger, file string, quit chan bool, autoStop chan bool, duration float64, closing *bool, interrupt chan bool) *UI {
	return &UI{
		logger:    logger,
		quit:      quit,
		autoStop:  autoStop,
		input:     make(chan string),
		file:      file,
		time:      "00:00:00",
		duration:  duration,
		closing:   closing,
		interrupt: interrupt,
	}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func LineBreaks(number int) {
	for i := 0; i < number; i++ {
		fmt.Print("\r\n")
	}
}

func (ui *UI) Enter() {
	go ui.read()
	ui.resetScreen()

	logChan := make(chan string)
	ui.logger.SubscribeToLogs(logChan)

	timeChan := make(chan string)
	ui.logger.SubscribeToTime(timeChan)

	stateChan := make(chan State)
	ui.logger.SubscribeToState(stateChan)

	for {
		select {
		case input := <-ui.input:
			if input == "q" {
				*ui.closing = true
				ui.resetScreen()
				ui.quit <- true
			} else {
				ui.resetScreen()
			}
		case time := <-timeChan:
			ui.time = time
			ui.updateTime()
		case log := <-logChan:
			ui.appendLog(log)
			ui.resetScreen()
		case state := <-stateChan:
			ui.showOverdriveWarning = state.overdriveWarning
			ui.resetScreen()
		case <-ui.autoStop:
			*ui.closing = true
			ui.quit <- true
		}
	}
}

func (ui *UI) read() {
	state, err := term.MakeRaw(0)
	if err != nil {
		ui.logger.Error(fmt.Sprintf("failed to read input %v", err))
	}
	defer func() {
		if err := term.Restore(0, state); err != nil {
			ui.logger.Error(fmt.Sprintf("failed to restore terminal %v", err))
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			ui.logger.Error(fmt.Sprintf("failed to read input %v", err))
		}
		if r == rune(3) {
			ui.interrupt <- true
		}
		ui.input <- string(r)
	}
}

func (ui *UI) resetScreen() {
	Clear()
	LineBreaks(1)
	fmt.Printf("%s %s", colored("Synth playing", colorBlueStrong), ui.file)
	LineBreaks(2)

	for _, log := range ui.logs {
		fmt.Print(log + "\r\n")
	}
	if len(ui.logs) > 0 {
		LineBreaks(1)
	}
	if ui.showOverdriveWarning {
		fmt.Printf("%s", colored("[WARNING] Volume exceeded 100%%", colorOrangeStorng))
		LineBreaks(2)
	}
	fmt.Printf("%s", ui.time)
	if ui.duration >= 0 {
		fmt.Printf(" - automatically stopping after %fs", ui.duration)
	}
	LineBreaks(1)
	fmt.Printf("%s ", colored("Type 'q' to quit:", colorBlueStrong))
}

func (ui *UI) updateTime() {
	// using ANSI escape sequences:
	// \0337 to save current cursor location
	// \033[1A to move cursor up one line
	// \r to move cursor to beginning of line
	// \0338 to restore original cursor location
	fmt.Printf("\0337\033[1A\r%s\0338", ui.time)
}

func (ui *UI) appendLog(log string) {
	logs := ui.logs
	logs = append(logs, log)
	if len(logs) > 10 {
		logs = logs[1:]
	}
	ui.logs = logs
}
