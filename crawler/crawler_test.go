package crawler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

var TestHTTPResponse = "pong"

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s", TestHTTPResponse)
}

func startServer() {
	http.HandleFunc("/ping", hello)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setup() {
	go startServer()
}

func teardown() {}

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

func TestGetTextFromURL(t *testing.T) {
	text, err := GetTextFromURL("http://127.0.0.1:8080/ping")
	if err != nil {
		t.Error(err)
	}
	if strings.TrimSpace(text) != TestHTTPResponse {
		t.Errorf("Expected %s, got %s", TestHTTPResponse, text)
	}
}
