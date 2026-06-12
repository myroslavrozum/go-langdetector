package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

// Store wraps a badger.DB instance to provide domain-specific database methods.
type Store struct {
	db *badger.DB
}

// NewStore initializes and returns a new Store.
func NewStore(database string) (*Store, error) {
	//opt := badger.DefaultOptions("").WithInMemory(true)
	opt := badger.DefaultOptions(database)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// Close cleanly shuts down the database.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) GetValue(key string) (string, error) {
	var rv string
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))

		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			rv = string(val)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return rv, nil
}

func (s *Store) PutValue(key string, value string) error {
	valLen := min(10, len(value))
	log.Printf("Putting %s -> %s to DB.....", key, value[:valLen])
	txn := s.db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) RestoreTrigrammes(language string) (map[string]float64, error) {
	log.Println("Restoring data.....")
	if language == "" {
		return nil, errors.New("language cannot be empty")
	}

	trigrammes_str, err := s.GetValue(language)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, fmt.Errorf("no trigrammes found for language %s", language)
		}
		return nil, err
	}

	var trigrammes map[string]float64
	if err := json.Unmarshal([]byte(trigrammes_str), &trigrammes); err != nil {
		return nil, err
	}
	return trigrammes, nil
}

func (s *Store) DumpTrigrammes(data map[string]map[string]float64) error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard() // Automatically rolled back if not committed

	for language, trigrammes := range data {
		trigrammes_str, err := json.Marshal(trigrammes)
		if err != nil {
			return err
		}

		strLen := min(20, len(trigrammes_str))
		log.Printf("Dumping data for %s language: %-20s......", language, string(trigrammes_str[:strLen]))

		if err := txn.Set([]byte(language), trigrammes_str); errors.Is(err, badger.ErrTxnTooBig) {
			if err := txn.Commit(); err != nil {
				return err
			}
			txn = s.db.NewTransaction(true)
			if err := txn.Set([]byte(language), trigrammes_str); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	err := txn.Commit()
	return err
}
