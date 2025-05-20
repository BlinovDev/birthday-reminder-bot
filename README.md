# birthday-reminder-bot

This project represents Telegram bot for storing people birthdays and notify about them.

## TODO list

- [x] Add buttons for main commands;
- [x] Develop birthday adding process;
- [x] Configure reminder task that can be initiated as cron;
- [x] Add update and delete functionality;
- [ ] Add `/help` response;

## Deploy

Before deploy you should have some server to run your bot on and in following example server name in local configs is `birthday-reminder-bot`.

```bash
# compile progect for Linux
GOOS=linux GOARCH=amd64 go build -o birthday-reminder-bot-v1.0 main.go

# upload executable file to server
scp birthday-reminder-bot-v1.0 birthday-reminder-bot:"/home/admin"

# open server' bash console
ssh birthday-reminder-bot

# go to folder with executable file
cd /home/admin

# start screen
screen -S birthday-reminder-bot

# run bot 
./birthday-reminder-bot-v1.0

# kill screen to stop the bot
screen -X -S birthday-reminder-bot kill
```
