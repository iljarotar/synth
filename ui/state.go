package ui

type state struct {
	Closed                  bool
	ShowingOverdriveWarning bool
	CurrentTime             int
}

var State = state{}
