package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type NoticeView struct {
	BaseView
	Notice string
}

func (nv *NoticeView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// because of border, we need two line to display one line content
	if v, err := g.SetView(nv.id, 0, maxY-4, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, nv.Notice)
	}
	return nil
}
