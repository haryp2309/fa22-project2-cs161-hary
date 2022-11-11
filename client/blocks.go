package client

import (
	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	"errors"
	"strings"
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

func (blocks Blocks) Encrypt(access Access) (encryptedBlocks Blocks) {
	encryptedBlocks = make(Blocks, 0)
	for _, block := range blocks {
		iv := userlib.RandomBytes(16)
		encryptedBlob := userlib.SymEnc(access.EncryptionKey, iv, block)
		encryptedBlocks = append(encryptedBlocks, encryptedBlob)
	}
	return
}

func (encryptedBlocks Blocks) Decrypt(access Access) (decryptedBlocks Blocks) {
	decryptedBlocks = make(Blocks, 0)
	for _, encryptedBlock := range encryptedBlocks {
		decryptedBlock := userlib.SymDec(access.EncryptionKey, encryptedBlock)
		decryptedBlocks = append(decryptedBlocks, decryptedBlock)
	}
	return
}

const MAX_BLOCK_SIZE = 32

func SplitBlobToBlocks(blob []byte) (blocks Blocks) {
	blocks = make(Blocks, 0)
	remainingDocument := blob
	for len(remainingDocument) > 0 {
		end := MAX_BLOCK_SIZE
		if len(remainingDocument) < end {
			end = len(remainingDocument)
		}

		blocks = append(blocks, remainingDocument[:end])
		remainingDocument = remainingDocument[end:]
	}
	return
}
