package ldapcon

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Credential struct {
	Username string
	Password string
}

var credential *Credential

func NewCredential() *Credential {
	if credential == nil {
		credential = &Credential{}
		var username string
		fmt.Print("Enter LDAP Username: ")
		fmt.Scanln(&username)

		fmt.Print("Enter Password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\nError reading password")
			return nil
		}
		password := string(passwordBytes)

		credential.Username = username
		credential.Password = password
		fmt.Println("")
	}
	return credential

}
