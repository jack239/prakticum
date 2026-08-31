package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	src := []byte("Ключ от сердца") // данные, которые хотим зашифровать
	fmt.Printf("original: %s\n", src)

	// будем использовать AES-256, создав ключ длиной 32 байта
	key, err := generateRandom(2 * aes.BlockSize) // ключ шифрования
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// NewCipher создает и возвращает новый cipher.Block.
	// Ключевым аргументом должен быть ключ AES, 16, 24 или 32 байта
	// для выбора AES-128, AES-192 или AES-256.
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
	fmt.Printf("encrypted: %x\n", dst)

	src2, err := aesgcm.Open(nil, nonce, dst, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("decrypted: %s\n", src2)
}
