# Boo

Boo is a telegram bot written in Go. It uses [telegram-bot-api](github.com/go-telegram-bot-api/telegram-bot-api) V5.

It might be useful if you're looking for golang realisation of telegram bot with following features:
- [x] 

## Install & run
```
git clone github.com/errorboi/boo
cd boo
cp .env.sample .env
```
Specify ENVs in .env file. Bare minimum is BOT_TOKEN and DB related variables.
When you're done, run:
```
docker-compose up --build
```



## Tech stuff

### Create migration
```
migrate create -ext sql -dir store/postgres/migrations -seq {migration_name}
```