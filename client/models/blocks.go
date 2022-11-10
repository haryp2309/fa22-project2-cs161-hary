package models

import (
	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

type Blocks [][]byte

func (blocks Blocks) MergeToBlob() (document []byte) {
	document = make([]byte, 0)
	for _, block := range blocks {
		document = append(document, block...)
	}
	return
}

func LoadBlocks(blockKeys []uuid.UUID) (blocks Blocks, err error) {

	blocks = make(Blocks, 0)

	for _, blockKey := range blockKeys {
		block, ok := userlib.DatastoreGet(blockKey)

		if !ok {
			err = errors.New(strings.ToTitle("block not found"))
			return
		}
		blocks = append(blocks, block)
	}

	return

}

func (blocks Blocks) Store() (blockKeys []uuid.UUID, err error) {
	blockKeys = make([]uuid.UUID, len(blocks))
	for i, block := range blocks {
		blockKeys[i] = uuid.New()
		userlib.DatastoreSet(blockKeys[i], block)
	}
	return
}
