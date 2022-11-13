package client

import (
	"encoding/json"
	"errors"
	"strings"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"
)

type Access struct {
	DocumentKey   uuid.UUID
	ToUsername    string
	FromUsername  string
	EncryptionKey []byte
	Signature     []byte
}

const ACCESS_PATH_PREFIX = "access/"

func getAccessPath(documentKey uuid.UUID, toUsername string) string {
	return ACCESS_PATH_PREFIX +
		documentKey.String() +
		toUsername
}

func InitAccess(fromUser *User, toUsername string, documentKey uuid.UUID, encryptionKey []byte) (access Access, err error) {

	var fromUsername string
	var signedSignature []byte

	if fromUser != nil {
		signature := InitSignature(
			toUsername,
			encryptionKey,
		)
		signedSignature, err = signature.Sign(fromUser.SignKey)
		if err != nil {
			return
		}

		fromUsername = fromUser.Username
	}

	access = Access{
		EncryptionKey: encryptionKey,
		ToUsername:    toUsername,
		FromUsername:  fromUsername,
		DocumentKey:   documentKey,
		Signature:     signedSignature,
	}
	return
}

func (access *Access) StoreAccess() (err error) {

	jsonAccess, err := json.Marshal(access)
	if err != nil {
		return
	}

	pubKey, ok := userlib.KeystoreGet(GetPKKeyStorePath(access.ToUsername))
	if !ok {
		err = errors.New(strings.ToTitle("public key not found"))
		return
	}
	encryptedAccess, err := HybridEncrypt(pubKey, jsonAccess)
	if err != nil {
		return
	}

	path := getAccessPath(access.DocumentKey, access.ToUsername)
	key, err := GenerateDataStoreKey(path)

	err = DatastoreSet(key, encryptedAccess)
	if err != nil {
		return
	}

	return

}

func LoadAccess(documentKey uuid.UUID, user User) (access Access, err error) {
	path := getAccessPath(documentKey, user.Username)
	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	encryptedAccess, ok, err := DatastoreGet(key)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(strings.ToTitle("access not found"))
		return
	}

	jsonAccess, err := HybridDecrypt(user.PrivKey, encryptedAccess)
	if err != nil {
		return
	}

	json.Unmarshal(jsonAccess, &access)

	err = access.isValid()
	return
}

func (access Access) isValid() (err error) {
	if access.FromUsername == "" {
		doc, erro := LoadDocument(access.DocumentKey)
		if erro != nil {
			return erro
		}

		erro = doc.IsOwner(access.ToUsername)
		if erro != nil {
			return erro
		}
		return
	}
	verifyKey, ok := userlib.KeystoreGet(GetDSKeyStorePath(access.FromUsername))
	if !ok {
		err = errors.New(strings.ToTitle("access not found"))
		return
	}

	signature := InitSignature(
		access.ToUsername,
		access.EncryptionKey,
	)

	err = signature.Verify(verifyKey, access.Signature)

	return
}

func RemoveAccess(documentKey uuid.UUID, toUsername string) (err error) {
	path := getAccessPath(documentKey, toUsername)
	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	userlib.DatastoreDelete(key)
	return
}
