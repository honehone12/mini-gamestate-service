package controller

import (
	"context"
	"errors"
	"mini-gamestate-service/db/models"
	"mini-gamestate-service/db/models/jewels"
	"strconv"
	"time"

	"github.com/redis/rueidis"
)

const (
	jewelKeyExtension = ":jewels"
)

var (
	ErrorValueNotSignedInteger = errors.New("the value is not expected to be smaller than 0")
)

type JewelController struct {
	kv      rueidis.Client
	timeOut time.Duration
}

func NewJewelController(kv rueidis.Client, timeOut time.Duration) JewelController {
	return JewelController{
		kv:      kv,
		timeOut: timeOut,
	}
}

func (c JewelController) IncrBy(
	userUuid string,
	color jewels.ColorCode,
	incr int64,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	n, err := c.Get(userUuid, color)
	if err != nil {
		return err
	}
	if incr < 0 && n < incr {
		return ErrorValueNotSignedInteger
	}

	field, err := color.ColorCodeToString()
	if err != nil {
		return err
	}
	key := userUuid + jewelKeyExtension
	incrBy := c.kv.B().Hincrby().Key(key).
		Field(field).Increment(incr).Build()
	res := c.kv.Do(ctx, incrBy)
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

func (c JewelController) Get(userUuid string, color jewels.ColorCode) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	field, err := color.ColorCodeToString()
	if err != nil {
		return 0, err
	}
	key := userUuid + jewelKeyExtension
	get := c.kv.B().Hget().Key(key).Field(field).Build()
	res := c.kv.Do(ctx, get)
	if res.Error() != nil {
		return 0, res.Error()
	}
	nInt, err := res.AsInt64()
	if err != nil {
		return 0, err
	}
	return nInt, nil
}

func (c JewelController) SetAll(userUuid string, j *models.Jewel) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	key := userUuid + jewelKeyExtension
	base := 10
	setAll := c.kv.B().Hmset().Key(key).FieldValue().
		FieldValue(jewels.RedField, strconv.FormatInt(j.Red, base)).
		FieldValue(jewels.BlueField, strconv.FormatInt(j.Blue, base)).
		FieldValue(jewels.GreenField, strconv.FormatInt(j.Green, base)).
		FieldValue(jewels.YellowField, strconv.FormatInt(j.Yellow, base)).
		FieldValue(jewels.BlackField, strconv.FormatInt(j.Black, base)).
		Build()
	res := c.kv.Do(ctx, setAll)
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

func (c JewelController) GetAll(userUuid string) (*models.Jewel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	key := userUuid + jewelKeyExtension
	getAll := c.kv.B().Hgetall().Key(key).Build()
	res := c.kv.Do(ctx, getAll)
	if res.Error() != nil {
		return nil, res.Error()
	}
	m, err := res.AsIntMap()
	if err != nil {
		return nil, err
	}

	return models.NewJewelFromMap(m), nil
}