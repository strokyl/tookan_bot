package generate_hash

import (
	"crypto/sha1"
	b64 "encoding/base64"
)

type Hashed string
type Secret string

func GenerateHash(salt Secret, secret string) Hashed {
	data := []byte(string(salt)+secret)
	bytes := sha1.Sum(data)
	return Hashed(b64.StdEncoding.EncodeToString(bytes[:]))
}
