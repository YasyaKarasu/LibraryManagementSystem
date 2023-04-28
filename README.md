# LibraryManagementSystem

A toy library management system, the lab 5 of ZJU Database Concept System course.

## Backend

Use Golang as the backend language.

| Database Driver     | Web Framework | Log Manager     | Config Manager |
| ------------------- | ------------- | --------------- | -------------- |
| go-sql-driver/mysql | labstack/echo | sirupsen/logrus | spf13/viper    |

### Code Structure

#### app

Model for the library management system and one implement with MySQL as the data storage. Any type that implement the following methods can be a possible library management system.

```go
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
```

#### conf

Manage configurations. Load `conf.yaml` file and provide MySQL login info to database connector.

#### model

Models for book, card and borrow record and functions to handle database.

#### utils

Utilities. Currently contains http response errors and string_to_snake function.

#### web

Use echo to implement http api for library operation.

## Frontend

Use Vue as the frontend language. Just learned for this lab so the code quality is poor :(

Use Naive-UI as the component library.