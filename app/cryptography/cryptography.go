package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const HashLength = 32

type Hash [HashLength]byte

var ZEROHASH Hash

// HashString calculates the SHA-256 hash of the input string and returns it as a hexadecimal string.
func HashString(uniqueData string) string {

	// Hash the unique data using SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(uniqueData))
	hash := hasher.Sum(nil)

	// Truncate the hash to 32 bytes
	truncatedHash := hash[:HashLength]

	// Convert the truncated hash to a hexadecimal string
	return hex.EncodeToString(truncatedHash)
}

// EncryptData encrypts the input data using AES encryption with the provided password.
func EncryptData(data []byte, password string) ([]byte, error) {
	// Derive the key from the password
	key := deriveKey(password)

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher block: %v", err)
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("error generating IV: %v", err)
	}

	// Encrypt the data
	ciphertext := make([]byte, aes.BlockSize+len(data))
	copy(ciphertext[:aes.BlockSize], iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// DecryptData decrypts the input ciphertext using AES decryption with the provided password.
func DecryptData(ciphertext []byte, password string) ([]byte, error) {
	// Derive the key from the password
	key := deriveKey(password)

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher block: %v", err)
	}

	// Check if the ciphertext is shorter than the block size
	if len(ciphertext) < aes.BlockSize {
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

// deriveKey derives a key of length HashLength from the password using PBKDF2.
func deriveKey(password string) []byte {

	salt := make([]byte, 16) // take a pinch of salt...

	iterations := 4096 // set a "timer"

	// apply salt
	return pbkdf2.Key( // cook and serve as directed
		[]byte(password), // mince password to byte
		salt,             // mix in some salt
		iterations,       // cook until the timer rings
		HashLength,       // serve on a bun
		sha256.New,       // with a bed of hash
	)
}
