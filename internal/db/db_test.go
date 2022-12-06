package db

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/deanfoley/netspeak-go-test/internal/testutils"
)

func Test_TaskTestScenario(t *testing.T) {
	usernames := make([]string, 0)
	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	for i := 0; i < 6; i++ {
		username := testutils.RandomString(12)
		CreateNewUser(username, resultChan, errorChan)

		<-resultChan

		usernames = append(usernames, username)
	}

	testUser := testutils.RandomString(12)
	CreateNewUser(testUser, resultChan, errorChan)

	<-resultChan

	for i := 0; i < 5; i++ {
		FollowUser(testUser, usernames[i], resultChan, errorChan)
		<-resultChan
	}

	for i := 0; i < 2; i++ {
		FollowUser(usernames[i], testUser, resultChan, errorChan)
		<-resultChan
	}

	close(resultChan)
	close(errorChan)

	testUserData := userDb[testUser]
	var followingUsernames []string
	for _, value := range testUserData.Following {
		followingUsernames = append(followingUsernames, value.UserName)
	}
	var followersUsernames []string
	for _, value := range testUserData.Followers {
		followersUsernames = append(followersUsernames, value.UserName)
	}

	if !reflect.DeepEqual(usernames[:5], followingUsernames) {
		t.Fatal("following lists did not match")
	}
	if !reflect.DeepEqual(usernames[:2], followersUsernames) {
		t.Fatal("follower lists did not match")
	}
}

func Test_CreateNewUser(t *testing.T) {
	username := testutils.RandomString(12)
	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	CreateNewUser(username, resultChan, errorChan)

	<-resultChan

	_, ok := userDb[username]
	if !ok {
		t.Fatal("User was not properly stored.")
	}
}

// This completely deadlocks
// Which has poor implications as to how my system might run under load
func Benchmark_CreateNew(b *testing.B) {
	b.StopTimer()

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	go func() {
		select {
		case <-resultChan:
			<-resultChan
			fmt.Println("ding")
		case <-errorChan:
			err := <-errorChan
			fmt.Println(err)
		}
	}()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		CreateNewUser(strconv.Itoa(n), resultChan, errorChan)
	}
}
