package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"regexp"
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
	// The commit hash is used for searching commit message only.
	// So strip it when display.
	re := regexp.MustCompile(`[0-9a-f]{40}`)
	graph = re.ReplaceAllLiteralString(graph, "")
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
