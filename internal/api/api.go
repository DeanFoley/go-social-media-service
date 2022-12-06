package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/deanfoley/netspeak-go-test/internal/data"
)

func Launch(readyCheck chan struct{}, shutdown chan struct{}, closeAnnounce chan struct{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("/followers/", GetFollowers)
	mux.HandleFunc("/following/", GetFollowing)
	mux.HandleFunc("/createNewUser/", CreateNewUser)
	mux.HandleFunc("/follow/", FollowUser)
	mux.HandleFunc("/unfollow/", UnfollowUser)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("Server started")

	readyCheck <- struct{}{}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	go func() {
		<-shutdown
		if err := srv.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("server Shutdown Failed:%+s", err)
		}
		log.Printf("Server exited properly")
		close(closeAnnounce)
	}()
}

func ParseToJSON(username string, operation string, userNames []string) []byte {
	response := data.Response{
		UserName:    username,
		RequestType: operation,
		Users:       userNames,
	}
	json, _ := json.Marshal(response)
	return json
}
