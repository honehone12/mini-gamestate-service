package jewels

import (
	"mini-gamestate-service/db/models"
	"mini-gamestate-service/db/models/colors"
	"mini-gamestate-service/server/context"
	"mini-gamestate-service/server/quick"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ParamUserOnly struct {
	UserUuid string `form:"uuid" validate:"required,uuid4,min=36,max=36"`
}

type ParamsOneColor struct {
	UserUuid  string           `form:"uuid" validate:"required,uuid4,min=36,max=36"`
	ColorCode colors.ColorCode `form:"color" validate:"required,min=0,max=99"`
}

type ParamsIncr struct {
	UserUuid  string           `form:"uuid" validate:"required,uuid4,min=36,max=36"`
	ColorCode colors.ColorCode `form:"color" validate:"required,min=0,max=99"`
	Incr      int64            `form:"incr" validate:"required,min=-99,max99"`
}

func Initialize(c echo.Context) error {
	formData := &ParamUserOnly{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	ctrl := c.(*context.Context).Jewel()
	if err := ctrl.SetAll(formData.UserUuid, &models.Jewel{
		Red:    0,
		Blue:   0,
		Green:  0,
		Yellow: 0,
		Black:  0,
	}); err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	return c.NoContent(http.StatusOK)
}

func GetAll(c echo.Context) error {
	formData := &ParamUserOnly{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	ctrl := c.(*context.Context).Jewel()
	j, err := ctrl.GetAll(formData.UserUuid)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	return c.JSON(http.StatusOK, j)
}

func Incr(c echo.Context) error {
	formData := &ParamsIncr{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	ctrl := c.(*context.Context).Jewel()
	if err := ctrl.IncrBy(
		formData.UserUuid,
		formData.ColorCode,
		formData.Incr,
	); err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	return c.NoContent(http.StatusOK)
}
