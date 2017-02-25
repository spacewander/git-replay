package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type StoryView struct {
}

func (sv *StoryView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("story", maxX*2/3, 0, maxX, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "storyView")
	}
	return nil
}
