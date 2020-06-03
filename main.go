package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/jrcichra/discord-day-trader/discord"
)

const credentials = "credentials"

func getCredentials() (string, string) {
	token := ""
	dsn := ""

	file, err := os.Open(credentials)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		if i == 0 {
			token = scanner.Text()
		} else if i == 1 {
			dsn = scanner.Text()
		} else {
			log.Fatal("Unexpected line in credentials file - it should be 1=token,2=dsn")
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return token, dsn
}

func main() {
	token := flag.String("token", "", "Token for your bot that handles the stock requests")
	dsn := flag.String("dsn", "", "Connection string to the database")
	flag.Parse()
	if *token == "" || *dsn == "" {
		//If we got a blank flag, check for a credentials file
		if _, err := os.Stat(credentials); err == nil || os.IsExist(err) {
			//There is a credentials file, parse it
			t, d := getCredentials()
			*token = t
			*dsn = d
		}
	}
	var d discord.Discord
	log.Println(*token, *dsn)
	d.New(*token, *dsn)
}
