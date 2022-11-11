package client

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
)

type Signature struct {
	ToUsername    string
	EncryptionKey []byte
}

func (signature Signature) Sign(signKey userlib.DSSignKey) (signedSignature []byte, err error) {

	marshalledSignature, err := json.Marshal(signature)
	if err != nil {
		return
	}
	err = signKey.PrivKey.Validate()
	if err != nil {
		return
	}
	signedSignature, err = userlib.DSSign(signKey, marshalledSignature)

	return
}

func (signature Signature) Verify(verifyKey userlib.DSVerifyKey, signedSignature []byte) (err error) {
	marshalledSignature, err := json.Marshal(signature)
	if err != nil {
		return
	}

	err = userlib.DSVerify(verifyKey, marshalledSignature, signedSignature)
	return
}

func InitSignature(toUsername string, encryptionKey []byte) (signature Signature) {
	return Signature{
		ToUsername:    toUsername,
		EncryptionKey: encryptionKey,
	}
}
