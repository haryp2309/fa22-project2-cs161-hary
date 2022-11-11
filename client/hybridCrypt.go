package client

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
)

type hybridCrypt struct {
	SymEncryptedMessage []byte
	PpkEncryptedSymKey  []byte
}

func hybridEncryptPure(publicKey userlib.PublicKeyType, message []byte) (hybridCryptObj hybridCrypt, err error) {
	symKey := userlib.RandomBytes(16)
	iv := userlib.RandomBytes(16)
	encryptedSymkey, err := userlib.PKEEnc(publicKey, symKey)
	hybridCryptObj = hybridCrypt{
		SymEncryptedMessage: userlib.SymEnc(symKey, iv, message),
		PpkEncryptedSymKey:  encryptedSymkey,
	}
	return
}

func (hybridCrypt hybridCrypt) hybridDecryptPure(privateKey userlib.PrivateKeyType) (message []byte, err error) {
	symKey, err := userlib.PKEDec(privateKey, hybridCrypt.PpkEncryptedSymKey)
	message = userlib.SymDec(symKey, hybridCrypt.SymEncryptedMessage)
	return
}

func HybridEncrypt(publicKey userlib.PublicKeyType, message []byte) (cipherObj []byte, err error) {
	hybridCryptObj, err := hybridEncryptPure(publicKey, message)
	if err != nil {
		return
	}
	cipherObj, err = json.Marshal(hybridCryptObj)
	return
}

func HybridDecrypt(privateKey userlib.PrivateKeyType, cipheObj []byte) (message []byte, err error) {
	var hybridCryptObj hybridCrypt
	err = json.Unmarshal(cipheObj, &hybridCryptObj)
	if err != nil {
		return
	}
	message, err = hybridCryptObj.hybridDecryptPure(privateKey)
	return
}
