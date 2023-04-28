package web

import (
	"LibManSys/app"
	"LibManSys/model"
	"LibManSys/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func createCard(c echo.Context) error {
	var request model.CardCreateRequest

	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind create request",
			Data: nil,
		})
	}

	if request.Name == nil || request.Department == nil || request.Type == nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "missing required field",
			Data: nil,
		})
	}

	card := model.Card{
		Name:       *request.Name,
		Department: *request.Department,
		Type:       *request.Type,
	}

	result := app.LMS.RegisterCard(&card)
	if !result.OK {
		logrus.Error(result.Message)
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Code: utils.E_INTERNAL_SERVER_ERROR,
			Msg:  result.Message,
			Data: nil,
		})
	}
	return c.JSON(http.StatusOK, utils.Success(card.CardID))
}

func queryCard(c echo.Context) error {
	var cid int
	err := echo.QueryParamsBinder(c).MustInt("cid", &cid).BindError()
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind card id param",
			Data: nil,
		})
	}

	result := app.LMS.QueryCard(cid)
	if !result.OK {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Code: utils.E_INTERNAL_SERVER_ERROR,
			Msg:  result.Message,
			Data: nil,
		})
	}
	return c.JSON(http.StatusOK, utils.Success(result.Payload))
}

func listCards(c echo.Context) error {
	result := app.LMS.ShowCards()
	if !result.OK {
		logrus.Error(result.Message)
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Code: utils.E_INTERNAL_SERVER_ERROR,
			Msg:  result.Message,
			Data: nil,
		})
	}

	return c.JSON(http.StatusOK, utils.Success(result.Payload))
}

func removeCard(c echo.Context) error {
	var cid int
	err := echo.QueryParamsBinder(c).MustInt("cid", &cid).BindError()
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind card id param",
			Data: nil,
		})
	}

	result := app.LMS.RemoveCard(cid)
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
