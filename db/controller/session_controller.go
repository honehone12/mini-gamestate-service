package controller

import (
	"context"
	"errors"
	"mini-gamestate-service/db/models"
	"strconv"
	"time"

	"github.com/redis/rueidis"
)

const (
	sessionKeyExtension = ":session"
	oneTimeIdField      = "one_time_id"
	createdAtField      = "created_at_unix"
)

var (
	ErrorFieldNotFound = errors.New("could not find the field")
)

type SessionController struct {
	kv      rueidis.Client
	timeOut time.Duration
}

func NewSessionController(kv rueidis.Client, timeOut time.Duration) SessionController {
	return SessionController{
		kv:      kv,
		timeOut: timeOut,
	}
}

func (c SessionController) Set(userUuid string, oneTimeSessionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	key := userUuid + sessionKeyExtension
	now := strconv.FormatInt(time.Now().Unix(), 10)
	set := c.kv.B().Hset().Key(key).FieldValue().
		FieldValue(oneTimeIdField, oneTimeSessionId).
		FieldValue(createdAtField, now).
		Build()
	res := c.kv.Do(ctx, set)
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

func (c SessionController) Get(userUuid string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	key := userUuid + sessionKeyExtension
	get := c.kv.B().Hgetall().Key(key).Build()
	res := c.kv.Do(ctx, get)
	if res.Error() != nil {
		return nil, res.Error()
	}
	m, err := res.AsStrMap()
	if err != nil {
		return nil, err
	}

	oneTimeId, ok := m[oneTimeIdField]
	if !ok {
		return nil, ErrorFieldNotFound
	}
	createdAtStr, ok := m[createdAtField]
	if !ok {
		return nil, ErrorFieldNotFound
	}
	createdAt, err := strconv.ParseInt(createdAtStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		OneTimeId:     oneTimeId,
		CreatedAtUnix: createdAt,
	}, nil
}
