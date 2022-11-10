package helpers

import (
	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"
)

func GenerateDataStoreKey(path string) (key uuid.UUID, err error) {
	var hashedPath = userlib.Hash([]byte(path))
	var shortenedHashedPath = make([]byte, 16)
	for i := 0; i < 16; i++ {
		shortenedHashedPath[i] = 0
		for j := 0; j < 4; j++ {
			shortenedHashedPath[i] += hashedPath[i+16*j]
		}
	}
	key, err = uuid.FromBytes(shortenedHashedPath)
	return
}
