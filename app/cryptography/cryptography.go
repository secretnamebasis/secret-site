package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

const HashLength = 32

type Hash [HashLength]byte

var ZEROHASH Hash

func HashString(uniqueData string) string {

	// Hash the unique data using SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(uniqueData))
	hash := hasher.Sum(nil)

	// Truncate the hash to 32 bytes
	truncatedHash := hash[:32]

	// Convert the truncated hash to a hexadecimal string
	return hex.EncodeToString(truncatedHash)
}

// Encrypt data using AES encryption with a password
func EncryptData(data []byte, password string) ([]byte, error) {
	// Generate a key from the password
	key := []byte(password)

	// Pad the key if its length is not valid
	paddedKey := make([]byte, 32)
	copy(paddedKey, key)

	block, err := aes.NewCipher(paddedKey)
	if err != nil {
		log.Printf("Error creating cipher block: %v", err)
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("Error generating IV: %v", err)
		return nil, err
	}

	// Encrypt the data
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// Prepend the IV to the ciphertext
	ciphertext = append(iv, ciphertext[aes.BlockSize:]...)

	return ciphertext, nil
}

// Decrypt data using AES decryption with a password
func DecryptData(ciphertext []byte, password string) ([]byte, error) {
	// Generate a key from the password
	key := []byte(password)

	// Pad the key if its length is not valid
	paddedKey := make([]byte, 32)
	copy(paddedKey, key)

	block, err := aes.NewCipher(paddedKey)
	if err != nil {
		log.Printf("Error creating cipher block: %v", err)
		return nil, fmt.Errorf("error creating cipher block: %v", err)
	}

	if len(ciphertext) < aes.BlockSize {
		log.Printf("Ciphertext too short")
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract the IV from the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt the data
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
