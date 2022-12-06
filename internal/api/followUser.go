package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deanfoley/netspeak-go-test/internal/data"
	"github.com/deanfoley/netspeak-go-test/internal/db"
)

// PATCH /follow/
func FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		fmt.Fprintf(w, "Incorrect REST method: %s, please use PATCH.", r.Method)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	var request data.Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("could not parse body: %s\n", err)
	}

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	go db.FollowUser(request.Follower, request.Target, resultChan, errorChan)

	select {
	case <-resultChan:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User %s is now following %s!", request.Follower, request.Target)))
	case err := <-errorChan:
		http.Error(w, err.Error(), 500)
	}

}
