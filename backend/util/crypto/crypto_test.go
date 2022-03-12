package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	helloWorld     = "Hello World!"
	publicKeyPath  = "./testdata/test_public_key.pem"
	privateKeyPath = "./testdata/test_private_key.pem"
)

func TestEncryptAndDecrypt(t *testing.T) {
	encryptData, err := Encrypt(helloWorld, publicKeyPath)
	assert.NoError(t, err)

	decryptData, err := Decrypt(encryptData, privateKeyPath)
	assert.NoError(t, err)
	assert.Equal(t, helloWorld, decryptData)
}
