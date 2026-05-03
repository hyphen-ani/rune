package storage

import (
	"encoding/json"
	"errors"
	"rune/internal/auth"
	"rune/internal/constants"

	bolt "go.etcd.io/bbolt"
)

const bucketName = "secrets"
const tokenBucket = "__rune_tokens"

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
		_, err = tx.CreateBucketIfNotExists([]byte(tokenBucket))
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
			return errors.New("[NOT FOUND] Secret not Found")
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
		if data == nil {
			return errors.New("salt was not found")
		}
		salt = make([]byte, len(data))
		copy(salt, data)
		return nil
	})

	return salt, err
}

func (b *BoltStore) IsValidToken(hash string) bool {
	//TODO implement me
	var exists bool

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tokenBucket))
		val := bucket.Get([]byte(hash))
		exists = val != nil
		return nil
	})
	if err != nil {
		return false
	}

	return exists
}

func (b *BoltStore) SaveToken(record auth.TokenRecord) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists([]byte(tokenBucket))
		data, _ := json.Marshal(record)
		return bucket.Put([]byte(record.ID), data)
	})
}

func (b *BoltStore) ListKeys() ([]string, error) {

	var keys []string
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(k, v []byte) error {
			key := string(k)

			if key == VerifyKey || key == "__rune_salt" {
				return nil
			}

			keys = append(keys, key)
			return nil
		})
	})
	return keys, err

}

func (b *BoltStore) GetTokenByHash(hash string) (*auth.TokenRecord, error) {

	var found *auth.TokenRecord

	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tokenBucket))
		if bucket == nil {
			return nil
		}

		bucket.ForEach(func(k, v []byte) error {
			var rec auth.TokenRecord
			json.Unmarshal(v, &rec)

			if rec.Hash == hash && !rec.Revoked {
				found = &rec
			}

			return nil
		})

		return nil
	})

	if found == nil {
		return nil, errors.New("[INVALID] Invalid Token")
	}

	return found, nil

}

func (b *BoltStore) ListTokens() ([]auth.TokenRecord, error) {
	var tokens []auth.TokenRecord

	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tokenBucket))
		if bucket == nil {
			return nil
		}

		bucket.ForEach(func(k, v []byte) error {
			var rec auth.TokenRecord
			json.Unmarshal(v, &rec)
			tokens = append(tokens, rec)
			return nil
		})

		return nil
	})

	return tokens, nil
}

func (b *BoltStore) RevokeToken(id string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tokenBucket))

		data := bucket.Get([]byte(id))
		if data == nil {
			return errors.New("[NOT FOUND] Token not Found")
		}

		var rec auth.TokenRecord
		json.Unmarshal(data, &rec)

		rec.Revoked = true

		updated, _ := json.Marshal(rec)
		return bucket.Put([]byte(id), updated)
	})
}

func (b *BoltStore) ListNamespaces() ([]string, error) {
	var list []string

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(constants.NamespacesBucket))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			list = append(list, string(k))
			return nil
		})
	})

	return list, err
}

func (b *BoltStore) NamespaceExists(namespace string) bool {
	var exists bool

	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(constants.NamespacesBucket))
		if bucket == nil {
			return nil
		}

		val := bucket.Get([]byte(namespace))
		exists = val != nil
		return nil
	})

	return exists
}

func (b *BoltStore) CreateNamespace(namespace string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists([]byte(constants.NamespacesBucket))
		return bucket.Put([]byte(namespace), []byte("1"))
	})
}
