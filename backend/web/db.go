package web

import (
	"LibManSys/app"
	"LibManSys/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func resetDatabase(c echo.Context) error {
	result := app.LMS.ResetDatabase()
	if !result.OK {
		logrus.Error(result.Message)
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Code: utils.E_INTERNAL_SERVER_ERROR,
			Msg:  result.Message,
			Data: nil,
		})
	}
	return c.JSON(http.StatusOK, utils.Success(nil))
}
