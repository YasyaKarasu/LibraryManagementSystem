package model

import "errors"

type Card struct {
	CardID     int    `json:"card_id" sql:"not null;autoIncrement;primaryKey"`
	Name       string `json:"name" sql:"not null;size:63;unique:card_unique"`
	Department string `json:"department" sql:"not null;size:63;unique:card_unique"`
	Type       string `json:"type" sql:"not null;char:1;unique:card_unique;check:type in ('T', 'S')"`
}

type CardList struct {
	Count int    `json:"count"`
	Cards []Card `json:"cards"`
}

type CardCreateRequest struct {
	Name       *string `json:"name,omitempty"`
	Department *string `json:"department,omitempty"`
	Type       *string `json:"type,omitempty"`
}

func (c *DatabaseConnector) RegisterCard(card *Card) error {
	if card == nil {
		return errors.New("card is nil")
	}

	var (
		insertSQL string
		args      []any
	)

	insertSQL = "INSERT INTO card (name, department, type) VALUES (?, ?, ?)"
	args = append(args, card.Name, card.Department, card.Type)

	result, err := c.DB.Exec(insertSQL, args...)
	if err != nil {
		return err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	card.CardID = int(insertedID)
	return nil
}

func (c *DatabaseConnector) RemoveCard(cardId int) error {
	var (
		querySQL  string
		deleteSQL string
	)

	querySQL = "SELECT * FROM borrow WHERE card_id = ? AND return_time = 0"

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query(querySQL, cardId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows.Next() {
		tx.Rollback()
		return errors.New("card has not returned all books")
	}
	rows.Close()

	deleteSQL = "DELETE FROM card WHERE card_id = ?"

	_, err = tx.Exec(deleteSQL, cardId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (c *DatabaseConnector) ShowCards() (*CardList, error) {
	querySQL := "SELECT * FROM card ORDER BY card_id"

	rows, err := c.DB.Query(querySQL)
	if err != nil {
		return nil, err
	}

	var cards []Card
	for rows.Next() {
		var card Card
		err := rows.Scan(&card.CardID, &card.Name, &card.Department, &card.Type)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return &CardList{
		Count: len(cards),
		Cards: cards,
	}, nil
}

func (c *DatabaseConnector) QueryCard(cardId int) (*Card, error) {
	querySQL := "SELECT * FROM card WHERE card_id = ?"

	row := c.DB.QueryRow(querySQL, cardId)

	var card Card
	err := row.Scan(&card.CardID, &card.Name, &card.Department, &card.Type)
	if err != nil {
		return nil, err
	}

	return &card, nil
}
