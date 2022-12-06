package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/deanfoley/netspeak-go-test/internal/data"
	"github.com/deanfoley/netspeak-go-test/internal/db"
	"github.com/deanfoley/netspeak-go-test/internal/testutils"
)

func TestMain(m *testing.M) {
	stubMux := http.NewServeMux()

	stubMux.HandleFunc("/followers/", GetFollowers)
	stubMux.HandleFunc("/following/", GetFollowing)
	stubMux.HandleFunc("/createNewUser/", CreateNewUser)
	stubMux.HandleFunc("/follow/", FollowUser)
	stubMux.HandleFunc("/unfollow/", UnfollowUser)

	os.Exit(m.Run())
}

// Test scenario from the exercise as an HTTP test
func Test_TaskTestScenario(t *testing.T) {
	usernames := make([]string, 0)
	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	for i := 0; i < 6; i++ {
		username := testutils.RandomString(12)
		db.CreateNewUser(username, resultChan, errorChan)

		<-resultChan

		usernames = append(usernames, username)
	}

	testUser := testutils.RandomString(12)
	db.CreateNewUser(testUser, resultChan, errorChan)

	<-resultChan

	for i := 0; i < 5; i++ {
		db.FollowUser(testUser, usernames[i], resultChan, errorChan)
		<-resultChan
	}

	for i := 0; i < 2; i++ {
		db.FollowUser(usernames[i], testUser, resultChan, errorChan)
		<-resultChan
	}

	close(resultChan)
	close(errorChan)

	getFollowingRequest, err := http.NewRequest("GET", fmt.Sprintf("/following/%s", testUser), nil)
	if err != nil {
		t.Fatal(err)
	}

	getFollowingWriter := httptest.NewRecorder()

	GetFollowing(getFollowingWriter, getFollowingRequest)

	getFollowingResponse := getFollowingWriter.Result()
	defer getFollowingResponse.Body.Close()

	if getFollowingResponse.StatusCode != http.StatusOK {
		t.Fatalf("Bad response from handler: %v", getFollowingResponse.StatusCode)
	}

	output, err := ioutil.ReadAll(getFollowingResponse.Body)
	if err != nil {
		t.Fatal(err)
	}
	var dat data.Response
	err = json.Unmarshal(output, &dat)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dat.Users, usernames[:5]) {
		t.Fatal("slices don't match!")
	}

	if testutils.Contains(dat.Users, usernames[5]) {
		t.Fatal("slices don't match!")
	}
}

func Test_CreateNewUser(t *testing.T) {
	username := testutils.RandomString(12)
	createNewUserRequest, err := http.NewRequest("POST", fmt.Sprintf("/createNewUser/%s", username), nil)
	if err != nil {
		t.Fatal(err)
	}

	createNewUserWriter := httptest.NewRecorder()

	CreateNewUser(createNewUserWriter, createNewUserRequest)

	createNewUserResponse := createNewUserWriter.Result()
	defer createNewUserResponse.Body.Close()

	if createNewUserResponse.StatusCode != http.StatusOK {
		t.Fatalf("Bad response from handler: %v", createNewUserResponse.StatusCode)
	}
}

func Test_GetFollowing(t *testing.T) {
	username := testutils.RandomString(12)

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	db.CreateNewUser(username, resultChan, errorChan)

	<-resultChan

	close(resultChan)
	close(errorChan)

	getFollowingRequest, err := http.NewRequest("GET", fmt.Sprintf("/following/%s", username), nil)
	if err != nil {
		t.Fatal(err)
	}

	getFollowingWriter := httptest.NewRecorder()

	GetFollowing(getFollowingWriter, getFollowingRequest)

	getFollowingResponse := getFollowingWriter.Result()
	defer getFollowingResponse.Body.Close()

	if getFollowingResponse.StatusCode != http.StatusOK {
		t.Fatalf("Bad response from handler: %v", getFollowingResponse.StatusCode)
	}
}

func Test_GetFollowingError(t *testing.T) {
	username := testutils.RandomString(12)

	getFollowingRequest, err := http.NewRequest("GET", fmt.Sprintf("/following/%s", username), nil)
	if err != nil {
		t.Fatal(err)
	}

	getFollowingErrorWriter := httptest.NewRecorder()

	GetFollowing(getFollowingErrorWriter, getFollowingRequest)

	getFollowingResponse := getFollowingErrorWriter.Result()
	defer getFollowingResponse.Body.Close()

	if getFollowingResponse.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Bad response from handler: %v", getFollowingResponse.StatusCode)
	}
}

func Test_GetFollowers(t *testing.T) {
	username := testutils.RandomString(12)

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	db.CreateNewUser(username, resultChan, errorChan)

	<-resultChan

	close(resultChan)
	close(errorChan)

	getFollowersRequest, err := http.NewRequest("GET", fmt.Sprintf("/followers/%s", username), nil)
	if err != nil {
		t.Fatal(err)
	}

	getFollowersWriter := httptest.NewRecorder()

	GetFollowers(getFollowersWriter, getFollowersRequest)

	getFollowersResponse := getFollowersWriter.Result()
	defer getFollowersResponse.Body.Close()

	if getFollowersResponse.StatusCode != http.StatusOK {
		t.Fatalf("Bad response from handler: %v", getFollowersResponse.StatusCode)
	}
}

func Test_GetFollowersError(t *testing.T) {
	username := testutils.RandomString(12)

	getFollowersErrorRequest, err := http.NewRequest("GET", fmt.Sprintf("/followers/%s", username), nil)
	if err != nil {
		t.Fatal(err)
	}

	getFollowersErrorWriter := httptest.NewRecorder()

	GetFollowers(getFollowersErrorWriter, getFollowersErrorRequest)

	getFollowersErrorResponse := getFollowersErrorWriter.Result()
	defer getFollowersErrorResponse.Body.Close()

	if getFollowersErrorResponse.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Bad response from handler: %v", getFollowersErrorResponse.StatusCode)
	}
}

func Test_FollowUser(t *testing.T) {
	user1 := testutils.RandomString(12)

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	db.CreateNewUser(user1, resultChan, errorChan)

	<-resultChan

	user2 := testutils.RandomString(12)

	db.CreateNewUser(user2, resultChan, errorChan)

	<-resultChan

	close(resultChan)
	close(errorChan)

	requestData := data.Request{
		Follower: user1,
		Target:   user2,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	followRequest, err := http.NewRequest("PATCH", "/follow/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	followWriter := httptest.NewRecorder()

	FollowUser(followWriter, followRequest)

	followUserResponse := followWriter.Result()
	defer followUserResponse.Body.Close()

	if followUserResponse.StatusCode != http.StatusOK {
		t.Fatalf("Bad response from handler: %v", followUserResponse.StatusCode)
	}
}

func Test_FollowUserError(t *testing.T) {
	user1 := testutils.RandomString(12)

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	db.CreateNewUser(user1, resultChan, errorChan)

	<-resultChan

	user2 := testutils.RandomString(12)

	close(resultChan)
	close(errorChan)

	requestData := data.Request{
		Follower: user1,
		Target:   user2,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	followRequest, err := http.NewRequest("PATCH", "/follow/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	followWriter := httptest.NewRecorder()

	FollowUser(followWriter, followRequest)

	followUserResponse := followWriter.Result()
	defer followUserResponse.Body.Close()

	if followUserResponse.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Bad response from handler: %v", followUserResponse.StatusCode)
	}
}

func Test_UnfollowUser(t *testing.T) {
	user1 := testutils.RandomString(12)

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	db.CreateNewUser(user1, resultChan, errorChan)

	<-resultChan

	user2 := testutils.RandomString(12)

	close(resultChan)
	close(errorChan)

	requestData := data.Request{
		Follower: user1,
		Target:   user2,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	followRequest, err := http.NewRequest("PATCH", "/unfollow/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	unfollowWriter := httptest.NewRecorder()

	UnfollowUser(unfollowWriter, followRequest)

	unfollowUserResponse := unfollowWriter.Result()
	defer unfollowUserResponse.Body.Close()

	if unfollowUserResponse.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Bad response from handler: %v", unfollowUserResponse.StatusCode)
	}
}

func Test_UnfollowUserError(t *testing.T) {
	user1 := testutils.RandomString(12)

	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	db.CreateNewUser(user1, resultChan, errorChan)

	<-resultChan

	user2 := testutils.RandomString(12)

	db.CreateNewUser(user2, resultChan, errorChan)

	<-resultChan

	db.FollowUser(user1, user2, resultChan, errorChan)

	<-resultChan

	close(resultChan)
	close(errorChan)

	requestData := data.Request{
		Follower: user1,
		Target:   user2,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	followRequest, err := http.NewRequest("PATCH", "/unfollow/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	unfollowWriter := httptest.NewRecorder()

	UnfollowUser(unfollowWriter, followRequest)

	unfollowUserResponse := unfollowWriter.Result()
	defer unfollowUserResponse.Body.Close()

	if unfollowUserResponse.StatusCode != http.StatusOK {
		t.Fatalf("Bad response from handler: %v", unfollowUserResponse.StatusCode)
	}
}
