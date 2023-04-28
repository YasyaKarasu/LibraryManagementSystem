package web

import (
	"LibManSys/app"
	"LibManSys/model"
	"LibManSys/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func borrowBook(c echo.Context) error {
	var request model.BorrowRequest
	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind borrow request",
			Data: nil,
		})
	}

	if request.CardID == nil || request.BookID == nil || request.BorrowTime == nil {
		logrus.Error("missing required field")
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "missing required field",
			Data: nil,
		})
	}

	result := app.LMS.BorrowBook(&model.Borrow{
		CardID:     *request.CardID,
		BookID:     *request.BookID,
		BorrowTime: *request.BorrowTime,
		ReturnTime: *request.ReturnTime,
	})
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

func returnBook(c echo.Context) error {
	var request model.BorrowRequest
	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind return request",
			Data: nil,
		})
	}

	if request.CardID == nil || request.BookID == nil || request.ReturnTime == nil {
		logrus.Error("missing required field")
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "missing required field",
			Data: nil,
		})
	}

	result := app.LMS.ReturnBook(&model.Borrow{
		CardID:     *request.CardID,
		BookID:     *request.BookID,
		BorrowTime: *request.BorrowTime,
		ReturnTime: *request.ReturnTime,
	})
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

func queryBorrowHistory(c echo.Context) error {
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

	result := app.LMS.ShowBorrowHistory(cid)
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
