package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"math/big"
)

// HashOTP computes an HMAC-SHA256 digest for a one-time code using the given
// pepper and the tuple userID:purpose:code. It returns the raw hash bytes
func HashOTP(pepper []byte, userID, purpose, code string) []byte {
	mac := hmac.New(sha256.New, pepper)
	mac.Write([]byte(userID))
	mac.Write([]byte(":"))
	mac.Write([]byte(purpose))
	mac.Write([]byte(":"))
	mac.Write([]byte(code))
	return mac.Sum(nil)
}

// EqualHashOTP compares two OTP hashes in constant time and returns true when they are equal
func EqualHashOTP(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// GenNumericCode returns a numeric string with the provided length
func GenNumericCode(length int) string {
	if length < 1 {
		panic("crypto: invalid given length")
	}
	buf := make([]byte, length)
	for i := range length {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		buf[i] = byte('0' + n.Int64())
	}
	return string(buf)
}
