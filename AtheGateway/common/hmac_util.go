package common

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"hash"
)

var hashAlgorithmMap = map[string]func() hash.Hash{
	"hmac-sha256": sha256.New,
	"hmac-sha512": sha512.New,
	"hmac-sha1":   sha1.New,
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(hashAlgorithmMap[algorithm], []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func SupportAlgorithm(algorithm string) bool {
	if _, ok := hashAlgorithmMap[algorithm]; !ok {
		return false
	}
	return true
}

func EncodeWithSha256(data []byte) string {
	sha := sha256.New()
	sha.Write(data)
	encodeData := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}


