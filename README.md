```
    ____                        __    __
   / __ \____ _______________ _/ /_  / /___ _____  _________ _
  / /_/ / __  / ___/ ___/ __  / __ \/ / __  / __ \/ ___/ __  /
 / ____/ /_/ (__  |__  ) /_/ / /_/ / / /_/ / / / / /__/ /_/ /
/_/    \__,_/____/____/\__,_/_.___/_/\__,_/_/ /_/\___/\__,_/
```

Passablanca is a command-line based password manager.

It uses 256-bit AES-GCM to encrypt the password list, and SHA-256 to derive the key.

### Installing

`
go get -u -v github.com/superp00t/passablanca
`

### Running

Assuming `$GOPATH/bin` is in your `PATH`, just type `passablanca`.