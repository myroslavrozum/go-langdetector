package db

import (
	"encoding/json"
	"errors"
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
	log.Printf("Putting %s -> %s to DB.....", key, value[:10])
	txn := db.NewTransaction(true)
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

	if trigrammes_str, err := GetValueFromDB(db, language); err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, errors.New("no trigrammes found for language " + language)
		}
		return nil, err
	} else {
		trigrammes := make(map[string]float64, 0)
		if err := json.Unmarshal([]byte(trigrammes_str), &trigrammes); err != nil {
			return nil, err
		}
		return trigrammes, nil
	}
}

func DumpTrigrammes(db *badger.DB, data map[string]map[string]float64) error {
	txn := db.NewTransaction(true)
	for language, trigrammes := range data {
		if trigrammes_str, err := json.Marshal(trigrammes); err != nil {
			return err
		} else {
			log.Printf("Dumping data: %s......", trigrammes_str[:10])
			if err := txn.Set([]byte(language), trigrammes_str); err == badger.ErrTxnTooBig {
				txn.Commit()
				txn = db.NewTransaction(true)
				txn.Set([]byte(language), []byte(trigrammes_str))
			}
		}
	}
	err := txn.Commit()
	return err
}
