package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

type Borrow struct {
	CardID     int   `sql:"not null;primaryKey;constraint:Card.CardID,OnDelete:CASCADE,OnUpdate:CASCADE"`
	BookID     int   `sql:"not null;primaryKey;constraint:Book.BookID,OnDelete:CASCADE,OnUpdate:CASCADE"`
	BorrowTime int64 `sql:"not null;primaryKey"`
	ReturnTime int64 `sql:"not null;default:0"`
}

type BorrowRequest struct {
	CardID     *int   `json:"card_id,omitempty"`
	BookID     *int   `json:"book_id,omitempty"`
	BorrowTime *int64 `json:"borrow_time,omitempty"`
	ReturnTime *int64 `json:"return_time,omitempty"`
}

type Item struct {
	CardID      int     `json:"card_id"`
	BookID      int     `json:"book_id"`
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Press       string  `json:"press"`
	PublishYear int     `json:"publish_year"`
	Author      string  `json:"author"`
	Price       myFloat `json:"price"`
	BorrowTime  int64   `json:"borrow_time"`
	ReturnTime  int64   `json:"return_time"`
}

type BorrowHistories struct {
	Count int    `json:"count"`
	Items []Item `json:"items"`
}

func queryBookBorrow(executor SQLExecutor, bookId int) ([]Borrow, error) {
	querySQL := "SELECT * FROM borrow WHERE book_id = ?"

	rows, err := executor.Query(querySQL, bookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrows []Borrow
	for rows.Next() {
		var borrow Borrow
		err = rows.Scan(&borrow.CardID, &borrow.BookID, &borrow.BorrowTime, &borrow.ReturnTime)
		if err != nil {
			return nil, err
		}
		borrows = append(borrows, borrow)
	}
	return borrows, nil
}

func (c *DatabaseConnector) QueryBookBorrow(bookId int) ([]Borrow, error) {
	return queryBookBorrow(c.DB, bookId)
}

func QueryBookBorrowInTx(tx *sql.Tx, bookId int) ([]Borrow, error) {
	return queryBookBorrow(tx, bookId)
}

func (c *DatabaseConnector) BorrowBook(borrow *Borrow) error {
	var (
		queryBookSQL   string
		queryBorrowSQL string
		updateSQL      string
		insertSQL      string
		args           []any
	)

	queryBookSQL = "SELECT stock FROM book WHERE book_id = ? FOR UPDATE"
	args = append(args, borrow.BookID)

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query(queryBookSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	var stock int
	if rows.Next() {
		rows.Scan(&stock)
	} else {
		tx.Rollback()
		return errors.New("book not found")
	}
	rows.Close()

	if stock <= 0 {
		tx.Rollback()
		return errors.New("stock not enough")
	}

	args = args[:0]
	queryBorrowSQL = "SELECT * FROM borrow WHERE book_id = ? AND card_id = ? AND return_time = 0"
	args = append(args, borrow.BookID, borrow.CardID)

	rows, err = tx.Query(queryBorrowSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows.Next() {
		tx.Rollback()
		return errors.New("book already borrowed")
	}
	rows.Close()

	args = args[:0]
	updateSQL = "UPDATE book SET stock = ? WHERE book_id = ?"
	args = append(args, stock-1, borrow.BookID)

	_, err = tx.Exec(updateSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	args = args[:0]
	insertSQL = "INSERT INTO borrow (card_id, book_id, borrow_time, return_time) VALUES (?, ?, ?, ?)"
	args = append(args, borrow.CardID, borrow.BookID, time.Now().Unix(), 0)

	_, err = tx.Exec(insertSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (c *DatabaseConnector) ReturnBook(borrow *Borrow) error {
	var (
		querySQL        string
		updateBorrowSQL string
		updateBookSQL   string
		args            []any
	)

	querySQL =
		"SELECT * FROM borrow WHERE book_id = ? AND card_id = ? AND return_time = 0 LOCK IN SHARE MODE"
	args = append(args, borrow.BookID, borrow.CardID)
	logrus.Info(querySQL, args)

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query(querySQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !rows.Next() {
		tx.Rollback()
		return errors.New("book not borrowed")
	}
	rows.Close()

	args = args[:0]
	updateBorrowSQL =
		"UPDATE borrow SET return_time = ? WHERE book_id = ? AND card_id = ? AND return_time = 0"
	args = append(args, time.Now().Unix(), borrow.BookID, borrow.CardID)

	_, err = tx.Exec(updateBorrowSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	args = args[:0]
	updateBookSQL = "UPDATE book SET stock = stock + 1 WHERE book_id = ?"
	args = append(args, borrow.BookID)

	_, err = tx.Exec(updateBookSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (c *DatabaseConnector) ShowBorrowHistory(cardId int) (*BorrowHistories, error) {
	var (
		queryCardSQL   string
		queryBorrowSQL string
		queryBookSQL   string
	)

	queryCardSQL = "SELECT * FROM card WHERE card_id = ?"

	tx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(queryCardSQL, cardId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if !rows.Next() {
		tx.Rollback()
		return nil, errors.New("card not found")
	}
	rows.Close()

	queryBorrowSQL = "SELECT * FROM borrow WHERE card_id = ?"

	rows, err = tx.Query(queryBorrowSQL, cardId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var borrows []Borrow

	for rows.Next() {
		var borrow Borrow
		err = rows.Scan(&borrow.CardID, &borrow.BookID, &borrow.BorrowTime, &borrow.ReturnTime)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		borrows = append(borrows, borrow)
	}
	rows.Close()

	histories := BorrowHistories{
		Count: 0,
		Items: make([]Item, 0),
	}
	queryBookSQL = "SELECT * FROM book WHERE book_id = ? LOCK IN SHARE MODE"
	stmt, err := tx.Prepare(queryBookSQL)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	for _, borrow := range borrows {
		bookRows, err := stmt.Query(borrow.BookID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		var book Book
		if bookRows.Next() {
			err = bookRows.Scan(
				&book.BookID,
				&book.Category,
				&book.Title,
				&book.Press,
				&book.PublishYear,
				&book.Author,
				&book.Price,
				&book.Stock,
			)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			tx.Rollback()
			return nil, errors.New("book not found")
		}
		bookRows.Close()

		histories.Count++
		histories.Items = append(histories.Items, Item{
			CardID:      cardId,
			BookID:      book.BookID,
			Category:    book.Category,
			Title:       book.Title,
			Press:       book.Press,
			PublishYear: book.PublishYear,
			Author:      book.Author,
			Price:       book.Price,
			BorrowTime:  borrow.BorrowTime,
			ReturnTime:  borrow.ReturnTime,
		})
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &histories, nil
}
