package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

type EncryptMethod string

const (
	AES_128     EncryptMethod = "AES-128"
	AES_128_ECB EncryptMethod = "AES-128-ECB"
)

func EncryptMethodFromString(method string) EncryptMethod {
	switch method {
	case "AES-128":
		return AES_128
	case "AES-128-ECB":
		return AES_128_ECB
	default:
		return ""
	}
}

type DecryptKey struct {
	Key string `json:"key"`
	IV  string `json:"iv"`
}

func Decrypt(content []byte, method EncryptMethod, decryptKey DecryptKey) ([]byte, error) {
	keyBytes, err := hex.DecodeString(decryptKey.Key)
	if err != nil {
		return nil, fmt.Errorf("decode key error: %w", err)
	}

	ivBytes, err := hex.DecodeString(decryptKey.IV)
	if err != nil {
		return nil, fmt.Errorf("decode iv error: %w", err)
	}

	switch method {
	case AES_128:
		content, err = AES128Decrypt(content, keyBytes, ivBytes, cipherModeCBC)
		if err != nil {
			return nil, fmt.Errorf("aes128 decrypt error: %w", err)
		}
	case AES_128_ECB:
		content, err = AES128Decrypt(content, keyBytes, ivBytes, cipherModeECB)
		if err != nil {
			return nil, fmt.Errorf("aes128 ecb decrypt error: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported encrypt method: %d", method)
	}

	return content, nil
}

type cipherMode int

const (
	cipherModeCBC cipherMode = iota
	cipherModeECB
)

func AES128Decrypt(encryptedBuff, keyByte, ivByte []byte, mode cipherMode) ([]byte, error) {
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}

	switch mode {
	case cipherModeCBC:
		mode := cipher.NewCBCDecrypter(block, ivByte)
		result := make([]byte, len(encryptedBuff))
		mode.CryptBlocks(result, encryptedBuff)
		b, err := removePKCS7Padding(result)
		if err != nil {
			return result, nil
		}
		return b, nil
	case cipherModeECB:
		return decryptECB(block, encryptedBuff)
	default:
		return nil, fmt.Errorf("unsupported cipher mode")
	}
}

func decryptECB(block cipher.Block, ciphertext []byte) ([]byte, error) {

	result := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(result[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}

	b, err := removePKCS7Padding(result)
	if err != nil {
		return result, nil
	}
	return b, nil
}

func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
