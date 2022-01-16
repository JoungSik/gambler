.PHONY: init clean build run

init: 
	go mod download github.com/bwmarrin/discordgo v0.23.2
	go mod download github.com/dustin/go-humanize v1.0.0
	go mod download github.com/joho/godotenv v1.4.0
	go mod download gorm.io/driver/mysql v1.2.3
	go mod download gorm.io/gorm v1.22.5
	go mod tidy

clean: 
	rm -rf ./bin go.sum
	
build:
	export GO111MODULE=on
	go build -o bin/main cmd/main.go 

run: clean init build