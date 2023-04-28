package model

import (
	"errors"
	"fmt"
)

type myFloat float32

func (f myFloat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%0.2f", f)), nil
}

type Book struct {
	BookID      int     `json:"book_id" sql:"not null;autoIncrement;primaryKey"`
	Category    string  `json:"category" sql:"not null;size:63;unique:book_unique"`
	Title       string  `json:"title" sql:"not null;size:63;unique:book_unique"`
	Press       string  `json:"press" sql:"not null;size:63;unique:book_unique"`
	PublishYear int     `json:"publish_year" sql:"not null;unique:book_unique"`
	Author      string  `json:"author" sql:"not null;size:63;unique:book_unique"`
	Price       myFloat `json:"price" sql:"not null;decimal:7,2;default:0.00"`
	Stock       int     `json:"stock" sql:"not null;default:0"`
}

type BookCreateRequest struct {
	Category    *string  `json:"category,omitempty"`
	Title       *string  `json:"title,omitempty"`
	Press       *string  `json:"press,omitempty"`
	PublishYear *int     `json:"publish_year,omitempty"`
	Author      *string  `json:"author,omitempty"`
	Price       *myFloat `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

type IncStockOption string

const (
	INC IncStockOption = "inc"
	DEC IncStockOption = "dec"
)

type BookStockUpdateRequest struct {
	BookID *int            `json:"book_id,omitempty"`
	Option *IncStockOption `json:"option,omitempty"`
	Delta  *int            `json:"delta,omitempty"`
}

type BookColumn string

const (
	BookColumnBookID      BookColumn = "book_id"
	BookColumnCategory    BookColumn = "category"
	BookColumnTitle       BookColumn = "title"
	BookColumnPress       BookColumn = "press"
	BookColumnPublishYear BookColumn = "publish_year"
	BookColumnAuthor      BookColumn = "author"
	BookColumnPrice       BookColumn = "price"
	BookColumnStock       BookColumn = "stock"
)

type SortOrder string

const (
	Ascending  SortOrder = "ASC"
	Descending SortOrder = "DESC"
)

type BookQueryConditions struct {
	Category       *string     `json:"category,omitempty"`
	Title          *string     `json:"title,omitempty"`
	Press          *string     `json:"press,omitempty"`
	MinPublishYear *int        `json:"min_publish_year,omitempty"`
	MaxPublishYear *int        `json:"max_publish_year,omitempty"`
	Author         *string     `json:"author,omitempty"`
	MinPrice       *float32    `json:"min_price,omitempty"`
	MaxPrice       *float32    `json:"max_price,omitempty"`
	SortBy         *BookColumn `json:"sort_by,omitempty"`
	SortOrder      *SortOrder  `json:"sort_order,omitempty"`
}

func NewBookQueryConditions() *BookQueryConditions {
	return &BookQueryConditions{}
}

func (c *BookQueryConditions) WithCategory(category string) *BookQueryConditions {
	c.Category = &category
	return c
}

func (c *BookQueryConditions) WithTitle(title string) *BookQueryConditions {
	c.Title = &title
	return c
}

func (c *BookQueryConditions) WithPress(press string) *BookQueryConditions {
	c.Press = &press
	return c
}

func (c *BookQueryConditions) WithMinPublishYear(minPublishYear int) *BookQueryConditions {
	c.MinPublishYear = &minPublishYear
	return c
}

func (c *BookQueryConditions) WithMaxPublishYear(maxPublishYear int) *BookQueryConditions {
	c.MaxPublishYear = &maxPublishYear
	return c
}

func (c *BookQueryConditions) WithAuthor(author string) *BookQueryConditions {
	c.Author = &author
	return c
}

func (c *BookQueryConditions) WithMinPrice(minPrice float32) *BookQueryConditions {
	c.MinPrice = &minPrice
	return c
}

func (c *BookQueryConditions) WithMaxPrice(maxPrice float32) *BookQueryConditions {
	c.MaxPrice = &maxPrice
	return c
}

func (c *BookQueryConditions) WithSortBy(sortBy BookColumn) *BookQueryConditions {
	c.SortBy = &sortBy
	return c
}

func (c *BookQueryConditions) WithSortOrder(sortOrder SortOrder) *BookQueryConditions {
	c.SortOrder = &sortOrder
	return c
}

func (c *BookQueryConditions) Build() *BookQueryConditions {
	return c
}

type BookQueryResult struct {
	Count   int    `json:"count"`
	Results []Book `json:"results"`
}

func (c *DatabaseConnector) StoreBook(book *Book) error {
	if book == nil {
		return errors.New("book is nil")
	}

	var (
		insertSQL string
		args      []any
	)
	insertSQL = "INSERT INTO book (category, title, press, publish_year, author, price, stock) VALUES (?, ?, ?, ?, ?, ?, ?)"
	args = append(args, book.Category, book.Title, book.Press, book.PublishYear, book.Author, book.Price, book.Stock)

	result, err := c.DB.Exec(insertSQL, args...)
	if err != nil {
		return err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	book.BookID = int(insertedID)
	return nil
}

func (c *DatabaseConnector) IncBookStock(bookId int, deltaStock int) error {
	var (
		querySQL  string
		updateSQL string
		args      []any
	)

	querySQL = "SELECT stock FROM book WHERE book_id = ? FOR UPDATE"
	args = append(args, bookId)

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query(querySQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	var stock int
	if rows.Next() {
		err = rows.Scan(&stock)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return errors.New("book not found")
	}
	rows.Close()

	if stock+deltaStock < 0 {
		tx.Rollback()
		return errors.New("stock not enough")
	}

	args = args[:0]
	updateSQL = "UPDATE book SET stock = ? WHERE book_id = ?"
	args = append(args, stock+deltaStock, bookId)

	_, err = tx.Exec(updateSQL, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (c *DatabaseConnector) StoreBooks(books []Book) error {
	if books == nil {
		return errors.New("books is nil")
	}
	if len(books) == 0 {
		return errors.New("books is empty")
	}

	var (
		insertSQL string
		args      []any
	)
	insertSQL = "INSERT INTO book (category, title, press, publish_year, author, price, stock) VALUES (?, ?, ?, ?, ?, ?, ?)"

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, book := range books {
		args = args[:0]
		args = append(args, book.Category, book.Title, book.Press, book.PublishYear, book.Author, book.Price, book.Stock)
		_, err = stmt.Exec(args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (c *DatabaseConnector) RemoveBook(bookId int) error {
	deleteSQL := "DELETE FROM book WHERE book_id = ?"

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	borrows, err := QueryBookBorrowInTx(tx, bookId)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, borrow := range borrows {
		if borrow.ReturnTime == 0 {
			tx.Rollback()
			return errors.New("book is borrowed")
		}
	}

	_, err = tx.Exec(deleteSQL, bookId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (c *DatabaseConnector) ModifyBookInfo(book *Book) error {
	if book == nil {
		return errors.New("book is nil")
	}

	var (
		updateSQL string
		args      []any
	)
	updateSQL = "UPDATE book SET category = ?, title = ?, press = ?, publish_year = ?, author = ?, price = ? WHERE book_id = ?"
	args = append(args, book.Category, book.Title, book.Press, book.PublishYear, book.Author, book.Price, book.BookID)

	_, err := c.DB.Exec(updateSQL, args...)
	return err
}

func (c *DatabaseConnector) QueryBook(condition *BookQueryConditions) (*BookQueryResult, error) {
	var (
		querySQL string
		args     []any
	)
	querySQL = "SELECT * FROM book"
	if condition != nil {
		if condition.Category != nil {
			querySQL += " WHERE category = ?"
			args = append(args, *condition.Category)
		}
		if condition.Title != nil {
			if len(args) == 0 {
				querySQL += " WHERE title LIKE ?"
			} else {
				querySQL += " AND title LIKE ?"
			}
			args = append(args, "%"+*condition.Title+"%")
		}
		if condition.Press != nil {
			if len(args) == 0 {
				querySQL += " WHERE press LIKE ?"
			} else {
				querySQL += " AND press LIKE ?"
			}
			args = append(args, "%"+*condition.Press+"%")
		}
		if condition.MinPublishYear != nil {
			if len(args) == 0 {
				querySQL += " WHERE publish_year >= ?"
			} else {
				querySQL += " AND publish_year >= ?"
			}
			args = append(args, *condition.MinPublishYear)
		}
		if condition.MaxPublishYear != nil {
			if len(args) == 0 {
				querySQL += " WHERE publish_year <= ?"
			} else {
				querySQL += " AND publish_year <= ?"
			}
			args = append(args, *condition.MaxPublishYear)
		}
		if condition.Author != nil {
			if len(args) == 0 {
				querySQL += " WHERE author LIKE ?"
			} else {
				querySQL += " AND author LIKE ?"
			}
			args = append(args, "%"+*condition.Author+"%")
		}
		if condition.MinPrice != nil {
			if len(args) == 0 {
				querySQL += " WHERE price >= ?"
			} else {
				querySQL += " AND price >= ?"
			}
			args = append(args, *condition.MinPrice)
		}
		if condition.MaxPrice != nil {
			if len(args) == 0 {
				querySQL += " WHERE price <= ?"
			} else {
				querySQL += " AND price <= ?"
			}
			args = append(args, *condition.MaxPrice)
		}
		if condition.SortBy != nil {
			querySQL += " ORDER BY " + string(*condition.SortBy)
			if condition.SortOrder != nil {
				if *condition.SortOrder != Ascending && *condition.SortOrder != Descending {
					return nil, errors.New("invalid sort order")
				}
				querySQL += " " + string(*condition.SortOrder)
			}
			querySQL += ", book_id ASC"
		} else {
			querySQL += " ORDER BY book_id ASC"
		}
	}

	rows, err := c.DB.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result BookQueryResult
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.BookID, &book.Category, &book.Title, &book.Press, &book.PublishYear, &book.Author, &book.Price, &book.Stock)
		if err != nil {
			return nil, err
		}
		result.Results = append(result.Results, book)
	}
	result.Count = len(result.Results)
	return &result, nil
}
