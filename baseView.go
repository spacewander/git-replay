package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mattn/go-runewidth"
)

type BaseView struct {
	id string
	g  *gocui.Gui
}

func (bv *BaseView) Show(message string) {
	bv.g.Execute(func(g *gocui.Gui) error {
		v, err := g.View(bv.id)
		if err != nil {
			return err
		}
		v.Clear()
		width, _ := bv.size()
		fmt.Fprintln(v, runewidth.Wrap(message, width))
		return nil
	})
}

func (bv *BaseView) size() (x, y int) {
	view, err := bv.g.View(bv.id)
	if err != nil {
		return 0, 0
	}
	return view.Size()
}
