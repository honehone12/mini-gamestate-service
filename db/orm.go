package db

import (
	"errors"
	"log"
	"mini-gamestate-service/db/controller"
	"os"
	"time"

	"github.com/redis/rueidis"
)

type Orm struct {
	kv      rueidis.Client
	timeOut time.Duration
}

func NewOrm(timeOut time.Duration) (Orm, error) {
	addr := os.Getenv("REDIS_ADDR")
	if len(addr) == 0 {
		return Orm{nil, 0}, errors.New("env param REDIS_ADDR is empty")
	}

	kv, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{addr},
	})
	if err != nil {
		return Orm{nil, 0}, err
	}

	log.Println("established new kv connection")

	return Orm{
		kv:      kv,
		timeOut: timeOut,
	}, nil
}

func (orm Orm) Session() controller.SessionController {
	return controller.NewSessionController(orm.kv, orm.timeOut)
}

func (orm Orm) Jewel() controller.JewelController {
	return controller.NewJewelController(orm.kv, orm.timeOut)
}
