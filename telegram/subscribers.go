package telegram

import (
	"errors"
	"strconv"
	"time"

	"github.com/roemerb/schiphol-runway-alerts/db"
)

// Subscriber is a Telegram user that has subscribed to the bot
type Subscriber struct {
	ID             int
	FirstName      string
	Username       string
	TelegramUserID int
	TelegramChatID int
	IsBot          bool
	LanguageCode   string
	RegisteredAt   string
}

// SubscriberRepository can be passed around and holds functions for
// storing/retrieving subscribers from the database
type SubscriberRepository struct{}

// GetOrCreateSubscriber checks if there is a subscriber in the database with user.ID
// and if not, creates it
func (s SubscriberRepository) GetOrCreateSubscriber(user *User) (*Subscriber, error) {
	sub, err := s.GetByTelegramID(user.ID)
	if sub == nil {
		sub, err = s.GetOrCreateSubscriber(user)
	}

	return sub, err
}

// GetByTelegramID searches the database for a user with Telegram ID id
func (s SubscriberRepository) GetByTelegramID(id int) (*Subscriber, error) {
	db, err := db.OpenDatabase()
	defer db.Close()
	if err != nil {
		return nil, errors.New("Could not connet to database: " + err.Error())
	}

	query := `
	SELECT * FROM subscribers
	WHERE telegram_user_id = ?;
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, errors.New("Error occured while executing query: " + err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No user with Telegram user ID " + strconv.Itoa(id))
	}

	var sub Subscriber
	rows.Scan(
		&sub.ID,
		&sub.FirstName,
		&sub.Username,
		&sub.TelegramUserID,
		&sub.TelegramChatID,
		&sub.IsBot,
		&sub.LanguageCode,
		&sub.RegisteredAt,
	)

	return &sub, nil
}

// SubscriberFromTelegramUser will store a subscriber in the database based on a Telegram user
// and initiates a Subscriber struct from the inserted user to return
func (s SubscriberRepository) SubscriberFromTelegramUser(user *User) (*Subscriber, error) {
	db, err := db.OpenDatabase()
	defer db.Close()
	if err != nil {
		return &Subscriber{}, errors.New("Could not connect to databse: " + err.Error())
	}

	query := `INSERT INTO subscribers (
		first_name,
		username,
		telegram_user_id,
		telegram_chat_id,
		is_bot,
		language_code,
		registered_at
	)
		VALUES (?,?,?,?,?,?,?);`
	stmt, err := db.Prepare(query)
	if err != nil {
		return &Subscriber{}, errors.New("Could not initiate query: " + err.Error())
	}
	defer stmt.Close()

	now := time.Now()
	registedAt := now.Format("2016-01-02 15:04:05")

	res, err := stmt.Exec(
		user.FirstName,
		user.Username,
		user.ID,
		user.ChatID,
		user.IsBot,
		user.LanguageCode,
		registedAt,
	)
	if err != nil {
		return &Subscriber{}, errors.New("Insertion error: " + err.Error())
	}

	ID, _ := res.LastInsertId()
	sub := Subscriber{
		ID:             int(ID),
		FirstName:      user.FirstName,
		Username:       user.Username,
		TelegramUserID: user.ID,
		TelegramChatID: user.ChatID,
		IsBot:          user.IsBot,
		LanguageCode:   user.LanguageCode,
		RegisteredAt:   registedAt,
	}

	return &sub, nil
}

// GetAllSubscribers retrieves all current subscribers from the database
func (s SubscriberRepository) GetAllSubscribers() ([]*Subscriber, error) {
	db, err := db.OpenDatabase()
	defer db.Close()
	if err != nil {
		return nil, errors.New("Could not connet to database: " + err.Error())
	}

	query := "SELECT * FROM subscribers;"
	rows, err := db.Query(query)

	var subs []*Subscriber
	for rows.Next() {
		var sub Subscriber
		rows.Scan(
			&sub.ID,
			&sub.FirstName,
			&sub.Username,
			&sub.TelegramUserID,
			&sub.TelegramChatID,
			&sub.IsBot,
			&sub.LanguageCode,
			&sub.RegisteredAt,
		)
		subs = append(subs, &sub)
	}

	return subs, nil
}
