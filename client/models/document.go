package models

import (
	"encoding/json"

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

type Document struct {
	BlockKeys []uuid.UUID
}

const DOCUMENT_BLOCK_SIZE = 32

func SplitDocumentToBlocks(document []byte) (blocks Blocks) {
	blocks = make(Blocks, 0)
	remainingDocument := document
	for len(remainingDocument) > 0 {
		end := DOCUMENT_BLOCK_SIZE
		if len(remainingDocument) < end {
			end = len(remainingDocument)
		}

		blocks = append(blocks, remainingDocument[:end])
		remainingDocument = remainingDocument[end:]
	}
	return
}

func LoadDocument(fileMapping FileMapping) (document Document, err error) {

	storageKey, err := fileMapping.LoadDocumentKey()
	if err != nil {
		return
	}

	documentJSON, ok := userlib.DatastoreGet(storageKey)
	if !ok {
		err = errors.New(strings.ToTitle("file not found"))
		return
	}

	err = json.Unmarshal(documentJSON, &document)

	if err != nil {
		return
	}

	return
}

func (document Document) Store(fileMapping FileMapping) (err error) {

	storageKey, err := fileMapping.LoadDocumentKey()
	if err != nil {
		return err
	}

	documentBytes, err := json.Marshal(document)
	if err != nil {
		return err
	}
	userlib.DatastoreSet(storageKey, documentBytes)
	return
}
