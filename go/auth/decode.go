package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	password = "x35k9f"
	msg      = `0ba7cd8c624345451df4710b81d1a349ce401e61bc7eb704ca` +
		`a84a8cde9f9959699f75d0d1075d676f1fe2eb475cf81f62ef` +
		`f701fee6a433cfd289d231440cf549e40b6c13d8843197a95f` +
		`8639911b7ed39a3aec4dfa9d286095c705e1a825b10a9104c6` +
		`be55d1079e6c6167118ac91318fe`
)

func main() {
	// 1) получите ключ из password, используя sha256.Sum256
	key := sha256.Sum256([]byte(password))

	// 2) создайте aesblock и aesgcm
	aesblock, _ := aes.NewCipher(key[:])
	aesgcm, _ := cipher.NewGCM(aesblock)
	// 3) получите вектор инициализации aesgcm.NonceSize() байт с конца ключа
	nonce := key[len(key)-aesgcm.NonceSize():]

	// 4) декодируйте сообщение msg в двоичный формат
	bin, _ := hex.DecodeString(msg)

	// 5) расшифровываем
	text, _ := aesgcm.Open(nil, nonce, bin, nil)

	fmt.Println(string(text))
}
