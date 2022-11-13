package client

import (
	"strconv"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	"errors"
	"strings"
)

type Blocks [][]byte

const DEBUG_BLOCKS = false

func getBlockPath(documentKey uuid.UUID, blockPosition int) (path string) {
	path = documentKey.String() + strconv.Itoa(blockPosition)
	return
}

func (blocks Blocks) MergeToBlob() (document []byte) {
	document = make([]byte, 0)
	for _, block := range blocks {
		document = append(document, block...)
	}
	return
}

func LoadBlocks(documentKey uuid.UUID, blocksCount int) (blocks Blocks, err error) {

	blocks = make(Blocks, 0)

	for i := 0; i < blocksCount; i++ {
		blockKey, err := GenerateDataStoreKey(getBlockPath(documentKey, i))
		if err != nil {
			return nil, err
		}
		block, ok, err := DatastoreGet(blockKey)
		if err != nil {
			return nil, err
		}
		if DEBUG_BLOCKS {
			userlib.DebugMsg("Loading a block...")
		}

		if !ok {
			err = errors.New(strings.ToTitle("block not found"))
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return

}

func (blocks Blocks) Store(documentKey uuid.UUID, blockStartPosition int) (blocksCount int, err error) {
	blocksCount = blockStartPosition
	for i, block := range blocks {
		if DEBUG_BLOCKS {
			userlib.DebugMsg("Storing a block...")
		}
		path := getBlockPath(documentKey, blockStartPosition+i)
		blockKey, err := GenerateDataStoreKey(path)
		if err != nil {
			return 0, err
		}
		err = DatastoreSet(blockKey, block)
		if err != nil {
			return 0, err
		}
		blocksCount++
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
