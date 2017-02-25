package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type CommitView struct {
}

func (cv *CommitView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("commit", maxX*2/3, maxY/2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "commitView")
	}
	return nil
}
