package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/deanfoley/netspeak-go-test/internal/db"
)

// POST /CreateUser/
// Accepts a payload with a desired username
// Returns confirmation if successful (username isn't taken), error if not
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "Incorrect REST method: %s, please use POST.", r.Method)
		return
	}

	userName := strings.TrimPrefix(r.URL.Path, "/createNewUser/")

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	go db.CreateNewUser(userName, resultChan, errorChan)

	select {
	case <-resultChan:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User %s created!", userName)))
	case err := <-errorChan:
		http.Error(w, err.Error(), 500)
	}

}
