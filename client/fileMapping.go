package client

import (
	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"
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

func (fileMapping FileMapping) LoadDocumentKey(symKey []byte) (documentKey uuid.UUID, ok bool, err error) {

	path := getFileMappingPath(fileMapping.Filename, fileMapping.Username)
	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}

	documentKeyBytes, ok, err := DatastoreGet(key)
	if err != nil || !ok {
		return
	}

	err = UnmarshalAndDecrypt(symKey, documentKeyBytes, &documentKey)
	if err != nil {
		return
	}

	if DEBUG_FILEMAPPING {
		userlib.DebugMsg("\n DEBUG: USER %s ACCESSED FILE %s, DOCUMENTKEY %s\n", fileMapping.Username, fileMapping.Filename, documentKey.String())
	}

	return
}

func (fileMapping FileMapping) StoreDocumentKey(documentKey uuid.UUID, symKey []byte) (err error) {
	path := getFileMappingPath(fileMapping.Filename, fileMapping.Username)
	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	documentKeyBytes, err := MarshalAndEncrypt(symKey, documentKey)
	if err != nil {
		return
	}

	err = DatastoreSet(key, documentKeyBytes)
	if err != nil {
		return
	}

	if DEBUG_FILEMAPPING {
		userlib.DebugMsg("\n DEBUG: USER %s STORED DOCUMENTKEY FOR FILE %s, DOCUMENTKEY %s\n", fileMapping.Username, fileMapping.Filename, documentKey.String())
	}
	return
}
