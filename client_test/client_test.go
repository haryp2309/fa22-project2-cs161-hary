package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.

	_ "encoding/hex"
	_ "errors"
	"strconv"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			//Skip("Skipping test")
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})
	Describe("Custom Tests", func() {

		Specify("Custom Basic Test: Testing Single User Store/Load.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		Specify("Existing user should not be able to create new account", func() {
			userlib.DebugMsg("Creating user alice")
			_, err := client.InitUser("alice", "password123")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Creating user alice again")
			_, err = client.InitUser("alice", "password321")
			Expect(err).ToNot(BeNil())
		})

		Specify("Usernames are case-sensitive: Bob and bob are different users", func() {
			userlib.DebugMsg("Creating user Bob")
			_, err := client.InitUser("Bob", "password123")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Creating user bob")
			_, err = client.InitUser("bob", "password123")
			Expect(err).To(BeNil())

		})

		Specify("The client SHOULD support passwords length equal to zero.", func() {
			userlib.DebugMsg("Creating user Bob with no password")
			_, err := client.InitUser("Bob", "")
			Expect(err).To(BeNil())
		})

		Specify(`The client MUST enforce that there is only a single copy of a file. 
			Sharing the file MAY NOT create a copy of the file.`, func() {
			userlib.DebugMsg("Creating user Bob")
			bob, err := client.InitUser("Bob", "bestpassword321")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating user Eve")
			eve, err := client.InitUser("Eve", "bestpassword321")
			Expect(err).To(BeNil())

			eve.StoreFile("file123", []byte("abc"))
			inv, err := eve.CreateInvitation("file123", "Bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("Eve", inv, "file321")
			Expect(err).To(BeNil())

			err = bob.AppendToFile("file321", []byte("cba"))
			Expect(err).To(BeNil())

			byteContent, err := eve.LoadFile("file123")
			Expect(err).To(BeNil())

			content := string(byteContent)
			Expect(content).To(Equal("abccba"))
		})

		Specify("Filenames MAY be any length, including zero (empty string).", func() {
			userlib.DebugMsg("Creating user Bob")
			bob, err := client.InitUser("Bob", "bestpassword321")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating file with empty name")
			err = bob.StoreFile("", []byte(""))
			Expect(err).To(BeNil())

		})

		Specify("User should only be able to revoke access to a file if he/she directly shared to that person", func() {
			userlib.DebugMsg("Creating user Bob, Alice and Eve")
			bob, err := client.InitUser("Bob", "bestpassword321")
			Expect(err).To(BeNil())
			eve, err := client.InitUser("Eve", "bestpassword321")
			Expect(err).To(BeNil())
			alice, err := client.InitUser("Alice", "bestpassword321")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob creates a file.")
			const FILENAME = "file"
			const CONTENT = "Very interessting document about absolutely nothing."
			bob.StoreFile(FILENAME, []byte(CONTENT))

			userlib.DebugMsg("Bob creates invite to eve")
			inv, err := bob.CreateInvitation(FILENAME, "Eve")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Eve accepts the invite")
			eve.AcceptInvitation("Bob", inv, FILENAME)

			userlib.DebugMsg("Eve creates invite to Alice")
			inv, err = eve.CreateInvitation(FILENAME, "Alice")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice accepts the invite")
			alice.AcceptInvitation("Alice", inv, FILENAME)

			userlib.DebugMsg("Bob tries to revoke access for Alice")
			err = bob.RevokeAccess(FILENAME, "Alice")
			Expect(err).ToNot(BeNil())
		})

		Specify("The client MUST NOT assume that filenames are globally unique.", func() {
			userlib.DebugMsg("Creating user Bob and Alice")
			bob, err := client.InitUser("Bob", "bestpassword123")
			Expect(err).To(BeNil())
			alice, err := client.InitUser("Alice", "bestpassword123")
			Expect(err).To(BeNil())

			const FILENAME = "filename"
			const CONTENT_1 = "Very interessting document about absolutely nothing."
			const CONTENT_2 = "Very boring document about absolutely nothing."

			userlib.DebugMsg("Bob stores a file with a filename")
			err = bob.StoreFile(FILENAME, []byte(CONTENT_1))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice stores another file with same filename")
			err = alice.StoreFile(FILENAME, []byte(CONTENT_2))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loads her file")
			byteContent, err := alice.LoadFile(FILENAME)
			Expect(err).To(BeNil())
			Expect(string(byteContent)).To(Equal(CONTENT_2))

		})

		Specify("Overwriting the contents of a file does not change who the file is shared with.", func() {
			userlib.DebugMsg("Creating user Bob and Alice")
			bob, err := client.InitUser("Bob", "bestpassword123")
			Expect(err).To(BeNil())
			alice, err := client.InitUser("Alice", "bestpassword123")
			Expect(err).To(BeNil())

			const FILENAME = "filename"
			const CONTENT = "Very interessting document about absolutely nothing."

			userlib.DebugMsg("Bob stores a file")
			err = bob.StoreFile(FILENAME, []byte(CONTENT))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob shares file with Alice")
			inv, err := bob.CreateInvitation(FILENAME, "Alice")
			Expect(err).To(BeNil())
			alice.AcceptInvitation("Bob", inv, FILENAME)

			userlib.DebugMsg("Alice loads file")
			byteContent, err := alice.LoadFile(FILENAME)
			Expect(err).To(BeNil())
			Expect(string(byteContent)).To(Equal(CONTENT))

			userlib.DebugMsg("Bob stores a file again with same filename")
			err = bob.StoreFile(FILENAME, []byte(CONTENT))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loads file again")
			byteContent, err = alice.LoadFile(FILENAME)
			Expect(err).To(BeNil())
			Expect(string(byteContent)).To(Equal(CONTENT))
		})

		Specify("Malicious tampering with keystore should be detected.", func() {
			userlib.DebugMsg("Creating user Bob ")
			_, err := client.InitUser("Bob", "bestpassword123")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Eve tries to change public key of bob")
			evePub, _, err := userlib.PKEKeyGen()
			Expect(err).To(BeNil())
			keystore := userlib.KeystoreGetMap()
			userlib.KeystoreClear()
			for key := range keystore {
				err = userlib.KeystoreSet(key, evePub)
				Expect(err).To(BeNil())
			}

			_, err = client.GetUser("Bob", "bestpassword123")
			Expect(err).ToNot(BeNil())
		})

		Specify("Keystore should not scale up with number of files.", func() {
			userlib.DebugMsg("Creating user Bob ")
			bob, err := client.InitUser("Bob", "bestpassword123")
			Expect(err).To(BeNil())

			map0 := userlib.KeystoreGetMap()

			userlib.DebugMsg("Bob stores 100 files")
			const CONTENT = "Very interessting document about absolutely nothing."
			const FILENAME = "filename"

			for i := 0; i < 100; i++ {
				//userlib.DebugMsg("Bob stores file" + strconv.Itoa(i))
				err = bob.StoreFile(FILENAME+strconv.Itoa(i), []byte(CONTENT))
				Expect(err).To(BeNil())
			}

			map1 := userlib.KeystoreGetMap()
			Expect(map0).To(Equal(map1))

		})
	})
})
