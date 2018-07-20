package tools

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"errors"
	"crypto/rand"
	"flag"
	"io/ioutil"
)

var decrypted string
var privateKey, publicKey []byte

func init() {
	var err error
	flag.StringVar(&decrypted, "d", "", "加密过的数据")
	flag.Parse()
	publicKey, err = ioutil.ReadFile("./config/RSA_key/public.pem")
	if err != nil {
		os.Exit(-1)
	}
	privateKey,err = ioutil.ReadFile("./config/RSA_key/private.pem")
	if err != nil {
		os.Exit(-1)
	}
}

// encoder
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// decoder
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}