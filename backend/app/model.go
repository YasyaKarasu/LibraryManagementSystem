package app

import "LibManSys/model"

type LibraryManagementSystem interface {
	Init() error
	Free()

	StoreBook(*model.Book) *ApiResult
	IncBookStock(bookId int, deltaStock int) *ApiResult
	StoreBooks([]model.Book) *ApiResult
	RemoveBook(bookId int) *ApiResult
	ModifyBookInfo(*model.Book) *ApiResult
	QueryBook(*model.BookQueryConditions) *ApiResult
	BorrowBook(*model.Borrow) *ApiResult
	ReturnBook(*model.Borrow) *ApiResult
	ShowBorrowHistory(cardId int) *ApiResult
	RegisterCard(*model.Card) *ApiResult
	QueryCard(cardId int) *ApiResult
	RemoveCard(cardId int) *ApiResult
	ShowCards() *ApiResult
	ResetDatabase() *ApiResult
}

type ApiResult struct {
	OK      bool
	Message string
	Payload any
}

func Success(payload any) *ApiResult {
	return &ApiResult{
		OK:      true,
		Message: "success",
		Payload: payload,
	}
}
