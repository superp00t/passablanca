package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/howeyc/gopass"
	"github.com/olekukonko/tablewriter"
	"github.com/superp00t/passablanca/cryptutil"
	readline "github.com/superp00t/readline"
)

const helpMessage = `To add an account:
	register <website URL> <account name> <password>

To list passwords:
	ls <optional URL match regex>

To copy a password to clipboard:
	clip <numeric ID>
    
To display this message:
	help

To generate a random password:
	generate
	
To quit Passablanca: 
	quit
`
const hdr = `    ____                        __    __                      
   / __ \____ _______________ _/ /_  / /___ _____  _________ _
  / /_/ / __  / ___/ ___/ __  / __ \/ / __  / __ \/ ___/ __  /
 / ____/ /_/ (__  |__  ) /_/ / /_/ / / /_/ / / / / /__/ /_/ / 
/_/    \__,_/____/____/\__,_/_.___/_/\__,_/_/ /_/\___/\__,_/  
`

type AccountEntry struct {
	URL      string
	Username string
	Password string
}

type Database struct {
	Accounts map[int]*AccountEntry
}

var dblocation string

func main() {
	usr, _ := user.Current()
	dblocation = path.Join(usr.HomeDir, ".passablanca_store")

	if _, err := os.Stat(dblocation); os.IsNotExist(err) {
		fmt.Println(hdr)
		fmt.Println("This is your first time using Passablanca.")
		fmt.Printf("Enter in a strong master password (password will not echo). ")

		passwd, err := gopass.GetPasswd()
		if err != nil {
			panic(err)
		}
		password := string(passwd)

		fmt.Println("Your password is " + password)

		mp := make(map[int]*AccountEntry)
		db := Database{
			Accounts: mp,
		}
		WriteDatabase(password, db)

		fmt.Println("Password list has been created at " + dblocation)
		fmt.Println("Type 'passablanca' again to get started.")
	} else {
		fmt.Println("Welcome back to Passablanca!")
		fmt.Printf("Please enter your master password (password will not echo). ")

		passwd, err := gopass.GetPasswd()
		password := string(passwd)
		if err != nil {
			panic(err)
		}

		db := ReadDatabase(password)

		// Begin readline loop
		rl, err := readline.NewEx(&readline.Config{
			UniqueEditLine: false,
		})
		if err != nil {
			panic(err)
		}

		for {
			rl.SetPrompt("[Passablanca] >>> ")
			ln := rl.Line()

			if ln.CanContinue() {
				continue
			} else if ln.CanBreak() {
				break
			}

			args := strings.Split(ln.Line, " ")

			switch args[0] {
			case "generate":
				fmt.Println(cryptutil.RandomString())
			case "register":
				if len(args) == 4 {
					ae := &AccountEntry{
						URL:      args[1],
						Username: args[2],
						Password: args[3],
					}

					db.Accounts[MaxInt(db.Accounts)+1] = ae

					WriteDatabase(password, db)
				} else {
					fmt.Println("Usage: register <website URL> <account name> <password>")
				}
			case "ls":
				Headers := []string{"ID", "URL", "Username", "Password"}
				var Body [][]string
				var rgx *regexp.Regexp
				if len(args) == 2 {
					rgx, err = regexp.Compile(args[1])
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
				}

				for k, v := range db.Accounts {
					if len(args) == 2 {
						mt := rgx.MatchString(v.URL)
						if mt {
							Body = append(Body, []string{
								fmt.Sprintf("%d", k),
								v.URL,
								v.Username,
								v.Password,
							})
						}
					} else {
						Body = append(Body, []string{
							fmt.Sprintf("%d", k),
							v.URL,
							v.Username,
							v.Password,
						})
					}
				}

				table := tablewriter.NewWriter(os.Stdout)
				for _, v := range Body {
					table.Append(v)
				}
				table.SetHeader(Headers)
				table.Render()
			case "clip":
				if len(args) == 2 {
					index, err := strconv.Atoi(args[1])
					if err != nil {
						fmt.Println(args[1] + " is not a valid ID. (" + err.Error() + ")")
					} else {
						if db.Accounts[index] == nil {
							fmt.Println(args[1] + " account entry does not exist.")
						} else {
							if err := clipboard.WriteAll(db.Accounts[index].Password); err != nil {
								fmt.Println("Clipboard paste failed :(")
							} else {
								fmt.Println("Password copied to clipboard!")
							}
						}
					}
				} else {
					fmt.Println("Usage: clip <numeric id>")
				}
			case "inspect":
				dat, err := json.MarshalIndent(db, "", "    ")
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dat))
			case "help":
				fmt.Println(helpMessage)
			case "quit":
				fmt.Println("Goodbye!")
				os.Exit(0)
			default:
				fmt.Println(args[0] + " is not a valid command.")
				fmt.Println(helpMessage)
			}
		}
	}
}

func ReadDatabase(password string) Database {
	encryptedDbData, err := ioutil.ReadFile(dblocation)
	if err != nil {
		panic(err)
	}

	decryptedDbData, err := cryptutil.Decrypt(password, encryptedDbData)
	if err != nil {
		fmt.Println("Could not decode the ~/.passablanca_store file. Maybe you entered in the wrong password?")
		os.Exit(-1)
	}

	var db Database

	err = json.Unmarshal(decryptedDbData, &db)
	if err != nil {
		fmt.Println("Could not decode the ~/.passablanca_store file. Maybe you entered in the wrong password?")
		os.Exit(-1)
	}
	return db
}

func WriteDatabase(password string, db Database) {
	dat, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		panic(err)
	}

	encryptedDbData := cryptutil.Encrypt(password, dat)
	ioutil.WriteFile(dblocation, encryptedDbData, 0755)
}

func MaxInt(m map[int]*AccountEntry) int {
	ret := 0
	for k, _ := range m {
		if k > ret {
			ret = k
		}
	}
	return ret
}
