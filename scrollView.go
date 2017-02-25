package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ScrollView struct {
	graph string
}

func (sv *ScrollView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("git log --graph", 0, 0, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, sv.graph)
	}
	return nil
}
