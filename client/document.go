package client

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	"strings"

	"errors"

	_ "strconv"
)

type Document struct {
	BlocksCount   int
	OwnerUsername string
}

const DEBUG_DOCUMENT = false

func LoadDocument(storageKey uuid.UUID) (document Document, err error) {

	documentJSON, ok, err := DatastoreGet(storageKey)
	if err != nil {
		return
	}
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

func (document Document) Store(storageKey uuid.UUID) (err error) {

	documentBytes, err := json.Marshal(document)
	if err != nil {
		return err
	}
	err = DatastoreSet(storageKey, documentBytes)
	if err != nil {
		return
	}

	return
}

func InitDocument(blocksCount int, ownerUsername string) (document Document) {
	document = Document{
		BlocksCount:   blocksCount,
		OwnerUsername: ownerUsername,
	}
	return
}

func (document Document) IsOwner(username string) (err error) {
	if DEBUG_DOCUMENT {
		userlib.DebugMsg("\nChecking if %s is the owner of document. The actual owner is %s.\n", username, document.OwnerUsername)
	}
	if username != document.OwnerUsername {
		err = errors.New(strings.ToTitle("user is not owner"))
	}

	return
}
