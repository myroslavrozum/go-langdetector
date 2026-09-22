# BUG-015: Badger Transaction Leak in `DumpTrigrammes` on Error

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `db`
- **Affects:** [`db/db.go:98-123`](file:///Users/myroslavrozum/src/go-langdetector/db/db.go#L98-L123)

## Description
In `db/db.go`:
```go
func (s *Store) DumpTrigrammes(data map[string]map[string]float64) error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard() // Automatically rolled back if not committed

	for language, trigrammes := range data {
		...
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
```
In Go, `defer txn.Discard()` evaluates the receiver expression `txn` at the moment the `defer` statement is executed. When `txn = s.db.NewTransaction(true)` reassigns `txn`, the deferred function still holds a reference to the first (already committed) transaction.

## Impact
If `txn.Set` fails on line 115 or line 118 after a split, the second transaction is never discarded. In BadgerDB, uncommitted and undiscarded write transactions hold memory and locks, leading to resource leaks.

## Suggested Fix
Use a closure so the defer evaluates the current value of `txn`, or handle transactions per item:
```go
txn := s.db.NewTransaction(true)
defer func() {
	txn.Discard()
}()
```
Or utilize `s.db.Update(...)` per language record.
