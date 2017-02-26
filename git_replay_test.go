package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
