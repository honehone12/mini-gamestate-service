package service

import (
	"log"
	"mini-gamestate-service/db"
	"mini-gamestate-service/server"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	RedisTimeout = time.Second
)

func parseParams() (
	string,
	string,
	string,
	db.Orm,
) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	orm, err := db.NewOrm(RedisTimeout)
	if err != nil {
		log.Fatal(err)
	}

	name := os.Getenv("SERVICE_NAME")
	if len(name) == 0 {
		log.Fatal("env param SERVICE_NAME is empty")
	}
	ver := os.Getenv("VERSION")
	if len(ver) == 0 {
		log.Fatal("env param VERSION is empty")
	}
	at := os.Getenv("SERVER_LISTEN_AT")
	if len(at) == 0 {
		log.Fatal("env param SERVER_LISTEN_AT is empty")
	}

	return name, ver, at, orm
}

func Run() {
	server.Run(parseParams())
}
