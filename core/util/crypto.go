/*
	util.go
	Purpose: Crypto utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package util

import (
	"crypto/rand"
	"encoding/base64"
	"unsafe"

	"golang.org/x/crypto/bcrypt"
)

var alphanum = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ12")

// RandStr return a random string containing 0-9a-zA-Z with length of the given size
func RandStr(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphanum[b[i]/4]
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Decrypt decrypts a message that is encrypted by the [Encrypt] function
func Decrypt(msg string, salt string) (string, error) {
	deb64, err := base64.URLEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}
	dec := xor(deb64, []byte(salt))
	return string(dec), nil
}

// Encrypt encrypts the msg using a xor method which depends on the salt
//
// This encryption is weak and should not use on confidential data.
func Encrypt(msg string, salt string) string {
	enc := xor([]byte(msg), []byte(salt))
	return base64.URLEncoding.EncodeToString(enc)
}

func xor(msg, key []byte) []byte {
	ret := make([]byte, len(msg))
	for i := 0; i < len(msg); i++ {
		ret[i] = msg[i] ^ key[i%len(key)]
	}
	return ret
}

// HashPassword returns the bcrypt hash of the password so that the system will not store the real password.
func HashPassword(password string) []byte {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes
}

// CompareHashed compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure.
func CompareHashed(hashed []byte, pwd string) error {
	return bcrypt.CompareHashAndPassword(hashed, []byte(pwd))
}
