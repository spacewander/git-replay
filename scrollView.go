package main

import (
	"github.com/jroimartin/gocui"
	"regexp"
)

type ScrollView struct {
	BaseView
}

func (sv *ScrollView) Layout(g *gocui.Gui) error {
	ow, oh := 0, 0
	width, height := sv.Size()
	if _, err := g.SetView(sv.id, ow, oh,
		ow+width+1, oh+height+1); err != nil && err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

// Size() will be used before Layout(*gocui.Gui) setup the view.
// So we define here instead of depending on BaseView's size()
func (sv *ScrollView) Size() (width, height int) {
	maxX, maxY := sv.g.Size()
	return maxX*2/3 - 1, maxY - 5
}

func (sv *ScrollView) ScrollToGraph(graph string) {
	// The commit hash is used for searching commit message only.
	// So strip it when display.
	re := regexp.MustCompile(`[0-9a-f]{40}`)
	graph = re.ReplaceAllLiteralString(graph, "")
	sv.BaseView.Show(graph)
}
