package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	info *CommitInfo
)

func TestMain(m *testing.M) {
	hash := "a592ab93c727bf21bbeeccb19b8a1362c6a73c96"
	err := InitRepo(".")
	if err != nil {
		panic(err)
	}
	commit, err := SearchCommit(hash)
	if err != nil {
		panic(err)
	}
	info = ExtractDataFromCommit(commit)

	m.Run()
}

func TestSubDateToIsoFormat(t *testing.T) {
	commit := `commit 35890e3fb07bc74976ad6c6f58480bbdcc708b7c
		Author: luozexuan <luozexuan@b.360.cn>
		Date:   Wed Feb 22 10:37:38 2017 +0800

		Date:   Wed Feb 22 10:37:38 2017 +0800
		index on release`
	new_commit := subDateToIsoFormat(commit)
	assert.Equal(t, `commit 35890e3fb07bc74976ad6c6f58480bbdcc708b7c
		Author: luozexuan <luozexuan@b.360.cn>
		Date:   2017-02-22 10:37:38 +0800

		Date:   Wed Feb 22 10:37:38 2017 +0800
		index on release`, new_commit)
}

func TestParseCommitMessage(t *testing.T) {
	info := &CommitInfo{
		message: `Merge commit '3550dbad2bfe220f2d01b7a95bf34ac0af9c5829' into feature/custom_chinatax

		Conflict:
		xxx`,
	}
	info.parseCommitMessage()
	assert.Equal(t, "Merge commit '3550dbad2bfe220f2d01b7a95bf34ac0af9c5829' into feature/custom_chinatax",
		info.title)
	assert.Equal(t, "3550dbad2bfe220f2d01b7a95bf34ac0af9c5829", info.merge_from)
	assert.Equal(t, "feature/custom_chinatax", info.merge_to)

	info = &CommitInfo{
		message: `Merge branch 'develop' into hotfix/6.0.0.2500`,
	}
	info.parseCommitMessage()
	assert.Equal(t, "Merge branch 'develop' into hotfix/6.0.0.2500", info.title)
	assert.Equal(t, "develop", info.merge_from)
	assert.Equal(t, "hotfix/6.0.0.2500", info.merge_to)

	info = &CommitInfo{
		message: `terminal ui with git log --graph

		blah blah blah`,
	}
	info.parseCommitMessage()
	assert.Equal(t, "terminal ui with git log --graph", info.title)
	assert.Equal(t, "", info.merge_from)
	assert.Equal(t, "", info.merge_to)
}

func TestExtractDataFromCommit(t *testing.T) {
	assert.Equal(t, "add gitignore and LICENSE\n", info.message)
	assert.Equal(t, "add gitignore and LICENSE", info.title)

	assert.Equal(t, "spacewander", info.name)
	assert.Equal(t, "spacewanderlzx@gmail.com", info.email)
	assert.Equal(t, "2017-02-24 17:13:30 +0800", info.date)
	assert.Equal(t, "spacewander", info.author_name)
	assert.Equal(t, "spacewanderlzx@gmail.com", info.author_email)
	assert.Equal(t, "2017-02-24 17:13:30 +0800", info.author_date)
	assert.Equal(t, "spacewander", info.committer_name)
	assert.Equal(t, "spacewanderlzx@gmail.com", info.committer_email)
	assert.Equal(t, "2017-02-25 14:47:53 +0800", info.committer_date)

	tag := info.tags[0]
	assert.Equal(t, "0.01", tag.name)
	assert.Equal(t, "tag for test\n", tag.message)
	assert.Equal(t, "2017-02-26 17:12:06 +0800", tag.date)
	assert.Equal(t, "spacewander", tag.tagger_name)
	assert.Equal(t, "spacewanderlzx@gmail.com", tag.tagger_email)
}

func TestExplainCommitInLuaScript(t *testing.T) {
	case_prefix := "fixture/dump_commit"
	defer func() {
		os.Remove(case_prefix + ".actual.txt")
	}()

	PlayWithCommitInfo(case_prefix+".lua", info)
	actual, err := ioutil.ReadFile(case_prefix + ".actual.txt")
	if err != nil {
		assert.Fail(t, err.Error())
	} else {
		expect, _ := ioutil.ReadFile(case_prefix + ".expect.txt")
		assert.Equal(t, strings.TrimRight(string(expect), "\n"), string(actual))
	}
}

func TestDebugLogging(t *testing.T) {
	defer func() {
		os.Remove(debugLogFilename)
	}()
	debugLogger = true
	debugLogger.Print("log info")
	log_info, err := ioutil.ReadFile(debugLogFilename)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.True(t, strings.HasSuffix(string(log_info), "log info\n"))
}

func TestFormatStatusLine(t *testing.T) {
	// required length: 2*(5 + 5 + 3) = 26
	status := map[string]string{
		"name1": "value",
		"name2": "value",
	}
	width := 30
	actual := formatStatusLine(width, status)
	assert.Equal(t, "  name1: value  name2: value", actual)
}

func TestFormatStatusLine_LongerThanRequired(t *testing.T) {
	status := map[string]string{
		"name1": "value",
		"name2": "value",
	}
	width := 26
	actual := formatStatusLine(width, status)
	assert.Equal(t, " name1: value name2: value", actual)
}

func TestFormatStatusLine_EmptyStatus(t *testing.T) {
	status := map[string]string{}
	width := 40
	actual := formatStatusLine(width, status)
	assert.Equal(t, "", actual)
}

func TestFormatStatusLine_NameOnlyStatus(t *testing.T) {
	status := map[string]string{
		"name1": "value",
		"name2": "",
		"name3": "value",
	}
	width := 36
	actual := formatStatusLine(width, status)
	assert.Equal(t, "  name1: value  name2  name3: value", actual)
}
