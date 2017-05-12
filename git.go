package main

import (
	"regexp"
	"strings"

	go_git "github.com/src-d/go-git"
	"srcd.works/go-git.v4/plumbing"
	"srcd.works/go-git.v4/plumbing/object"
)

var (
	repo *go_git.Repository

	commitToTag = make(map[plumbing.Hash][]plumbing.Hash)
)

// A plain data object contains commit info,  which we will pass to Lua script later.
// Note that all empty string would be nil in Lua side.
type CommitInfo struct {
	author_name  string
	author_email string
	// yyyy-MM-dd mm:hh:ss +tz
	author_date string

	committer_name  string
	committer_email string
	// yyyy-MM-dd mm:hh:ss +tz
	committer_date string

	hash string

	// alias to author_xxx
	name  string
	email string
	date  string

	message string
	// title is the first line of message (without line break)
	title string

	// Empty if it is not a merged commit.
	merge_from string
	merge_to   string

	// Tags attached with given commit.
	tags []*TagInfo
}

type TagInfo struct {
	name string
	date string
	// message attached with `git tag -m`, notice that there is a trailing line break.
	message string
	// tagger is the guy who creates the tag
	tagger_name  string
	tagger_email string
}

func InitRepo(path string) (err error) {
	debugLogger.Println("init repo ", path)
	repo, err = go_git.PlainOpen(path)
	if err != nil {
		return err
	}

	// Only annotated tags(tag created with git tag -a) are listed.
	// You can figure them out with `git for-each-ref "refs/tags"`
	iter, err := repo.Tags()
	if err != nil {
		return nil
	}
	for tag, err := iter.Next(); err == nil; tag, err = iter.Next() {
		debugLogger.Println("tag: ", tag.Name, " ", tag.Hash.String())
		if commit, err := tag.Commit(); err == nil {
			debugLogger.Println("commit attached by tag: ", commit.Hash.String())
			if prev_tags, ok := commitToTag[commit.Hash]; ok {
				commitToTag[commit.Hash] = append(prev_tags, tag.Hash)
			} else {
				commitToTag[commit.Hash] = []plumbing.Hash{tag.Hash}
			}
		}
	}
	return nil
}

func SearchCommit(hash string) (*object.Commit, error) {
	commit, err := repo.Commit(plumbing.NewHash(hash))
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func ExtractDataFromCommit(commit *object.Commit) *CommitInfo {
	formatStr := "2006-01-02 15:04:05 -0700"
	info := &CommitInfo{
		hash:            commit.Hash.String(),
		message:         commit.Message,
		name:            commit.Author.Name,
		email:           commit.Author.Email,
		date:            commit.Author.When.Format(formatStr),
		author_name:     commit.Author.Name,
		author_email:    commit.Author.Email,
		author_date:     commit.Author.When.Format(formatStr),
		committer_name:  commit.Committer.Name,
		committer_email: commit.Committer.Email,
		committer_date:  commit.Committer.When.Format(formatStr),
	}
	if tagHashes, ok := commitToTag[commit.Hash]; ok {
		info.tags = []*TagInfo{}
		for _, tagHash := range tagHashes {
			if tag, err := repo.Tag(tagHash); err == nil {
				debugLogger.Println("found tag: ", tag.Name, " ", tag.Hash.String())
				tagInfo := &TagInfo{
					name:         tag.Name,
					message:      tag.Message,
					date:         tag.Tagger.When.Format(formatStr),
					tagger_name:  tag.Tagger.Name,
					tagger_email: tag.Tagger.Email,
				}
				info.tags = append(info.tags, tagInfo)
			}
		}
	}
	info.parseCommitMessage()
	return info
}

func (info *CommitInfo) parseCommitMessage() {
	info.title = strings.SplitN(info.message, "\n", 2)[0]
	re := regexp.MustCompile(`Merge \w+ '([^']+)' into ([^\s]+)`)
	matchGroup := re.FindStringSubmatch(info.title)
	if matchGroup != nil {
		info.merge_from = matchGroup[1]
		info.merge_to = matchGroup[2]
		debugLogger.Println("merge from: ", info.merge_from, " to: ", info.merge_to)
	}
}
