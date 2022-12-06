package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/deanfoley/netspeak-go-test/internal/app"
	"github.com/deanfoley/netspeak-go-test/internal/data"
	"github.com/deanfoley/netspeak-go-test/internal/db"
)

// GET /following/{username}
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "Incorrect REST method: %s, please use GET.", r.Method)
		return
	}

	userName := strings.TrimPrefix(r.URL.Path, "/following/")

	resultChan := make(chan []*data.User, 1)
	errorChan := make(chan error, 1)

	go db.GetFollowing(userName, resultChan, errorChan)

	select {
	case result := <-resultChan:
		following := app.ParseUserList(result)
		w.WriteHeader(200)
		w.Write(ParseToJSON(userName, "following", following))
	case err := <-errorChan:
		http.Error(w, err.Error(), 500)
	}
}
