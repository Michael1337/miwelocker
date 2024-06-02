package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	saltSize       = 16
	keySize        = 32
	iterations     = 100000
	defaultExt     = "miwelocked"
	usageMessage   = "Usage: go run miwelocker.go <encrypt/decrypt> <file> <password> [ID] [extension]"
)

func main() {
	if len(os.Args) < 4 || len(os.Args) > 6 {
		fmt.Println(usageMessage)
		return
	}

	action := os.Args[1]
	filePath := os.Args[2]
	password := os.Args[3]

	switch action {
	case "encrypt":
		if len(os.Args) < 5 {
			fmt.Println("Usage for encryption: go run miwelocker.go encrypt <file> <password> <ID> [extension]")
			return
		}
		id := os.Args[4]
		extension := defaultExt
		if len(os.Args) == 6 {
			extension = os.Args[5]
			extension = strings.ReplaceAll(extension, ".", "") // remove additional dots in custom extensions for func removeIDAndExtension to work correctly
		}
		err := encryptFile(filePath, password, id, extension)
		if err != nil {
			fmt.Println("Error encrypting file:", err)
		} else {
			fmt.Println("File encrypted successfully.")
		}
	case "decrypt":
		if len(os.Args) != 4 {
			fmt.Println("Usage for decryption: go run miwelocker.go decrypt <file> <password>")
			return
		}
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

	salt := make([]byte, saltSize)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
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
	encryptedData = append(salt, encryptedData...)

	newFilePath := fmt.Sprintf("%s.%s.%s", filePath, id, extension)
	return ioutil.WriteFile(newFilePath, encryptedData, 0600)
}

func decryptFile(filePath, password string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(data) < saltSize {
		return fmt.Errorf("ciphertext too short")
	}

	salt, data := data[:saltSize], data[saltSize:]
	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
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
	return ioutil.WriteFile(originalFilePath, plainData, 0600)
}

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

func removeIDAndExtension(filePath string) string {
	lastDot := strings.LastIndex(filePath, ".")
	if strings.Count(filePath, ".") == 1 || lastDot == -1 {
		return filePath // if string has no or just one dot, i. e. if someone removed the ID and MiweLocker extension
	}
	beforeLastDot := strings.LastIndex(filePath[:lastDot], ".")
	if beforeLastDot == -1 {
		return filePath[:lastDot] // if someone removed the MiweLocker extension, but not the ID
	}
	return filePath[:beforeLastDot]
}
