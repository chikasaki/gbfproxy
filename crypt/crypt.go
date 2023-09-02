package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"log"
	"os"
)

type Config struct {
	KeyFile string
}

type Crypt interface {
	Encrypt(input []byte) []byte
	Decrypt(input []byte) []byte
}

type AESCrypt struct {
	key []byte
}

func NewCrypt(config *Config) (Crypt, error) {
	key, err := os.ReadFile(config.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	return &AESCrypt{
		key: key,
	}, nil
}

func (c *AESCrypt) Encrypt(input []byte) (output []byte) {
	block, _ := aes.NewCipher(c.key)
	blockSize := block.BlockSize()
	input = PKCS5Padding(input, blockSize)
	output = make([]byte, len(input))
	encrypter := cipher.NewCBCEncrypter(block, c.key[:blockSize])
	encrypter.CryptBlocks(output, input)
	return
}

func (c *AESCrypt) Decrypt(input []byte) (output []byte) {
	output = make([]byte, len(input))
	block, _ := aes.NewCipher(c.key)
	blockSize := block.BlockSize()
	decrypter := cipher.NewCBCDecrypter(block, c.key[:blockSize])
	decrypter.CryptBlocks(output, input)
	output = PKCS5UnPadding(output)
	return
}

func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
