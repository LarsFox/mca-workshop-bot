package storage

import (
	"strconv"

	"github.com/boltdb/bolt"
)

// Client saves data to Bolt.
type Client struct {
	db *bolt.DB
}

var (
	bucketUsersMessages = []byte("users_messages")
)

// NewClient returns a new client creating a DB file if it does not exist.
func NewClient(dbPath string) (*Client, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketUsersMessages); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

// GetUserMessage returns a user message or an empty string if the message is not found.
func (c *Client) GetUserMessage(userID int64) (string, error) {
	var msg string
	if err := c.db.View(func(tx *bolt.Tx) error {
		val := tx.Bucket(bucketUsersMessages).Get(uID(userID))
		if val == nil {
			return nil
		}
		msg = string(val)
		return nil
	}); err != nil {
		return "", err
	}
	return msg, nil
}

// SaveUserMessage saves a user message.
func (c *Client) SaveUserMessage(userID int64, msg string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketUsersMessages).Put(uID(userID), []byte(msg))
	})
}

// DeleteUserMessage deletes a user message.
func (c *Client) DeleteUserMessage(userID int64) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketUsersMessages).Delete(uID(userID))
	})
}

func uID(userID int64) []byte {
	return []byte(strconv.FormatInt(userID, 10))
}
