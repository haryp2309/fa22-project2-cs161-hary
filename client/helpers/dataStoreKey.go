package helpers

import (
	"github.com/google/uuid"
)

func GenerateDataStoreKey(path string) (key uuid.UUID, err error) {
	shortenedHashedPath := ByteToHash16([]byte(path))
	key, err = uuid.FromBytes(shortenedHashedPath)
	return
}
