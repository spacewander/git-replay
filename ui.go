package main

import (
	"github.com/jroimartin/gocui"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SpeedControlSignal int

const (
	SlowDown SpeedControlSignal = iota
	Paused
	SpeedUp
)

var (
	done             = make(chan struct{})
	speedControlChan = make(chan SpeedControlSignal)

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
		Notice: `press 'q' to quit`,
	}
	statusLineView = &StatusLineView{
		BaseView: BaseView{
			id: "status line",
			g:  g,
		},
	}
	statusLineView.InitStatus()
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
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, togglePaused); err != nil {
		errorLogger.Panicln(err)
	}
	if err := g.SetKeybinding("", '>', gocui.ModNone, speedUp); err != nil {
		errorLogger.Panicln(err)
	}
	if err := g.SetKeybinding("", '<', gocui.ModNone, slowDown); err != nil {
		errorLogger.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		errorLogger.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func togglePaused(g *gocui.Gui, v *gocui.View) error {
	debugLogger.Println("toggle paused")
	speedControlChan <- Paused
	return nil
}

func speedUp(g *gocui.Gui, v *gocui.View) error {
	debugLogger.Println("speed up")
	speedControlChan <- SpeedUp
	return nil
}

func slowDown(g *gocui.Gui, v *gocui.View) error {
	debugLogger.Println("slow down")
	speedControlChan <- SlowDown
	return nil
}

func tick(g *gocui.Gui) {
	_, windowSize := g.Size()
	_, topPadding := scrollView.Size()
	lastLogIdx := len(gitLogs) - 1

	re := regexp.MustCompile(`[0-9a-f]{40}`)

	// 4 lines per second
	scrollTicker := time.NewTicker(250 * time.Millisecond)
	defer scrollTicker.Stop()

	speed := 1
	paused := false
	end := false
	statusLineView.UpdateStatus("speed", "1")
	for {
		select {
		case <-done:
			return
		case <-scrollTicker.C:
			select {
			case speedSignal := <-speedControlChan:
				switch speedSignal {
				case Paused:
					if paused {
						statusLineView.RemoveStatus("paused")
					} else {
						statusLineView.UpdateStatus("paused", "")
					}
					paused = !paused
				case SpeedUp:
					if speed <= 8 {
						speed *= 2
						statusLineView.UpdateStatus("speed", strconv.Itoa(speed))
					}
				case SlowDown:
					if speed > 1 {
						speed /= 2
						statusLineView.UpdateStatus("speed", strconv.Itoa(speed))
					}
				}
				debugLogger.Printf("Current speed %v, is paused %v", speed, paused)

			default:
			}

			if paused || end {
				continue
			}

			for i := 0; i < speed; i++ {
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
					end = true
					// halt in the rest time
					break
				}
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
