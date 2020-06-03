package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/jrcichra/discord-day-trader/common"

	"github.com/jrcichra/discord-day-trader/db"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!trade"
const balance = "balance"
const help = "help"

//Discord - handles all the discord interactions
type Discord struct {
	discord *discordgo.Session
	db      db.Database
}

//all the context a routine needs to do its job
type context struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}

func (d *Discord) getBalance(c *context, u *common.User) string {
	log.Println("We got a balance request")
	//Check if we have an account for this discord user in the database
	retMsg := ""
	for _, v := range d.db.GetAccounts(u.UserID) {
		for range d.db.GetAccountTransactions(v.AccountID) {
			retMsg += fmt.Sprintf("Doing a thing\n")
		}
	}
	return "You have lots of money"
}

//getAccount - get the account of the user
func (d *Discord) getAccount(c *context) (*common.User, error) {
	//Pull out their user ID and cast its type
	userID := common.UserID(c.Message.Author.ID)
	//Pull account data for this userID from the database
	u, err := d.db.GetUser(userID)
	return u, err
}

func (d *Discord) makeAccount(c *context) (*common.User, error) {
	var u common.User
	userID := common.UserID(c.Message.Author.ID)
	u.UserID = userID
	u.Username = c.Message.Author.Username
	_, err := d.db.CreateUser(&u)
	return &u, err
}

func (d *Discord) listAccounts(id common.UserID) string {
	msg := ""
	accounts := d.db.GetAccounts(id)
	for _, account := range accounts {
		id := account.AccountID
		name := account.Name
		created := account.Created
		msg += strconv.Itoa(int(id)) + "\t" + name + "\t" + created.Time.String() + "\n"
	}
	return msg
}

func (d *Discord) getHelp(c *context) string {
	return `
		DayTrader Discord Bot
		
	`
}

//Welcome - sequence for new users
func (d *Discord) welcome(c *context) string {
	msg := `
		Welcome to the DayTrader Discord Bot!
		We've made you an account with $100,000 dollars to start.
		Here's the output of: '!trade accounts' for your account:
	`
	msg += d.listAccounts(common.UserID(c.Message.Author.ID))
	msg += "\n\nType !trade help for a list of other useful commands!"
	return msg
}

func (d *Discord) handler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// See if the message was for us
	if strings.HasPrefix(m.Content, prefix) {
		//They want to interact with our bot
		//Prepare the context for our handler functions
		c := &context{s, m}
		//Switch case into the different commands we provide
		//remove the prefix from the rest of the logic
		request := strings.ToLower(m.Content[len(prefix):])
		var retMsg string
		//check if the first char is a space, if so, remove it
		if len(request) > 0 && request[0] == ' ' {
			request = request[1:]
		}
		log.Println("I'm processing command:'" + request + "' from: " + m.Author.Username)
		//check if the user has an account
		user, err := d.getAccount(c)
		if user.UserID == "" {
			//The don't have an account - make them one
			user, err = d.makeAccount(c)
			if err != nil {
				panic(err)
			}
			//Introduce them into the program
			retMsg = d.welcome(c)

		} else {
			switch request {
			case balance:
				retMsg = d.getBalance(c, user)
			case help:
				retMsg = d.getHelp(c)
			default:
				retMsg = "Unknown command: `" + request + "`"
			}
		}
		// Send back the constructed message
		s.ChannelMessageSend(m.ChannelID, retMsg)
	}
}

//New - creates a new discord bot
func (d *Discord) New(token string, dsn string) {
	var err error
	d.discord, err = discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	d.discord.AddHandler(d.handler)

	err = d.discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	//Make the db connection
	d.db = db.Database{}
	d.db.Connect(dsn)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	d.discord.Close()
}
