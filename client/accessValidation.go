package client

import (
	"errors"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"
)

type AccessValidation struct {
	DocumentKey       uuid.UUID
	ValidatedUsername string
	SignedValidation  []byte
	FromUsername      string
}

const DEBUG_ACCESS_VALIDATION = false

func getAccessValidationPath(documentKey uuid.UUID, validatedUsername string) string {
	return "AccessValidation/" + validatedUsername + documentKey.String()
}

func InitAccessValidation(documentKey uuid.UUID, validatedUsername string, user *User) (accessValidation AccessValidation, err error) {
	var signedValidation []byte
	var username string
	if user != nil {
		signedValidation, err = userlib.DSSign(user.SignKey, append(documentKey[:], []byte(validatedUsername)...))
		if err != nil {
			return
		}

		username = user.Username
	}
	accessValidation = AccessValidation{
		DocumentKey:       documentKey,
		ValidatedUsername: validatedUsername,
		SignedValidation:  signedValidation,
		FromUsername:      username,
	}
	return
}

func (accessValidation AccessValidation) Store() (err error) {
	path := getAccessValidationPath(accessValidation.DocumentKey, accessValidation.ValidatedUsername)

	encryptedAccessValidation, err := MarshalAndEncrypt([]byte(path), accessValidation)
	if err != nil {
		return
	}

	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	err = DatastoreSet(key, encryptedAccessValidation)
	if err != nil {
		return
	}
	return
}

func (accessValidation AccessValidation) validate() (err error) {
	verifyKey, ok := userlib.KeystoreGet(GetDSKeyStorePath(accessValidation.FromUsername))
	if DEBUG_ACCESS_VALIDATION {
		userlib.DebugMsg(
			"\nValidating access to documentKey %s for user %s\n",
			accessValidation.DocumentKey.String(),
			accessValidation.ValidatedUsername,
		)
	}
	if !ok {
		err = errors.New("VERIFY KEY OF USER NOT FOUND")
		return
	}
	err = userlib.DSVerify(
		verifyKey,
		append(accessValidation.DocumentKey[:], []byte(accessValidation.ValidatedUsername)...),
		accessValidation.SignedValidation,
	)

	return
}

func loadAccessValidation(validatedUsername string, documentKey uuid.UUID) (accessValidation AccessValidation, err error) {
	path := getAccessValidationPath(documentKey, validatedUsername)
	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	marshalledAccessValidation, ok, err := DatastoreGet(key)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("ACCESS VALIDATION NOT FOUND")
		return
	}

	err = UnmarshalAndDecrypt([]byte(path), marshalledAccessValidation, &accessValidation)
	if err != nil {
		return
	}
	return
}

func CheckAccessValidation(validatedUsername string, documentKey uuid.UUID) (err error) {
	accessValidation, err := loadAccessValidation(validatedUsername, documentKey)
	if err != nil {
		return
	}

	if accessValidation.FromUsername != "" {
		err = accessValidation.validate()
		if err != nil {
			return err
		}
		err = CheckAccessValidation(accessValidation.FromUsername, documentKey)
		if err != nil {
			return err
		}
	} else {
		doc, err := LoadDocument(documentKey)
		if err != nil {
			return err
		}
		err = doc.IsOwner(accessValidation.ValidatedUsername)
		if err != nil {
			return err
		}
	}
	return
}

func RemoveAccessValidation(documentKey uuid.UUID, toUsername string) (err error) {
	path := getAccessValidationPath(documentKey, toUsername)
	key, err := GenerateDataStoreKey(path)
	if err != nil {
		return
	}
	err = DatastoreDelete(key)
	return
}
