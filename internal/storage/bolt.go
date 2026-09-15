package storage

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

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

func (b *BoltStore) Delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("[NOT FOUND] Secret not Found")
		}
		return bucket.Delete([]byte(key))
	})
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

func (b *BoltStore) GetStats() (VaultStats, error) {
	var stats VaultStats
	secretsByNS := map[string]int{}
	rotatedKeys := map[string]bool{}
	tokensByMonth := map[string]int{}

	err := b.db.View(func(tx *bolt.Tx) error {
		// --- secrets pass ---
		sb := tx.Bucket([]byte(bucketName))
		if sb != nil {
			sb.ForEach(func(k, _ []byte) error {
				key := string(k)
				if key == VerifyKey || key == "__rune_salt" {
					return nil
				}
				if idx := strings.Index(key, "@v"); idx != -1 {
					vStr := key[idx+2:]
					if v, err := strconv.Atoi(vStr); err == nil && v >= 2 {
						rotatedKeys[key[:idx]] = true
					}
				} else {
					ns := key
					if i := strings.Index(key, "/"); i != -1 {
						ns = key[:i]
					}
					secretsByNS[ns]++
					stats.TotalSecrets++
				}
				return nil
			})
		}

		// --- tokens pass ---
		tb := tx.Bucket([]byte(tokenBucket))
		if tb != nil {
			tb.ForEach(func(_, v []byte) error {
				var rec auth.TokenRecord
				if err := json.Unmarshal(v, &rec); err != nil {
					return nil
				}
				stats.TotalTokens++
				if rec.Revoked {
					stats.RevokedTokens++
				} else {
					stats.ActiveTokens++
				}
				if t, err := time.Parse(time.RFC3339, rec.CreatedAt); err == nil {
					tokensByMonth[t.Format("2006-01")]++
				}
				return nil
			})
		}

		return nil
	})

	if err != nil {
		return stats, err
	}

	stats.RotatedSecrets = len(rotatedKeys)
	stats.TotalNamespaces = len(secretsByNS)

	for ns, count := range secretsByNS {
		stats.SecretsByNamespace = append(stats.SecretsByNamespace, NamespaceCount{Namespace: ns, Count: count})
	}
	sort.Slice(stats.SecretsByNamespace, func(i, j int) bool {
		return stats.SecretsByNamespace[i].Count > stats.SecretsByNamespace[j].Count
	})

	for month, count := range tokensByMonth {
		stats.TokensByMonth = append(stats.TokensByMonth, MonthCount{Month: month, Count: count})
	}
	sort.Slice(stats.TokensByMonth, func(i, j int) bool {
		return stats.TokensByMonth[i].Month < stats.TokensByMonth[j].Month
	})

	return stats, nil
}

func (b *BoltStore) DeleteNamespace(namespace string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(constants.NamespacesBucket))
		if bucket == nil {
			return nil
		}
		return bucket.Delete([]byte(namespace))
	})
}
