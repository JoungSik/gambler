package models

import (
	"encoding/json"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"gorm.io/gorm"
)

type CommandType string

type Command struct {
	Type     CommandType
	Author   string
	Argument string
}

const (
	Create   = CommandType("생성")
	Gambling = CommandType("도박")
	Bankrupt = CommandType("파산")
	Reset    = CommandType("초기화")
	All      = CommandType("올인")
)

func (c *Command) Execute(db *gorm.DB) string {
	var author discordgo.User
	json.Unmarshal([]byte(c.Author), &author)

	switch c.Type {
	case Create:
		user := User{ID: author.ID, Name: author.Username, Amount: DEFAULT_AMOUNT, InitCount: 0}
		if dividend := db.Create(&user); dividend.Error != nil {
			return "이미 생성되었어요!"
		}
		return "어서오세요! 초기 자본금은 " + strconv.Itoa(DEFAULT_AMOUNT) + "입니다!"
	case Gambling:
		dividend, err := strconv.ParseInt(c.Argument, 10, 64)
		if err != nil || dividend < 1000 {
			return "최소 배팅금은 1,000부터 입니다!"
		}

		var user User
		db.Find(&user, "id = ?", author.ID)
		if user.ID == "" {
			return "```!생성을 먼저 진행해주세요```"
		}

		origin := user.Amount
		if origin <= 0 {
			return "파산했어요.. ```!파산```으로 회복하세요!"
		} else if origin < dividend {
			return "욕심쟁이, 소지금보다 배팅금이 크면 못해요!"
		}

		num := int64(generator(origin))
		amount := dividend * num

		user.Amount += -int64(dividend) + int64(amount)
		db.Save(user)

		history := History{UserID: user.ID, Invest: dividend, Principal: origin, Result: amount, Total: user.Amount, Diameter: num}
		db.Create(&history)

		message := "유저: " + user.Name + "\n" +
			"투자금: " + humanize.Comma(int64(dividend)) + "\n" +
			"원금: " + humanize.Comma(origin) + "\n" +
			"결과: " + humanize.Comma(int64(amount)) + "\n" +
			"남은금액: " + humanize.Comma(user.Amount) + "\n" +
			"배율: " + strconv.FormatInt(num, 10)

		return message
	case Bankrupt:
		return "다음 업데이트를 기대해주세요!"
	case Reset:
		var user User
		db.Find(&user, "id = ?", author.ID)
		if user.ID == "" {
			return "```!생성을 먼저 진행해주세요```"
		}

		user.Amount = DEFAULT_AMOUNT
		user.InitCount += 1
		db.Save(user)

		return "당신은 " + humanize.Comma(int64(user.InitCount)) + "번 비겁하게 초기화 하셨습니다."
	case All:
		var user User
		db.Find(&user, "id = ?", author.ID)
		if user.ID == "" {
			return "```!생성을 먼저 진행해주세요```"
		}

		origin := user.Amount
		if origin <= 0 {
			return "파산했어요.. ```!파산```으로 회복하세요!"
		} else if origin < int64(1000) {
			return "욕심쟁이, 소지금보다 최소 배팅금이 크면 못해요!"
		}

		num := int64(generator(origin))
		amount := origin * num

		user.Amount += -int64(origin) + int64(amount)
		db.Save(user)

		history := History{UserID: user.ID, Invest: origin, Principal: origin, Result: amount, Total: user.Amount, Diameter: num}
		db.Create(&history)

		message := "유저: " + user.Name + "\n" +
			"투자금: " + humanize.Comma(int64(origin)) + "\n" +
			"원금: " + humanize.Comma(origin) + "\n" +
			"결과: " + humanize.Comma(int64(amount)) + "\n" +
			"남은금액: " + humanize.Comma(user.Amount) + "\n" +
			"배율: " + strconv.FormatInt(num, 10)

		return message
	default:
		return "없는 명령어에요"
	}
}

func Parse(message string, m *discordgo.MessageCreate) (*Command, error) {
	create, _ := regexp.Compile("!생성")
	gambling, _ := regexp.Compile("!도박")
	bankrupt, _ := regexp.Compile("!파산")
	reset, _ := regexp.Compile("!초기화")
	all, _ := regexp.Compile("!올인")

	var command *Command = nil
	data, _ := json.Marshal(m.Author)
	switch {
	case create.MatchString(message):
		command = &Command{Type: Create, Author: string(data)}
	case gambling.MatchString(message):
		regex, _ := regexp.Compile("[0-9]+$")
		_, err := strconv.Atoi(regex.FindString(m.Content))
		if err != nil {
			return nil, errors.New("금액을 입력해주세요")
		}
		command = &Command{Type: Gambling, Author: string(data), Argument: regex.FindString(m.Content)}
	case bankrupt.MatchString(message):
		command = &Command{Type: Bankrupt, Author: string(data)}
	case reset.MatchString(message):
		command = &Command{Type: Reset, Author: string(data)}
	case all.MatchString(message):
		command = &Command{Type: All, Author: string(data)}
	default:
		return nil, errors.New("unknown command type")
	}

	return command, nil
}

func generator(origin int64) int {
	var items [10]int
	rand.Seed(time.Now().UnixNano())
	min, max := 1, 10

	for i := range items {
		num := 0
		for num == 0 {
			num = rand.Intn(max)
			items[i] = num
		}
	}

	switch {
	case int(origin) >= 100000000:
		items[0] = items[0] * -1
		items[2] = items[2] * -1
		items[4] = items[4] * -1
		items[6] = items[6] * -1
		items[8] = items[8] * -1
	case int(origin) >= 1000000:
		items[0] = items[0] * -1
		items[2] = items[2] * -1
	}

	return items[rand.Intn(max-min)]
}
