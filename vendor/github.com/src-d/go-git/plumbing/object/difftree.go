package object

import (
	"bytes"
	"os"

	"srcd.works/go-git.v4/utils/merkletrie"
	"srcd.works/go-git.v4/utils/merkletrie/noder"
)

// DiffTree compares the content and mode of the blobs found via two
// tree objects.
func DiffTree(a, b *Tree) (Changes, error) {
	from := newTreeNoder(a)
	to := newTreeNoder(b)

	merkletrieChanges, err := merkletrie.DiffTree(from, to, hashEqual)
	if err != nil {
		return nil, err
	}

	return newChanges(merkletrieChanges)
}

// check if the hash of the contents is different, if not, check if
// the permissions are different (but taking into account deprecated
// file modes).  On a treenoder, the hash of the contents is codified
// in the first 20 bytes of the data returned by Hash() and the last
// 4 bytes is the mode.
func hashEqual(a, b noder.Hasher) bool {
	hashA, hashB := a.Hash(), b.Hash()
	contentsA, contentsB := hashA[:20], hashB[:20]

	sameContents := bytes.Equal(contentsA, contentsB)
	if !sameContents {
		return false
	}

	modeA, modeB := hashA[20:], hashB[20:]

	return equivalentMode(modeA, modeB)
}

func equivalentMode(a, b []byte) bool {
	if isFilish(a) && isFilish(b) {
		return true
	}
	return bytes.Equal(a, b)
}

var (
	file           = modeToBytes(FileMode)
	fileDeprecated = modeToBytes(FileModeDeprecated)
	// remove this by fixing plumbing.Object mode ASAP
	fileGoGit = modeToBytes(os.FileMode(0644))
)

func isFilish(b []byte) bool {
	return bytes.Equal(b, file) ||
		bytes.Equal(b, fileDeprecated) ||
		bytes.Equal(b, fileGoGit)
}
