package client

import (
	"fmt"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	"github.com/cs161-staff/project2-starter-code/client/helpers"

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

func getFileMappingPath(filename string, username string) string {
	return "Documents/" + filename + username
}

const DEBUG_FILEMAPPING = false

func (fileMapping FileMapping) LoadDocumentKey() (documentKey uuid.UUID, err error) {

	path := getFileMappingPath(fileMapping.Filename, fileMapping.Username)
	key, err := helpers.GenerateDataStoreKey(path)
	if err != nil {
		return
	}

	documentKeyBytes, ok := userlib.DatastoreGet(key)

	if !ok {
		err = errors.New("WRONG USERNAME OR PASSWORD")
		return
	}

	err = helpers.UnmarshalAndDecrypt([]byte(path), documentKeyBytes, &documentKey)
	if err != nil {
		return
	}

	if DEBUG_FILEMAPPING {
		fmt.Printf("\n DEBUG: USER %s ACCESSED FILE %s, DOCUMENTKEY %s\n", fileMapping.Username, fileMapping.Filename, documentKey.String())
	}

	return
}

func (fileMapping FileMapping) StoreDocumentKey(documentKey uuid.UUID) (err error) {
	path := getFileMappingPath(fileMapping.Filename, fileMapping.Username)
	key, err := helpers.GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	documentKeyBytes, err := helpers.MarshalAndEncrypt([]byte(path), documentKey)
	if err != nil {
		return
	}

	userlib.DatastoreSet(key, documentKeyBytes)

	if DEBUG_FILEMAPPING {
		fmt.Printf("\n DEBUG: USER %s STORED DOCUMENTKEY FOR FILE %s, DOCUMENTKEY %s\n", fileMapping.Username, fileMapping.Filename, documentKey.String())
	}
	return
}
