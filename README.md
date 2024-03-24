# HUBLISH Backend (Go + Fiber + Gorm + PostgreSQL)

## Descriptions 📋

This is a sub repository for [Hublish](https://github.com/mwongsatorn/hublish). The backend of hublish built using

- Go
- Fiber
- Gorm
- PostgreSQL
- Docker

## Requirements 🛠

- Go (version 1.21.4 +)
- Docker

## Get Started 🏃

1. Run this command to clone this repository

```bash
git clone https://github.com/mwongsatorn/hublish-be-go
```

2. Go to the project directory
3. Run this command to install all of the packages

```go
go mod download
```

4. Run this command to run the database containers

```bash
docker compose up -d
```

5. Run this command to populate data inside the database (You can skip this step if you are not running the project first time or you don't want to )

```
# If you have make install on your local machine

make seed

# If you do not have make install on your local machine

go run cmd/api/main.go -seed
```

6. Run this command to start the server

```
# If you have make install on your local machine

make run

# If you do not have make install on your local machine

go run cmd/api/main.go
```
