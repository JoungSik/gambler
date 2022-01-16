package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/JoungSik/gambler/cmd/models"
	"github.com/JoungSik/gambler/configs"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

func main() {
	config := configs.NewConfig()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	db, err := configs.InitDB(true)
	if err != nil {
		fmt.Println("error creating Database session,", err)
		return
	}

	if m.Content == "!생성" {
		user := models.User{ID: m.Author.ID, Server: m.GuildID, Name: m.Author.Username, Amount: models.DEFAULT_AMOUNT, InitCount: 0}
		if result := db.Create(&user); result.Error != nil {
			s.ChannelMessageSend(m.ChannelID, "이미 생성되었어요!")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "어서오세요! 초기 자본금은 "+strconv.Itoa(models.DEFAULT_AMOUNT)+"입니다!")
	}

	if metched, _ := regexp.MatchString("!도박", m.Content); metched {
		regex, _ := regexp.Compile("[0-9]+$")
		result, err := strconv.Atoi(regex.FindString(m.Content))
		if err != nil {
			return
		}

		if result < 1000 {
			s.ChannelMessageSend(m.ChannelID, "최소 배팅금은 1,000부터 입니다!")
			return
		}

		amount := result * generator()

		var user models.User
		db.Find(&user, "id = ? AND server = ?", m.Author.ID, m.GuildID)
		origin := user.Amount

		if origin <= 0 {
			s.ChannelMessageSend(m.ChannelID, "파산했어요.. ```!파산```으로 회복하세요!")
			return
		}

		if origin < int64(result) {
			s.ChannelMessageSend(m.ChannelID, "욕심쟁이, 소지금보다 배팅금이 크면 못해요!")
			return
		}

		user.Amount += -int64(result) + int64(amount)
		db.Save(user)

		message := "투자금: " + humanize.Comma(int64(result)) + "\n원금: " + humanize.Comma(origin) + "\n결과: " + humanize.Comma(int64(amount)) + "\n남은 금액: " + humanize.Comma(user.Amount)
		s.ChannelMessageSend(m.ChannelID, message)
	}

	if m.Content == "!파산" {
		s.ChannelMessageSend(m.ChannelID, "다음 업데이트를 기대해주세요!")
	}

	if m.Content == "!초기화" {
		var user models.User
		db.Find(&user, "id = ? AND server = ?", m.Author.ID, m.GuildID)
		user.Amount = models.DEFAULT_AMOUNT
		user.InitCount += 1
		db.Save(user)

		s.ChannelMessageSend(m.ChannelID, "당신은 "+humanize.Comma(int64(user.InitCount))+"번 비겁하게 초기화 하셨습니다.")
	}

}

func generator() int {
	rand.Seed(time.Now().UnixNano())
	min := -10
	max := 10
	return rand.Intn(max-min+1) + min
}
