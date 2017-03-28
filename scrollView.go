package main

import (
	"github.com/jroimartin/gocui"
	"regexp"
)

type ScrollView struct {
	BaseView
}

func (sv *ScrollView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView(sv.id, 0, 0, maxX*2/3, maxY-4); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func (sv *ScrollView) ScrollToGraph(graph string) {
	// The commit hash is used for searching commit message only.
	// So strip it when display.
	re := regexp.MustCompile(`[0-9a-f]{40}`)
	graph = re.ReplaceAllLiteralString(graph, "")
	sv.BaseView.Show(graph)
}
