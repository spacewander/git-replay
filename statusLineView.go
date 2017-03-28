package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mattn/go-runewidth"
	"sort"
	"strings"
)

type StatusLineView struct {
	BaseView
	status map[string]string
}

func (slv *StatusLineView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// because of border, we need two line to display one line content
	if _, err := g.SetView(slv.id, 0, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		slv.LayoutStatusLine()
	}
	return nil
}

func (slv *StatusLineView) LayoutStatusLine() {
	width, _ := slv.g.Size()
	slv.BaseView.Show(formatStatusLine(width, slv.status))
}

func formatStatusLine(width int, status map[string]string) string {
	statusLine := []string{}
	requiredSpace := 0
	for name, value := range status {
		// 3 == len(defaultSeparatorSize) + len(": ")
		requiredSpace += runewidth.StringWidth(name) + runewidth.StringWidth(value) + 3
	}
	separatorSize := 1
	if requiredSpace < width {
		separatorSize += (width - requiredSpace) / (len(status) + 1)
	}
	for name, value := range status {
		statusLine = append(statusLine, fmt.Sprintf("%s%s: %s",
			strings.Repeat(" ", separatorSize), name, value))
	}
	sort.Strings(statusLine)
	return strings.Join(statusLine, "")
}
