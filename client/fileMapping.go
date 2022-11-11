package client

import (
	"encoding/json"
	"fmt"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	"fa22-project2-cs161-hary/client/helpers"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

type FileMapping struct {
	Username string
	Filename string
}

func InitFileMapping(username string, filename string) (filemapping FileMapping) {
	filemapping = FileMapping{
		Username: username,
		Filename: filename,
	}
	return
}

const DEBUG_FILEMAPPING = false

func (fileMapping FileMapping) LoadDocumentKey() (documentKey uuid.UUID, err error) {

	key, err := helpers.GenerateDataStoreKey(fileMapping.Filename + fileMapping.Username)
	if err != nil {
		return
	}

	documentKeyBytes, ok := userlib.DatastoreGet(key)

	if !ok {
		err = errors.New("WRONG USERNAME OR PASSWORD")
		return
	}

	err = json.Unmarshal(documentKeyBytes, &documentKey)
	if err != nil {
		return
	}

	if DEBUG_FILEMAPPING {
		fmt.Printf("\n DEBUG: USER %s ACCESSED FILE %s, DOCUMENTKEY %s\n", fileMapping.Username, fileMapping.Filename, documentKey.String())
	}

	return
}

func (fileMapping FileMapping) StoreDocumentKey(documentKey uuid.UUID) (err error) {
	key, err := helpers.GenerateDataStoreKey(fileMapping.Filename + fileMapping.Username)
	if err != nil {
		return
	}
	documentKeyBytes, err := json.Marshal(documentKey)
	if err != nil {
		return
	}

	userlib.DatastoreSet(key, documentKeyBytes)

	if DEBUG_FILEMAPPING {
		fmt.Printf("\n DEBUG: USER %s STORED DOCUMENTKEY FOR FILE %s, DOCUMENTKEY %s\n", fileMapping.Username, fileMapping.Filename, documentKey.String())
	}
	return
}
