package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"

	"fa22-project2-cs161-hary/client/helpers"

	"github.com/google/uuid"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string
	PrivKey  userlib.PrivateKeyType
	SignKey  userlib.DSSignKey

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdata *User, err error) {
	var user User
	userdata = &user
	userdata.Username = username

	var key userlib.UUID
	key, err = helpers.GenerateDataStoreKey(username + password)

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

	var passSymKey = userlib.Argon2Key([]byte(password), []byte("password?"), 16)
	var encJsonUser = userlib.SymEnc(passSymKey, userlib.RandomBytes(16), jsonUser)

	userlib.DatastoreSet(key, encJsonUser)

	return
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	key, err := helpers.GenerateDataStoreKey(username + password)
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

	var passSymKey = userlib.Argon2Key([]byte(password), []byte("password?"), 16)
	var jsonUser = userlib.SymDec(passSymKey, encJsonUser)

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
	invitationPtr, err = InitInvitation(documentKey).Store()
	return
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	invitation, err := LoadInvitation(invitationPtr)
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
