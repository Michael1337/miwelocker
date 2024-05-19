# MiweLocker

MiweLocker is a simple file encryption and decryption tool written in Go.

## Usage

### Encryption

To encrypt a file, run the following command:

```sh
go run miwelocker.go encrypt <file> <password> <ID> <extension>
```

- `<file>`: The path to the file you want to encrypt.
- `<password>`: The password used for encryption.
- `<ID>`: An identifier for the encrypted file.
- `<extension>`: The custom extension for the encrypted file.

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
