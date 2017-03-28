package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"regexp"
	"srcd.works/go-git.v4/plumbing/object"
	"time"
)

type CommitView struct {
	BaseView
}

func (cv *CommitView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(cv.id, maxX*2/3, maxY/2, maxX, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "commitView")
	}
	return nil
}

func (cv *CommitView) DisplayCommit(commit *object.Commit) {
	commitMessage := subDateToIsoFormat(commit.String())
	cv.BaseView.Show(commitMessage)
}

// substitute 'Date:   Wed Feb 17 16:20:26 2016 +0800' to
// 'Date:   2016-02-17 16:20:26 +0800'
func subDateToIsoFormat(commitMessage string) string {
	re := regexp.MustCompile(`Date:   [\w\s:]+ [-|+]\d{4}`)
	found := false
	// One-shot replacer
	return re.ReplaceAllStringFunc(commitMessage, func(date string) string {
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
