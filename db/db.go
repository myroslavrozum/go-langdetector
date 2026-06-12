package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

// Ensure all functions are exported for use outside this package

func InitDB(database string) (*badger.DB, error) {
	//opt := badger.DefaultOptions("").WithInMemory(true)
	opt := badger.DefaultOptions(database)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetValueFromDB(db *badger.DB, key string) (string, error) {
	var rv string
	err := db.View(func(txn *badger.Txn) error {
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

func PutValueToDB(db *badger.DB, key string, value string) error {
	valLen := min(10, len(value))
	log.Printf("Putting %s -> %s to DB.....", key, value[:valLen])
	txn := db.NewTransaction(true)
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

func RestoreTrigrammes(db *badger.DB, language string) (map[string]float64, error) {
	log.Println("Restoring data.....")
	if language == "" {
		return nil, errors.New("language cannot be empty")
	}

	trigrammes_str, err := GetValueFromDB(db, language)
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

func DumpTrigrammes(db *badger.DB, data map[string]map[string]float64) error {
	txn := db.NewTransaction(true)
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
			txn = db.NewTransaction(true)
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
