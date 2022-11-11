package helpers

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
)

func MarshalAndEncrypt(byteKey []byte, value interface{}) (out []byte, err error) {
	key := ByteToHash16(byteKey)
	iv := userlib.RandomBytes(16)

	marshalledValue, err := json.Marshal(value)
	if err != nil {
		return
	}

	out = userlib.SymEnc(key, iv, marshalledValue)
	return
}

func UnmarshalAndDecrypt(byteKey []byte, cipherValue []byte, out interface{}) (err error) {
	key := ByteToHash16(byteKey)
	marshalledValue := userlib.SymDec(key, cipherValue)

	err = json.Unmarshal(marshalledValue, out)
	return

}
