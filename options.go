package cronic

import "github.com/gdamore/tcell"

type Options struct {
	color tcell.Color
	path  string
}

func NewOptions() *Options {
	return &Options{
		color: tcell.GetColor("white"),
		path:  "/etc/cron.d",
	}
}

func (o *Options) SetColor(color string) *Options {
	o.color = tcell.GetColor(color)
	return o
}

func (o *Options) SetPath(path string) *Options {
	o.path = path
	return o
}
