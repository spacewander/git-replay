package main

import (
	"github.com/jroimartin/gocui"
)

func DrawUI(gitLogs string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		errorLogger.Panicln(err)
	}
	defer g.Close()

	scrollView := &ScrollView{
		graph: gitLogs,
	}
	storyView := &StoryView{}
	commitView := &CommitView{}
	g.SetManager(scrollView, storyView, commitView)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		errorLogger.Panicln(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		errorLogger.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		errorLogger.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
