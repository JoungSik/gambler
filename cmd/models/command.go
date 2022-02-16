package models

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"gorm.io/gorm"
)

type CommandType string

type Command struct {
	Type     CommandType
	Author   string
	Server   string
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

	user, _ := strconv.ParseInt(author.ID, 10, 64)
	server, _ := strconv.ParseInt(c.Server, 10, 64)

	switch c.Type {
	case Create:
		client := http.Client{}

		form := url.Values{}
		form.Add("user[id]", author.ID)
		form.Add("user[name]", author.Username)
		form.Add("user[email]", author.ID+"@example.com")
		form.Add("user[password]", author.Username)

		req, _ := http.NewRequest("POST", "http://localhost:3000/users", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		newServer := Server{ID: server}
		db.Select("ID").Create(&newServer)

		newAccount := Account{UserId: user, ServerId: server, Amount: DEFAULT_AMOUNT, InitCount: 0}
		if resultAccount := db.First(&newAccount); resultAccount.RowsAffected > 0 {
			return "이미 생성되었어요!"
		}

		db.Create(&newAccount)
		return "어서오세요! 초기 자본금은 " + strconv.Itoa(DEFAULT_AMOUNT) + "입니다!"
	case Gambling:
		dividend, err := strconv.ParseInt(c.Argument, 10, 64)
		if err != nil || dividend < 1000 {
			return "최소 배팅금은 1,000부터 입니다!"
		}

		var account Account
		result := db.Find(&account, "user_id = ? and server_id = ?", user, server)
		if result.Error != nil {
			return "```!생성을 먼저 진행해주세요```"
		}

		origin := account.Amount
		if origin <= 0 {
			return "파산했어요.. ```!파산```으로 회복하세요!"
		} else if origin < dividend {
			return "욕심쟁이, 소지금보다 배팅금이 크면 못해요!"
		} else if dividend < int64(float64(origin)*0.09) {
			return "쫄?, 최소 금액은 " + humanize.Comma(int64(float64(origin)*0.09)) + " 입니다."
		}

		num := int64(generator(origin))
		amount := dividend * num
		tax := payTax(amount)

		account.Amount += -int64(dividend) + int64(amount) - tax

		var bonus int64
		if account.Amount < 0 {
			bonus = payTax(origin * (num * -1))
			account.Amount += bonus
		}

		db.Save(user)

		history := History{UserId: user, ServerId: server, Invest: dividend, Principal: origin, Result: amount, Tax: tax, Total: account.Amount, Diameter: num}
		db.Create(&history)

		var message string
		if account.Amount > 0 {
			message = "유저: " + author.Username + "\n" +
				"투자금: " + humanize.Comma(int64(dividend)) + "\n" +
				"원금: " + humanize.Comma(origin) + "\n" +
				"결과: " + humanize.Comma(int64(amount)) + "\n" +
				"세금: " + humanize.Comma(tax) + "\n" +
				"남은금액: " + humanize.Comma(account.Amount) + "\n" +
				"배율: " + strconv.FormatInt(num, 10)
		} else {
			message = "유저: " + author.Username + "\n" +
				"투자금: " + humanize.Comma(int64(origin)) + "\n" +
				"원금: " + humanize.Comma(origin) + "\n" +
				"결과: " + humanize.Comma(int64(amount)) + "\n" +
				"세금: " + humanize.Comma(tax) + "\n" +
				"위로금: " + humanize.Comma(bonus) + "\n" +
				"남은금액: " + humanize.Comma(account.Amount) + "\n" +
				"배율: " + strconv.FormatInt(num, 10) + "\n\n" +
				"ㅋ 감사합니다. 고객님, 설거지는 저쪽입니다."
		}

		return message
	case Bankrupt:
		return "다음 업데이트를 기대해주세요!"
	case Reset:
		var account Account
		result := db.Find(&account, "user_id = ? and server_id = ?", user, server)
		if result.Error != nil {
			return "```!생성을 먼저 진행해주세요```"
		}

		account.Amount = DEFAULT_AMOUNT
		account.InitCount += 1
		db.Save(account)

		return "당신은 " + humanize.Comma(int64(account.InitCount)) + "번 비겁하게 초기화 하셨습니다."
	case All:
		var account Account
		result := db.Find(&account, "user_id = ? and server_id = ?", user, server)
		if result.Error != nil {
			return "```!생성을 먼저 진행해주세요```"
		}

		origin := account.Amount
		if origin <= 0 {
			return "파산했어요.. ```!파산```으로 회복하세요!"
		} else if origin < int64(1000) {
			return "욕심쟁이, 소지금보다 최소 배팅금이 크면 못해요!"
		}

		num := int64(generator(origin))
		amount := origin * num
		tax := payTax(amount)

		account.Amount += -int64(origin) + int64(amount) - tax

		var bonus int64
		if account.Amount < 0 {
			bonus = payTax(origin * (num * -1))
			account.Amount += bonus
		}

		db.Save(account)

		history := History{UserId: account.UserId, ServerId: account.ServerId, Invest: origin, Principal: origin, Result: amount, Tax: tax, Total: account.Amount, Diameter: num}
		db.Create(&history)

		var message string
		if account.Amount > 0 {
			message = "유저: " + author.Username + "\n" +
				"투자금: " + humanize.Comma(int64(origin)) + "\n" +
				"원금: " + humanize.Comma(origin) + "\n" +
				"결과: " + humanize.Comma(int64(amount)) + "\n" +
				"세금: " + humanize.Comma(tax) + "\n" +
				"남은금액: " + humanize.Comma(account.Amount) + "\n" +
				"배율: " + strconv.FormatInt(num, 10)
		} else {
			message = "유저: " + author.Username + "\n" +
				"투자금: " + humanize.Comma(int64(origin)) + "\n" +
				"원금: " + humanize.Comma(origin) + "\n" +
				"결과: " + humanize.Comma(int64(amount)) + "\n" +
				"세금: " + humanize.Comma(tax) + "\n" +
				"위로금: " + humanize.Comma(bonus) + "\n" +
				"남은금액: " + humanize.Comma(account.Amount) + "\n" +
				"배율: " + strconv.FormatInt(num, 10) + "\n\n" +
				"ㅋ 감사합니다. 고객님, 설거지는 저쪽입니다."
		}

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
		command = &Command{Type: Create, Author: string(data), Server: m.GuildID}
	case gambling.MatchString(message):
		regex, _ := regexp.Compile("[0-9]+$")
		_, err := strconv.Atoi(regex.FindString(m.Content))
		if err != nil {
			return nil, errors.New("금액을 입력해주세요")
		}
		command = &Command{Type: Gambling, Author: string(data), Server: m.GuildID, Argument: regex.FindString(m.Content)}
	case bankrupt.MatchString(message):
		command = &Command{Type: Bankrupt, Author: string(data), Server: m.GuildID}
	case reset.MatchString(message):
		command = &Command{Type: Reset, Author: string(data), Server: m.GuildID}
	case all.MatchString(message):
		command = &Command{Type: All, Author: string(data), Server: m.GuildID}
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

func payTax(amount int64) int64 {
	switch {
	case amount > 0 && amount <= 12000000:
		return (int64((float64(amount) * 0.06)))
	case amount > 12000000 && amount <= 46000000:
		return (int64((float64(amount) * 0.15)) - 1080000)
	case amount > 46000000 && amount <= 88000000:
		return (int64((float64(amount) * 0.24)) - 5220000)
	case amount > 88000000 && amount <= 150000000:
		return (int64((float64(amount) * 0.35)) - 14900000)
	case amount > 150000000 && amount <= 300000000:
		return (int64((float64(amount) * 0.38)) - 19400000)
	case amount > 300000000 && amount <= 500000000:
		return (int64((float64(amount) * 0.4)) - 25400000)
	case amount > 500000000:
		return (int64((float64(amount) * 0.42)) - 35400000)
	default:
		return 0
	}
}
