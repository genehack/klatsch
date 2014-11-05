package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func initStuff(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(getDatabaseFile()); err == nil {
		fmt.Println("A Klatsch database already exists. Not re-init-ing.")
		fmt.Println("Delete the database and rerun the command if you wish to proceed.")
		return
	}

	secrets, err := readSecrets()
	if err != nil {
		log.Fatal(err)
	}

	err = writeSecrets(secrets)
	if err != nil {
		log.Fatal(err)
	}

}

func readOne(prompt string) (value string, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)
	value, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	value = strings.TrimSpace(value)
	if value == "" {
		fmt.Println("ERROR: you must enter a value!\n\n")
		return readOne(prompt)
	}

	return
}

func readSecrets() (secrets map[string]string, err error) {
	secrets = make(map[string]string, 4)

	secrets["consumerKey"], err = readOne("Enter your application (consumer) key: ")
	if err != nil {
		return
	}

	secrets["consumerSecret"], err = readOne("Enter your application (consumer) secret: ")
	if err != nil {
		return
	}

	// TODO move this into the process of adding an account
	secrets["accessToken"], err = readOne("Enter your account access token: ")
	if err != nil {
		return
	}

	// TODO move this into the process of adding an account
	secrets["accessSecret"], err = readOne("Enter your account access secret: ")
	if err != nil {
		return
	}

	return
}

func writeSecrets(secrets map[string]string) error {
	db := getDatabaseHandle()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.Prepare("INSERT INTO config (key,val) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	for key, val := range secrets {
		_, err = statement.Exec(key, val)
		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}
