package client

import (
	"encoding/json"
	"errors"

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

func LoadInvitation(invitationPointer uuid.UUID) (invitation Invitation, err error) {
	serializedInvitation, ok := userlib.DatastoreGet(invitationPointer)
	if !ok {
		err = errors.New("INVITATION NOT FOUND")
		return
	}
	err = json.Unmarshal(serializedInvitation, &invitation)
	return
}

func (invitation Invitation) Store() (invitationPointer uuid.UUID, err error) {
	serializedInvitation, err := json.Marshal(invitation)
	invitationPointer = uuid.New()
	userlib.DatastoreSet(invitationPointer, serializedInvitation)
	return
}
