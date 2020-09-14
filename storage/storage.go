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
	bucketUsersModels = []byte("users_models")
)

// NewClient returns a new client creating a DB file if it does not exist.
func NewClient(dbPath string) (*Client, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketUsersModels); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

// GetUserModel returns a user chosen model or an empty string if none was chosen.
func (c *Client) GetUserModel(userID int64) (string, error) {
	var msg string
	if err := c.db.View(func(tx *bolt.Tx) error {
		val := tx.Bucket(bucketUsersModels).Get(uID(userID))
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

// SaveUserModel saves a user chosen model.
func (c *Client) SaveUserModel(userID int64, model string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketUsersModels).Put(uID(userID), []byte(model))
	})
}

func uID(userID int64) []byte {
	return []byte(strconv.FormatInt(userID, 10))
}
