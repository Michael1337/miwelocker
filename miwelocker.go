package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run miwelocker.go <encrypt/decrypt> <file> <password> [ID] [extension]")
		return
	}

	action := os.Args[1]
	filePath := os.Args[2]
	password := os.Args[3]

	switch action {
	case "encrypt":
		if len(os.Args) < 6 {
			fmt.Println("Usage for encryption: go run miwelocker.go encrypt <file> <password> <ID> <extension>")
			return
		}
		id := os.Args[4]
		extension := os.Args[5]
		err := encryptFile(filePath, password, id, extension)
		if err != nil {
			fmt.Println("Error encrypting file:", err)
		} else {
			fmt.Println("File encrypted successfully.")
		}
	case "decrypt":
		err := decryptFile(filePath, password)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
		} else {
			fmt.Println("File decrypted successfully.")
		}
	default:
		fmt.Println("Unknown action. Use 'encrypt' or 'decrypt'.")
	}
}

func encryptFile(filePath, password, id, extension string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(deriveKey(password))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)
	newFilePath := fmt.Sprintf("%s.%s.%s", filePath, id, extension)
	return ioutil.WriteFile(newFilePath, encryptedData, 0644)
}

func decryptFile(filePath, password string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(deriveKey(password))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plainData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	originalFilePath := removeIDAndExtension(filePath)
	return ioutil.WriteFile(originalFilePath, plainData, 0644)
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func removeIDAndExtension(filePath string) string {
	lastDot := strings.LastIndex(filePath, ".")
	if lastDot == -1 {
		return filePath
	}
	beforeLastDot := strings.LastIndex(filePath[:lastDot], ".")
	if beforeLastDot == -1 {
		return filePath[:lastDot]
	}
	return filePath[:beforeLastDot]
}
