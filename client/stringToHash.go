package client

import userlib "github.com/cs161-staff/project2-userlib"

func ByteToHash16(input []byte) (hashedText16 []byte) {
	var hashedText = userlib.Hash(input)
	hashedText16 = make([]byte, 16)
	for i := 0; i < 16; i++ {
		hashedText16[i] = 0
		for j := 0; j < 4; j++ {
			hashedText16[i] += hashedText[i+16*j]
		}
	}
	return
}
