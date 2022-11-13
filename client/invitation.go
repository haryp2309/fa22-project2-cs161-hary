package client

import (
	"errors"

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
	serializedInvitation, ok, err := DatastoreGet(invitationPointer)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("INVITATION NOT FOUND")
		return
	}
	err = UnmarshalAndDecrypt([]byte(recievingUsername), serializedInvitation, &invitation)
	return
}

func (invitation Invitation) Store(recievingUsername string) (invitationPointer uuid.UUID, err error) {
	serializedInvitation, err := MarshalAndEncrypt([]byte(recievingUsername), invitation)
	if err != nil {
		return
	}
	invitationPointer = uuid.New()
	err = DatastoreSet(invitationPointer, serializedInvitation)
	if err != nil {
		return
	}
	return
}
