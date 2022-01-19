package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
)

func Encrypt(data, publicKeyPath string) (string, error) {
	rsaPublicKey, err := readPublicKey(publicKeyPath)
	if err != nil {
		return "", err
	}

	encryptPKCS1v15, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(data))
	if err != nil {
		return "", err
	}

	encryptString := base64.StdEncoding.EncodeToString(encryptPKCS1v15)
	return encryptString, err
}

func Decrypt(base64data, privateKeyPath string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return "", err
	}

	rsaPrivateKey, err := readPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}

	decryptPKCS1v15, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, decodeString)
	return string(decryptPKCS1v15), err
}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	context, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(context)

	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	return privateKey, err
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	readFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(readFile)

	var pkixPublicKey interface{}
	if pemBlock.Type == "RSA PUBLIC KEY" {
		pkixPublicKey, err = x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	} else if pemBlock.Type == "PUBLIC KEY" {
		pkixPublicKey, err = x509.ParsePKIXPublicKey(pemBlock.Bytes)
	}
	if err != nil {
		return nil, err
	}

	publicKey := pkixPublicKey.(*rsa.PublicKey)
	return publicKey, nil
}
