package main

import (
	"github.com/jroimartin/gocui"
	"regexp"
	"strings"
	"time"
)

var (
	done = make(chan struct{})

	scrollView     *ScrollView
	storyView      *StoryView
	commitView     *CommitView
	noticeView     *NoticeView
	statusLineView *StatusLineView

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
		BaseView: BaseView{
			id: "scroll",
			g:  g,
		},
	}
	storyView = &StoryView{
		BaseView: BaseView{
			id: "story",
			g:  g,
		},
	}
	commitView = &CommitView{
		BaseView: BaseView{
			id: "commit",
			g:  g,
		},
	}
	noticeView = &NoticeView{
		BaseView: BaseView{
			id: "notice",
			g:  g,
		},
		notice: `press 'q' to quit`,
	}
	statusLineView = &StatusLineView{
		BaseView: BaseView{
			id: "status line",
			g:  g,
		},
		status: map[string]string{
			"speed": "1",
		},
	}
	g.SetManager(
		scrollView,
		storyView,
		commitView,
		noticeView,
		statusLineView,
	)

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

	re := regexp.MustCompile(`[0-9a-f]{40}`)

	// 4 lines per second
	scrollTicker := time.NewTicker(250 * time.Millisecond)
	defer scrollTicker.Stop()
	for {
		select {
		case <-done:
			return
		case <-scrollTicker.C:
			hash := re.FindString(gitLogs[lastLogIdx])
			if hash != "" {
				if commit, err := SearchCommit(hash); err == nil {
					debugLogger.Println("commit: ", commit.Hash.String())
					commitView.DisplayCommit(commit)
					if scriptName != "" {
						commitInfo := ExtractDataFromCommit(commit)
						if err := PlayWithCommitInfo(scriptName, commitInfo); err != nil {
							errorLogger.Panicln(err)
						}
					}
				}
			}
			scrollLog(lastLogIdx, topPadding, windowSize)
			if topPadding > 0 {
				topPadding--
			}
			lastLogIdx--
			if lastLogIdx < 0 {
				return
			}
		}
	}
}

func scrollLog(lastLogIdx, topPadding, windowSize int) {
	var graph string
	if topPadding > 0 {
		if lastLogIdx < 0 {
			return
		}
		graph = strings.Repeat("\n", topPadding) +
			strings.Join(gitLogs[lastLogIdx:], "\n")
	} else {
		if lastLogIdx < windowSize {
			return
		}
		graph = strings.Join(gitLogs[lastLogIdx-windowSize:lastLogIdx], "\n")
	}
	scrollView.ScrollToGraph(graph)
}
