package web

import (
	"LibManSys/app"
	"LibManSys/model"
	"LibManSys/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func createBook(c echo.Context) error {
	var request model.BookCreateRequest

	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind create request",
			Data: nil,
		})
	}

	if request.Category == nil ||
		request.Title == nil ||
		request.Press == nil ||
		request.PublishYear == nil ||
		request.Author == nil ||
		request.Price == nil ||
		request.Stock == nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "missing required field",
			Data: nil,
		})
	}

	book := model.Book{
		Category:    *request.Category,
		Title:       *request.Title,
		Press:       *request.Press,
		PublishYear: *request.PublishYear,
		Author:      *request.Author,
		Price:       *request.Price,
		Stock:       *request.Stock,
	}
	result := app.LMS.StoreBook(&book)
	if !result.OK {
		logrus.Error(result.Message)
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Code: utils.E_INTERNAL_SERVER_ERROR,
			Msg:  result.Message,
			Data: nil,
		})
	}

	return c.JSON(http.StatusOK, utils.Success(book.BookID))
}

func createBookBatch(c echo.Context) error {
	var request []model.BookCreateRequest

	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind create request",
			Data: nil,
		})
	}

	books := make([]model.Book, 0)
	for _, book := range request {
		if book.Category == nil ||
			book.Title == nil ||
			book.Press == nil ||
			book.PublishYear == nil ||
			book.Author == nil ||
			book.Price == nil ||
			book.Stock == nil {
			return c.JSON(http.StatusBadRequest, utils.Error{
				Code: utils.E_BAD_PARAM,
				Msg:  "missing required field",
				Data: nil,
			})
		}
		books = append(books, model.Book{
			Category:    *book.Category,
			Title:       *book.Title,
			Press:       *book.Press,
			PublishYear: *book.PublishYear,
			Author:      *book.Author,
			Price:       *book.Price,
			Stock:       *book.Stock,
		})
	}

	result := app.LMS.StoreBooks(books)
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

func listAllBooksMatched(c echo.Context) error {
	var request model.BookQueryConditions

	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind list request",
			Data: nil,
		})
	}

	result := app.LMS.QueryBook(&request)
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

func updateBook(c echo.Context) error {
	var request model.Book

	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind update request",
			Data: nil,
		})
	}

	result := app.LMS.ModifyBookInfo(&request)
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

func updateBookStock(c echo.Context) error {
	var request model.BookStockUpdateRequest

	err := c.Bind(&request)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind update request",
			Data: nil,
		})
	}

	if request.BookID == nil || request.Option == nil || request.Delta == nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "missing required field",
			Data: nil,
		})
	}

	delta := 0
	if *request.Option == model.INC {
		delta = *request.Delta
	} else if *request.Option == model.DEC {
		delta = -(*request.Delta)
	} else {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "invalid option",
			Data: nil,
		})
	}

	result := app.LMS.IncBookStock(*request.BookID, delta)
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

func removeBook(c echo.Context) error {
	var bid int
	err := echo.QueryParamsBinder(c).MustInt("bid", &bid).BindError()
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, utils.Error{
			Code: utils.E_BAD_PARAM,
			Msg:  "failed to bind book id param",
			Data: nil,
		})
	}

	result := app.LMS.RemoveBook(bid)
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
