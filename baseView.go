package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
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
		fmt.Fprintln(v, message)
		return nil
	})
}
