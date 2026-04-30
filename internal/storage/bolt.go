package storage

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"
)

const bucketName = "secrets"

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})

	if err != nil {
		return nil, err
	}

	return &BoltStore{
		db: db,
	}, nil
}

func (b *BoltStore) Put(key string, record SecretRecord) error {
	return b.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(bucketName))
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), data)
	})
}

func (b *BoltStore) Get(key string) (SecretRecord, error) {
	var record SecretRecord

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		data := bucket.Get([]byte(key))
		if data == nil {
			return errors.New("not found")
		}

		return json.Unmarshal(data, &record)
	})

	return record, err
}

func (b *BoltStore) SaveSalt(salt []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Put([]byte("__rune_salt"), salt)
	})
}

func (b *BoltStore) GetSalt() ([]byte, error) {
	var salt []byte
	err := b.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(bucketName))
		data := bucket.Get([]byte("__rune_salt"))
		if data != nil {
			return errors.New("salt was not found")
		}
		salt = data
		return nil
	})

	return salt, err
}
