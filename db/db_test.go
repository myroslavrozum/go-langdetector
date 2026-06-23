package db

import (
	"log"
	"os"
	"testing"
)

var store *Store
var dbErr error
var dirPath = "data/langdetector-badger-db"
var testData = map[string]map[string]float64{
	"en": {
		"wo":  0.09090909090909091,
		"w":   0.09090909090909091,
		"Hel": 0.09090909090909091,
		"ell": 0.09090909090909091,
		"ld!": 0.09090909090909091,
		"llo": 0.09090909090909091,
		"lo,": 0.09090909090909091,
		"o, ": 0.09090909090909091,
		"orl": 0.09090909090909091,
		"rld": 0.09090909090909091,
		"wor": 0.09090909090909091,
	}}

func setup() {
	log.Println("Performing setup before all tests...")
	if store, dbErr = NewStore(dirPath); dbErr != nil {
		log.Fatal(dbErr)
	}
}

func teardown() {
	log.Println("Performing teardown after all tests...")
	store.Close()
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Printf("Failed to delete directory: %v\n", err)
		return
	}
}

func TestMain(m *testing.M) {
	// 1. Run your setup function
	setup()
	// 2. Execute all the tests in the package
	exitCode := m.Run()

	// 3. Run your teardown function (optional)
	teardown()

	// 4. Exit with the correct code returned by the test suite
	os.Exit(exitCode)
}

func TestNewStore(t *testing.T) {
	if dbErr != nil {
		t.Errorf(`Error initializing store: %v`, dbErr)
	}
}

func TestPutValue(t *testing.T) {
	if err := store.PutValue(`test`, `test value`); err != nil {
		t.Errorf(`Error puting test value to the database: %v`, err)
	}
}

func TestClose(t *testing.T) {}

func TestGetValue(t *testing.T) {
	store.PutValue(`test`, `test value`)
	value, err := store.GetValue(`test`)
	if err != nil {
		t.Errorf(`Error getting test value from the database: %v`, err)
	}
	if value != `test value` {
		t.Errorf("Value was read from the database but is is different form expected `test value` != %s", value)
	}
}

func TestRestoreTrigrammes(t *testing.T) {
	store.DumpTrigrammes(testData)
	trigrammes, err := store.RestoreTrigrammes("en")
	if err != nil {
		t.Errorf(`Error restoring trigrammes: %v`, err)
	}
	if len(trigrammes) != len(testData["en"]) {
		t.Errorf("Error restoring trigrammes: expected length %d, got %d", len(testData["en"]), len(trigrammes))
	}
}

func TestDumpTrigrammes(t *testing.T) {
	store.DumpTrigrammes(testData)
}
