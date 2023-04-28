package web

import (
	"LibManSys/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func addRoutes(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, utils.Success("pong"))
	})

	book := e.Group("/book")
	book.POST("/create", createBook)
	book.POST("/create/batch", createBookBatch)
	book.POST("/list", listAllBooksMatched)
	book.PUT("/update", updateBook)
	book.PUT("/stock/update", updateBookStock)
	book.DELETE("/remove", removeBook)

	card := e.Group("/card")
	card.POST("/create", createCard)
	card.GET("/get", queryCard)
	card.GET("/list", listCards)
	card.DELETE("/remove", removeCard)

	borrow := e.Group("/borrow")
	borrow.GET("/list", queryBorrowHistory)
	borrow.PUT("/borrow", borrowBook)
	borrow.PUT("/return", returnBook)

	db := e.Group("/db")
	db.POST("/reset", resetDatabase)
}
