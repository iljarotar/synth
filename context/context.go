package context

import "github.com/gordonklaus/portaudio"

type Context struct {
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Init() {
	portaudio.Initialize()
}

func (c *Context) Terminate() {
	portaudio.Terminate()
}
