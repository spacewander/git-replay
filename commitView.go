package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"regexp"
	"time"
)

type CommitView struct {
	g *gocui.Gui
}

func (cv *CommitView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("commit", maxX*2/3, maxY/2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "commitView")
	}
	return nil
}

func (cv *CommitView) Show(hash string) {
	cv.g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("commit")
		if err != nil {
			return err
		}
		v.Clear()
		commit, err := SearchCommit(hash)
		if err != nil {
			return err
		}
		commit = subDateToIsoFormat(commit)
		fmt.Fprintln(v, commit)
		return nil
	})
}

// substitute 'Date:   Wed Feb 17 16:20:26 2016 +0800' to
// 'Date:   2016-02-17 16:20:26 +0800'
func subDateToIsoFormat(commit string) string {
	re := regexp.MustCompile(`Date:   [\w\s:]+ [-|+]\d{4}`)
	found := false
	// One-shot replacer
	return re.ReplaceAllStringFunc(commit, func(date string) string {
		if found {
			return date
		}
		found = true
		return convertDateToIsoFormat(date)
	})
}

func convertDateToIsoFormat(date string) string {
	end := len(date)
	prefix := date[:8]
	timeFormat := date[8 : end-6]
	suffix := date[end-6:]
	dt, err := time.Parse(time.ANSIC, timeFormat)
	if err != nil {
		debugLogger.Printf("format %s to iso format failed: %s", timeFormat, err)
		return date
	}
	return prefix + dt.Format("2006-01-02 15:04:05") + suffix
}
