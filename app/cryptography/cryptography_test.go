package cryptography_test

import (
	"bytes"
	"testing"

	"github.com/secretnamebasis/secret-site/app/cryptography"
)

func TestEncryptDecrypt(t *testing.T) {
	// Test data
	data := []byte("hello, world!")
	password := "secretPassword"

	// Encrypt the data
	encryptedData, err := cryptography.EncryptData(data, password)
	if err != nil {
		t.Errorf("Error encrypting data: %v", err)
	}

	// Decrypt the encrypted data
	decryptedData, err := cryptography.DecryptData(encryptedData, password)
	if err != nil {
		t.Errorf("Error decrypting data: %v", err)
	}

	// Verify that decrypted data matches the original data
	if !bytes.Equal(decryptedData, data) {
		t.Errorf("Decrypted data does not match original data. Expected: %s, Got: %s", string(data), string(decryptedData))
	}
}

func TestDecryptInvalidData(t *testing.T) {
	// Test data
	invalidData := []byte("secret")
	password := "secretPassword"

	// Attempt to decrypt invalid data
	_, err := cryptography.DecryptData(invalidData, password)

	// Verify that the expected error is returned
	if err == nil || err.Error() != "ciphertext too short" {
		t.Errorf("Expected 'ciphertext too short' error, but got: %v", err)
	}
}

func TestDecryptWithIncorrectPassword(t *testing.T) {
	// Test data
	data := []byte("hello, world! I am secret.")
	password := "secretPassword"
	incorrectPassword := "incorrectPassword"

	// Encrypt the data with the correct password
	encryptedData, _ := cryptography.EncryptData(data, password)

	// Attempt to decrypt with an incorrect password
	decryptedData, err := cryptography.DecryptData(encryptedData, incorrectPassword)

	// Verify when an error is returned
	if err != nil {
		t.Errorf("Expected no error when decrypting with an incorrect password: %s", err)
	}

	// Ensure that decrypted data does not match the original data
	if bytes.Equal(decryptedData, data) {
		t.Errorf("Decrypted data matches original data. Expected mismatch for incorrect password. Got: %s", string(decryptedData))
	}
}

func TestEncryptDecryptEmptyData(t *testing.T) {
	// Test data
	var emptyData []byte
	password := "secretPassword"

	// Encrypt the empty data
	encryptedData, err := cryptography.EncryptData(emptyData, password)
	if err != nil {
		t.Errorf("Error encrypting empty data: %v", err)
	}

	// Decrypt the encrypted empty data
	decryptedData, err := cryptography.DecryptData(encryptedData, password)
	if err != nil {
		t.Errorf("Error decrypting empty data: %v", err)
	}

	// Verify that decrypted empty data matches the original data
	if !bytes.Equal(decryptedData, emptyData) {
		t.Errorf("Decrypted empty data does not match original data. Expected: %v, Got: %v", emptyData, decryptedData)
	}
}
