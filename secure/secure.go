package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func DecryptAES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func CustomEncrypt(plainText, key string) string {
	encrypted := make([]byte, len(plainText))

	keyIndex := 0
	for i := 0; i < len(plainText); i++ {
		char := plainText[i]
		if 'a' <= char && char <= 'z' {
			shift := key[keyIndex] - 'A'
			encrypted[i] = s(char, 'a', int(shift))
			keyIndex = (keyIndex + 1) % len(key)
		} else if 'A' <= char && char <= 'Z' {
			shift := key[keyIndex] - 'A'
			encrypted[i] = s(char, 'A', int(shift))
			keyIndex = (keyIndex + 1) % len(key)
		} else {
			encrypted[i] = char
		}
	}

	return string(encrypted)
}

func s(letter byte, base byte, shift int) byte {
	shifted := int(letter-base+byte(shift)) % 26
	if shifted < 0 {
		shifted += 26
	}
	return base + byte(shifted)
}

func CustomDecrypt(encrypted, key string) string {
	decrypted := make([]byte, len(encrypted))

	keyIndex := 0
	for i := 0; i < len(encrypted); i++ {
		char := encrypted[i]
		if 'a' <= char && char <= 'z' {
			shift := key[keyIndex] - 'A'
			decrypted[i] = s(char, 'a', -int(shift))
			keyIndex = (keyIndex + 1) % len(key)
		} else if 'A' <= char && char <= 'Z' {
			shift := key[keyIndex] - 'A'
			decrypted[i] = s(char, 'A', -int(shift))
			keyIndex = (keyIndex + 1) % len(key)
		} else {
			decrypted[i] = char
		}
	}

	return string(decrypted)
}
