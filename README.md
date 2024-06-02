# MiweLocker

MiweLocker is a simple file encryption and decryption tool written in Go.

## Usage

### Setup / Initializing Go Modules

Before running the tool, you need to initialize the Go modules and download the required packages.

1. Navigate to the directory containing `miwelocker.go`.
2. Initialize Go modules:

```sh
go mod init miwelocker
go get golang.org/x/crypto/pbkdf2
```

### Encryption

To encrypt a file, run the following command:

```sh
go run miwelocker.go encrypt <file> <password> <ID> [extension]
```

- `<file>`: The path to the file you want to encrypt.
- `<password>`: The password used for encryption.
- `<ID>`: An identifier for the encrypted file.
- `[extension]`: The custom extension for the encrypted file (optional). If not provided, the default extension `miwelocked` will be used.

Example:

```sh
go run miwelocker.go encrypt example.txt mypassword 12345 enc
```

This will produce an encrypted file named `example.txt.12345.enc`.

### Decryption

To decrypt a file, run the following command:

```sh
go run miwelocker.go decrypt <file> <password>
```

- `<file>`: The path to the encrypted file you want to decrypt.
- `<password>`: The password used for decryption.

Example:

```sh
go run miwelocker.go decrypt example.txt.12345.enc mypassword
```

This will decrypt the file `example.txt.12345.enc` and restore it to its original state.

## Code Explanation

### Main Function

The `main` function checks if the required arguments are provided and then calls the appropriate function (`encryptFile` or `decryptFile`) based on the action (`encrypt` or `decrypt`).

### Encryption

The `encryptFile` function performs the following steps:

1. Reads the contents of the file.
2. Generates a salt.
3. Derives a key from the password and salt using PBKDF2.
4. Encrypts the data using AES-GCM with the derived key.
5. Appends the salt to the encrypted data and writes it to a new file with the specified ID and extension.

### Decryption

The `decryptFile` function performs the following steps:

1. Reads the encrypted data from the file.
2. Extracts the salt from the encrypted data.
3. Derives the key from the password and salt using PBKDF2.
4. Decrypts the data using AES-GCM with the derived key.
5. Writes the decrypted data to a new file without the ID and extension.

### Key Derivation

The `deriveKey` function uses PBKDF2 (Password-Based Key Derivation Function 2) to derive a key from the password and salt. This process makes it computationally expensive for attackers to brute-force the password.

### File Extension

The default extension for encrypted files is `.miwelocked`, but you can specify a custom extension when encrypting a file.
