package sessions

import (
	"errors"
	"mini-gamestate-service/server/context"
	"mini-gamestate-service/server/quick"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type SessionForm struct {
	UserUuid         string `form:"uuid" validate:"required,alphanum,min=36,max=36"`
	OneTimeSessionId string `form:"id" validate:"required,alphanum,min=44,max=44"`
}

var (
	ErrorOnetimeIdExpired = errors.New("onetime id is expired")
	ErrorInvalidOneTimeId = errors.New("onetime id is invalid")
)

func Set(c echo.Context) error {
	formData := &SessionForm{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	ctrl := c.(*context.Context).Session()
	if err := ctrl.Set(formData.UserUuid, formData.OneTimeSessionId); err != nil {
		c.Logger().Warn(err)
		return quick.ServiceError()
	}

	return c.NoContent(http.StatusOK)
}

func Verify(c echo.Context) error {
	formData := &SessionForm{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	ctrl := c.(*context.Context).Session()
	s, err := ctrl.Get(formData.UserUuid)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	expiration := s.CreatedAtUnix + 60*60
	now := time.Now().Unix()
	if now > expiration {
		c.Logger().Warn(ErrorOnetimeIdExpired)
		return quick.BadRequest()
	}

	if strings.Compare(formData.OneTimeSessionId, s.OneTimeId) != 0 {
		c.Logger().Warn(ErrorInvalidOneTimeId)
		return quick.BadRequest()
	}

	return c.NoContent(http.StatusOK)
}
