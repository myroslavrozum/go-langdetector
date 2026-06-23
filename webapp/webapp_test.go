package webapp_test

import (
	"go-langdetector/db"
	"go-langdetector/webapp"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		store  *db.Store
		logger chan string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webapp.Run(tt.store, tt.logger)
		})
	}
}
