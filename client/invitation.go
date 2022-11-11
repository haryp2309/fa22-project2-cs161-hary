package client

import (
	"errors"

	"github.com/cs161-staff/project2-starter-code/client/helpers"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"
)

type Invitation struct {
	DocumentKey uuid.UUID
}

func InitInvitation(documentKey uuid.UUID) Invitation {
	return Invitation{
		DocumentKey: documentKey,
	}
}

func LoadInvitation(invitationPointer uuid.UUID, recievingUsername string) (invitation Invitation, err error) {
	serializedInvitation, ok := userlib.DatastoreGet(invitationPointer)
	if !ok {
		err = errors.New("INVITATION NOT FOUND")
		return
	}
	err = helpers.UnmarshalAndDecrypt([]byte(recievingUsername), serializedInvitation, &invitation)
	return
}

func (invitation Invitation) Store(recievingUsername string) (invitationPointer uuid.UUID, err error) {
	serializedInvitation, err := helpers.MarshalAndEncrypt([]byte(recievingUsername), invitation)
	invitationPointer = uuid.New()
	userlib.DatastoreSet(invitationPointer, serializedInvitation)
	return
}
