package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/google/uuid"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).

	// Useful for creating new error messages to return using errors.New("...")
	"errors"
)

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string
	PrivKey  userlib.PrivateKeyType
	SignKey  userlib.DSSignKey
	SymKey   []byte
}

func getUserKey(username string, password string) (argonKey []byte) {
	argonKey = userlib.Argon2Key([]byte(password), []byte(username), 16)
	return
}

func doUserExist(username string) (ok bool) {
	path := GetPKKeyStorePath(username)
	_, ok = userlib.KeystoreGet(path)
	return
}

func (user User) validateUserKeys() (err error) {
	const VALIDATION_FAILED = "USER VALIDATION FAILED; MALICIOUS ACTIVITY DETECTED; "
	pub, ok := userlib.KeystoreGet(GetPKKeyStorePath(user.Username))
	if !ok {
		err = errors.New(VALIDATION_FAILED + "PUBLIC PKKEY NOT FOUND")
		return
	}

	const CONTENT = "CONTENT"
	encContent, err := userlib.PKEEnc(pub, []byte(CONTENT))
	if err != nil {
		err = errors.New(VALIDATION_FAILED + err.Error())
		return
	}
	decContent, err := userlib.PKEDec(user.PrivKey, encContent)
	if err != nil {
		err = errors.New(VALIDATION_FAILED + err.Error())
		return
	}
	if string(decContent) != CONTENT {
		err = errors.New(VALIDATION_FAILED + "PRIVATE AND PUBLIC KEY DOES NOT MATCH")
		return
	}

	return

}

func InitUser(username string, password string) (userdata *User, err error) {
	var user User
	userdata = &user
	userdata.Username = username
	userdata.SymKey = userlib.RandomBytes(16)

	if doUserExist(username) {
		err = errors.New(ERROR_USER_EXISTS)
		return
	}
	userKey := getUserKey(username, password)
	key, err := uuid.FromBytes(userKey)

	if err != nil {
		return
	}

	var pub userlib.PublicKeyType
	pub, userdata.PrivKey, err = userlib.PKEKeyGen()

	if err != nil {
		return
	}

	userlib.KeystoreSet(GetPKKeyStorePath(username), pub)

	var verifyKey userlib.DSVerifyKey
	userdata.SignKey, verifyKey, err = userlib.DSKeyGen()
	if err != nil {
		return
	}
	userlib.KeystoreSet(GetDSKeyStorePath(username), verifyKey)

	var jsonUser []byte
	jsonUser, err = json.Marshal(userdata)

	if err != nil {
		return
	}

	var encJsonUser = userlib.SymEnc(userKey, userlib.RandomBytes(16), jsonUser)

	err = DatastoreSet(key, encJsonUser)
	if err != nil {
		return
	}

	err = user.validateUserKeys()
	if err != nil {
		return
	}

	return
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	userKey := getUserKey(username, password)
	key, err := uuid.FromBytes(userKey)
	if err != nil {
		return
	}

	var encJsonUser []byte
	var ok bool
	encJsonUser, ok, err = DatastoreGet(key)
	if err != nil {
		return
	}

	if !ok {
		err = errors.New("WRONG USERNAME OR PASSWORD")
		return
	}

	var jsonUser = userlib.SymDec(userKey, encJsonUser)

	err = json.Unmarshal(jsonUser, &userdata)
	if err != nil {
		return
	}

	userdataptr = &userdata

	err = userdata.validateUserKeys()
	if err != nil {
		return
	}

	return userdataptr, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	fileMapping := InitFileMapping(
		userdata.Username,
		filename,
	)

	err = fileMapping.StoreDocumentKey(uuid.New(), userdata.SymKey)
	if err != nil {
		return err
	}

	documentKey, err := fileMapping.LoadDocumentKey(userdata.SymKey)
	if err != nil {
		return err
	}

	encryptionKey := userlib.RandomBytes(16)

	validation, err := InitAccessValidation(documentKey, userdata.Username, nil)
	if err != nil {
		return
	}
	err = validation.Store()
	if err != nil {
		return
	}

	access, err := InitAccess(nil, userdata.Username, documentKey, encryptionKey)
	if err != nil {
		return
	}

	err = access.StoreAccess()
	if err != nil {
		return
	}

	blocks := SplitBlobToBlocks(content).Encrypt(access)

	blocksCount, err := blocks.Store(documentKey, 0)
	if err != nil {
		return err
	}

	document := InitDocument(
		blocksCount,
		userdata.Username,
	)

	err = document.Store(documentKey)

	return
}

func (userdata *User) AppendToFile(filename string, content []byte) (err error) {

	fileMapping := InitFileMapping(
		userdata.Username,
		filename,
	)

	documentKey, err := fileMapping.LoadDocumentKey(userdata.SymKey)
	if err != nil {
		return
	}

	err = CheckAccessValidation(userdata.Username, documentKey)
	if err != nil {
		return
	}

	access, err := LoadAccess(documentKey, *userdata)
	if err != nil {
		return
	}

	blocks := SplitBlobToBlocks(content).Encrypt(access)

	document, err := LoadDocument(documentKey)
	if err != nil {
		return
	}

	document.BlocksCount, err = blocks.Store(documentKey, document.BlocksCount)
	if err != nil {
		return err
	}

	err = document.Store(documentKey)

	return
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {

	fileMapping := InitFileMapping(
		userdata.Username,
		filename,
	)

	documentKey, err := fileMapping.LoadDocumentKey(userdata.SymKey)
	if err != nil {
		return
	}

	err = CheckAccessValidation(userdata.Username, documentKey)
	if err != nil {
		return
	}

	document, err := LoadDocument(documentKey)
	if err != nil {
		return nil, err
	}

	access, err := LoadAccess(documentKey, *userdata)
	if err != nil {
		return
	}

	blocks, err := LoadBlocks(documentKey, document.BlocksCount)

	content = blocks.Decrypt(access).MergeToBlob()
	return content, err
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	documentKey, err := InitFileMapping(userdata.Username, filename).LoadDocumentKey(userdata.SymKey)
	if err != nil {
		return
	}
	accessValidation, err := InitAccessValidation(documentKey, recipientUsername, userdata)
	if err != nil {
		return
	}
	err = accessValidation.Store()
	if err != nil {
		return
	}

	parentAccess, err := LoadAccess(documentKey, *userdata)
	if err != nil {
		return
	}
	access, err := InitAccess(userdata, recipientUsername, documentKey, parentAccess.EncryptionKey)
	if err != nil {
		return
	}
	err = access.StoreAccess()
	if err != nil {
		return
	}
	invitationPtr, err = InitInvitation(documentKey).Store(recipientUsername)
	return
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	invitation, err := LoadInvitation(invitationPtr, userdata.Username)
	if err != nil {
		return err
	}
	err = InitFileMapping(userdata.Username, filename).StoreDocumentKey(invitation.DocumentKey, userdata.SymKey)
	if err != nil {
		return err
	}
	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {

	fileMapping := InitFileMapping(
		userdata.Username,
		filename,
	)

	documentKey, err := fileMapping.LoadDocumentKey(userdata.SymKey)
	if err != nil {
		return err
	}

	accessValidation, err := loadAccessValidation(recipientUsername, documentKey)
	if err != nil {
		return err
	}

	if accessValidation.FromUsername != userdata.Username {
		err = errors.New("USER DOES NOT HAVE DIRECT ACCESS")
		return err
	}

	err = RemoveAccess(documentKey, recipientUsername)
	if err != nil {
		return err
	}

	err = RemoveAccessValidation(documentKey, recipientUsername)
	if err != nil {
		return err
	}

	return nil
}
