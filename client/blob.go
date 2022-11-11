package client

import userlib "github.com/cs161-staff/project2-userlib"

type Blob []byte

type EncryptedBlob []byte

func (blob Blob) Encrypt(symKey []byte) (encryptedBlob EncryptedBlob) {
	iv := userlib.RandomBytes(16)
	encryptedBlob = userlib.SymEnc(symKey, iv, blob)
	return
}

func (encryptedBlob EncryptedBlob) Decrypt(symKey []byte) (blob Blob) {
	blob = userlib.SymDec(symKey, encryptedBlob)
	return
}
