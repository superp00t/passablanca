# Passablanca

Passablanca is a command-line based password manager.

It uses 256-bit AES-GCM to encrypt the password list, and SHA-256 to derive the key.

### Installing

`
go get -u -v github.com/superp00t/passablanca
`

### Running

Assuming `$GOPATH/bin` is in your `PATH`, just type `passablanca`.