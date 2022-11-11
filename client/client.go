package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client/helpers"

	"github.com/google/uuid"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string
	PrivKey  userlib.PrivateKeyType
	SignKey  userlib.DSSignKey
}

func getUserKey(username string, password string) (argonKey []byte) {
	argonKey = userlib.Argon2Key([]byte(password), []byte(username), 16)
	return
}

func doUserExist(username string) (ok bool) {
	path := helpers.GetPKKeyStorePath(username)
	_, ok = userlib.KeystoreGet(path)
	return
}

func InitUser(username string, password string) (userdata *User, err error) {
	var user User
	userdata = &user
	userdata.Username = username

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

	userlib.KeystoreSet(helpers.GetPKKeyStorePath(username), pub)

	var verifyKey userlib.DSVerifyKey
	userdata.SignKey, verifyKey, err = userlib.DSKeyGen()
	if err != nil {
		return
	}
	userlib.KeystoreSet(helpers.GetDSKeyStorePath(username), verifyKey)

	var jsonUser []byte
	jsonUser, err = json.Marshal(userdata)

	if err != nil {
		return
	}

	var encJsonUser = userlib.SymEnc(userKey, userlib.RandomBytes(16), jsonUser)

	userlib.DatastoreSet(key, encJsonUser)

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
	encJsonUser, ok = userlib.DatastoreGet(key)

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
	return userdataptr, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	fileMapping := InitFileMapping(
		userdata.Username,
		filename,
	)

	err = fileMapping.StoreDocumentKey(uuid.New())
	if err != nil {
		return err
	}

	documentKey, err := fileMapping.LoadDocumentKey()
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

	blockKeys, err := blocks.Store()
	if err != nil {
		return err
	}

	document := InitDocument(
		blockKeys,
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

	documentKey, err := fileMapping.LoadDocumentKey()
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

	newBlockKeys, err := blocks.Store()
	if err != nil {
		return err
	}

	document, err := LoadDocument(documentKey)
	if err != nil {
		return
	}

	document.BlockKeys = append(document.BlockKeys, newBlockKeys...)

	err = document.Store(documentKey)

	return
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {

	fileMapping := InitFileMapping(
		userdata.Username,
		filename,
	)

	documentKey, err := fileMapping.LoadDocumentKey()
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

	blocks, err := LoadBlocks(document.BlockKeys)

	content = blocks.Decrypt(access).MergeToBlob()
	return content, err
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	documentKey, err := InitFileMapping(userdata.Username, filename).LoadDocumentKey()
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
	err = InitFileMapping(userdata.Username, filename).StoreDocumentKey(invitation.DocumentKey)
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

	documentKey, err := fileMapping.LoadDocumentKey()
	if err != nil {
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
