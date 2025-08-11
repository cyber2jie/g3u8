package util

import (
	"encoding/hex"
	"strings"
)

func Hex2Byte(hexStr string) ([]byte, error) {
	hexStr = strings.TrimSpace(hexStr)
	if strings.HasPrefix(hexStr, "0x") || strings.HasPrefix(hexStr, "0X") {
		hexStr = hexStr[2:]
	}
	return hex.DecodeString(hexStr)
}

func Byte2Hex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
