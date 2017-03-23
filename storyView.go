package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type StoryView struct {
	g *gocui.Gui
}

func (sv *StoryView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("story", maxX*2/3, 0, maxX, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "waiting for your story...")
	}
	return nil
}

func (sv *StoryView) Show(story string) {
	sv.g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("story")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, story)
		return nil
	})
}
