package main

import (
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

var (
	done = make(chan struct{})

	scrollView *ScrollView
	storyView  *StoryView
	commitView *CommitView

	gitLogs []string
)

func DrawUI(_gitLogs []string) {
	gitLogs = _gitLogs
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		errorLogger.Panicln(err)
	}
	defer g.Close()

	scrollView = &ScrollView{
		g: g,
	}
	storyView = &StoryView{}
	commitView = &CommitView{}
	g.SetManager(scrollView, storyView, commitView)

	go tick(g)
	defer close(done)

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

func tick(g *gocui.Gui) {
	_, windowSize := g.Size()
	topPadding := windowSize - 1
	lastLogIdx := len(gitLogs) - 1

	scrollTicker := time.NewTicker(250 * time.Millisecond)
	defer scrollTicker.Stop()
	for {
		select {
		case <-done:
			return
		case <-scrollTicker.C:
			eof := scrollLog(lastLogIdx, topPadding, windowSize)
			if eof {
				return
			}
			lastLogIdx--
			if topPadding > 0 {
				topPadding--
			}
		}
	}
}

func scrollLog(lastLogIdx, topPadding, windowSize int) (eof bool) {
	var graph string
	if topPadding > 0 {
		if lastLogIdx < 0 {
			return true
		}
		graph = strings.Repeat("\n", topPadding) +
			strings.Join(gitLogs[lastLogIdx:], "\n")
	} else {
		if lastLogIdx < windowSize {
			return true
		}
		graph = strings.Join(gitLogs[lastLogIdx-windowSize:lastLogIdx], "\n")
	}
	scrollView.ScrollToGraph(graph)
	return false
}
