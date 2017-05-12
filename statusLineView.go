package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mattn/go-runewidth"
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
		slv.drawStatusLine()
	}
	return nil
}

func (slv *StatusLineView) InitStatus() {
	if slv.status == nil {
		slv.status = map[string]string{}
	}
}

func (slv *StatusLineView) UpdateStatus(name string, value string) {
	slv.status[name] = value
	slv.drawStatusLine()
}

func (slv *StatusLineView) RemoveStatus(name string) {
	delete(slv.status, name)
	slv.drawStatusLine()
}

func (slv *StatusLineView) drawStatusLine() {
	width, _ := slv.g.Size()
	slv.BaseView.Show(formatStatusLine(width, slv.status))
}

func formatStatusLine(width int, status map[string]string) string {
	statusLine := []string{}
	requiredSpace := 0
	for name, value := range status {
		valueWidth := runewidth.StringWidth(value)
		// if not value, show name only
		if valueWidth == 0 {
			// 1 == len(defaultSeparatorSize)
			requiredSpace += runewidth.StringWidth(name) + 1
		} else {
			// 3 == len(defaultSeparatorSize) + len(": ")
			requiredSpace += runewidth.StringWidth(name) + valueWidth + 3
		}
	}
	separatorSize := 1
	if requiredSpace < width {
		separatorSize += (width - requiredSpace) / (len(status) + 1)
	}
	for name, value := range status {
		valueWidth := runewidth.StringWidth(value)
		if valueWidth == 0 {
			statusLine = append(statusLine, fmt.Sprintf("%s%s",
				strings.Repeat(" ", separatorSize), name))
		} else {
			statusLine = append(statusLine, fmt.Sprintf("%s%s: %s",
				strings.Repeat(" ", separatorSize), name, value))

		}
	}
	sort.Strings(statusLine)
	return strings.Join(statusLine, "")
}
