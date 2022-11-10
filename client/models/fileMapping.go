package models

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	"fa22-project2-cs161-hary/client/helpers"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

type FileMapping struct {
	Username    string
	Filename    string
	documentKey uuid.UUID // Automatically set when stored
}

func (fileMapping FileMapping) LoadDocumentKey() (documentKey uuid.UUID, err error) {

	if fileMapping.documentKey != uuid.Nil {
		documentKey = fileMapping.documentKey
		return
	}

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

	return
}

func (fileMapping *FileMapping) StoreDocumentKey(documentKey uuid.UUID) {
	key, err := helpers.GenerateDataStoreKey(fileMapping.Filename + fileMapping.Username)
	if err != nil {
		return
	}
	documentKeyBytes, err := json.Marshal(documentKey)
	if err != nil {
		return
	}

	userlib.DatastoreSet(key, documentKeyBytes)

	fileMapping.documentKey = documentKey
}
