package app

import (
	"LibManSys/model"

	"github.com/sirupsen/logrus"
)

type LibraryManagementSystemImpl struct {
	Connector *model.DatabaseConnector
}

var LMS LibraryManagementSystem

func NewLibraryManagementSystemImpl(config *model.ConnectConfig) *LibraryManagementSystemImpl {
	return &LibraryManagementSystemImpl{Connector: &model.DatabaseConnector{
		Config: *config,
	}}
}

func (l *LibraryManagementSystemImpl) Connect() error {
	return l.Connector.Connect()
}

func (l *LibraryManagementSystemImpl) Close() error {
	return l.Connector.Close()
}

func (l *LibraryManagementSystemImpl) Init() error {
	err := l.Connect()
	if err != nil {
		return err
	}
	l.AutoMigrate(
		model.Book{},
		model.Card{},
		model.Borrow{},
	)
	return nil
}

func (l *LibraryManagementSystemImpl) Free() {
	l.Close()
}

func (l *LibraryManagementSystemImpl) AutoMigrate(args ...any) {
	l.Connector.AutoMigrate(args...)
}

func (l *LibraryManagementSystemImpl) StoreBook(book *model.Book) *ApiResult {
	err := l.Connector.StoreBook(book)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) IncBookStock(bookId int, deltaStock int) *ApiResult {
	err := l.Connector.IncBookStock(bookId, deltaStock)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) StoreBooks(books []model.Book) *ApiResult {
	err := l.Connector.StoreBooks(books)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) RemoveBook(bookId int) *ApiResult {
	err := l.Connector.RemoveBook(bookId)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) ModifyBookInfo(book *model.Book) *ApiResult {
	err := l.Connector.ModifyBookInfo(book)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) QueryBook(conditions *model.BookQueryConditions) *ApiResult {
	result, err := l.Connector.QueryBook(conditions)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(result)
}

func (l *LibraryManagementSystemImpl) BorrowBook(borrow *model.Borrow) *ApiResult {
	err := l.Connector.BorrowBook(borrow)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) ReturnBook(borrow *model.Borrow) *ApiResult {
	err := l.Connector.ReturnBook(borrow)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) ShowBorrowHistory(cardId int) *ApiResult {
	borrowHistory, err := l.Connector.ShowBorrowHistory(cardId)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(borrowHistory)
}

func (l *LibraryManagementSystemImpl) RegisterCard(card *model.Card) *ApiResult {
	err := l.Connector.RegisterCard(card)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) QueryCard(cardId int) *ApiResult {
	card, err := l.Connector.QueryCard(cardId)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(card)
}

func (l *LibraryManagementSystemImpl) RemoveCard(cardId int) *ApiResult {
	err := l.Connector.RemoveCard(cardId)
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}

func (l *LibraryManagementSystemImpl) ShowCards() *ApiResult {
	cardList, err := l.Connector.ShowCards()
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(cardList)
}

func (l *LibraryManagementSystemImpl) ResetDatabase() *ApiResult {
	err := l.Connector.ResetDatabase()
	if err != nil {
		logrus.Error(err)
		return &ApiResult{
			OK:      false,
			Message: err.Error(),
		}
	}
	return Success(nil)
}
