package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ScrollView struct {
	g *gocui.Gui
}

func (sv *ScrollView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("log-graph", 0, 0, maxX*2/3, maxY); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func (sv *ScrollView) ScrollToGraph(graph string) {
	sv.g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("log-graph")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, graph)
		return nil
	})
}
