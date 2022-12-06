package db

import (
	"fmt"

	"github.com/deanfoley/netspeak-go-test/internal/data"
)

var userDb map[string]*data.User

var cmds chan request

type request struct {
	instruction string
	userName    string
	target      string
	resultChan  chan struct{}
	dataChan    chan []*data.User
	errorChan   chan error
}

func init() {
	userDb = make(map[string]*data.User)

	cmds = make(chan request, 100)

	go func(userDb map[string]*data.User, cmds chan request) {
		for cmd := range cmds {
			switch cmd.instruction {
			case "getFollowers":
				if val, ok := userDb[cmd.userName]; ok {
					cmd.dataChan <- val.Followers
				} else {
					cmd.errorChan <- fmt.Errorf("error from database: user %s not found", cmd.userName)
				}
			case "getFollowing":
				if val, ok := userDb[cmd.userName]; ok {
					cmd.dataChan <- val.Following
				} else {
					cmd.errorChan <- fmt.Errorf("error from database: user %s not found", cmd.userName)
				}
			case "createNewUser":
				if _, ok := userDb[cmd.userName]; !ok {
					userDb[cmd.userName] = &data.User{
						UserName:  cmd.userName,
						Followers: make([]*data.User, 0),
						Following: make([]*data.User, 0),
					}
					cmd.resultChan <- struct{}{}
				} else {
					cmd.errorChan <- fmt.Errorf("error from database: user %s already exists", cmd.userName)
				}
			case "followUser":
				if val1, ok := userDb[cmd.userName]; ok {
					if _, ok := userDb[cmd.target]; ok {
						// val.Following = append(val.Following, &val2)
						userDb[cmd.target].AddNewFollower(val1)
						cmd.resultChan <- struct{}{}
					} else {
						cmd.errorChan <- fmt.Errorf("error from database: user %s not found", cmd.target)
					}
				} else {
					cmd.errorChan <- fmt.Errorf("error from database: user %s not found", cmd.userName)
				}
			case "unfollowUser":
				if val, ok := userDb[cmd.userName]; ok {
					if val2, ok := userDb[cmd.target]; ok {
						val2.RemoveFollower(val)
						cmd.resultChan <- struct{}{}
					} else {
						cmd.errorChan <- fmt.Errorf("error from database: user %s not found", cmd.target)
					}
				} else {
					cmd.errorChan <- fmt.Errorf("error from database: user %s not found", cmd.userName)
				}
			}
		}
	}(userDb, cmds)
}

func GetFollowers(userName string, result chan []*data.User, errors chan error) {
	replyChan := make(chan []*data.User, 1)
	errorChan := make(chan error, 1)
	request := request{
		instruction: "getFollowers",
		userName:    userName,
		dataChan:    replyChan,
		errorChan:   errorChan,
	}
	cmds <- request
	select {
	case val := <-replyChan:
		result <- val
	case err := <-errorChan:
		errors <- err
	}
	close(replyChan)
	close(errorChan)
}

func GetFollowing(userName string, result chan []*data.User, errors chan error) {
	replyChan := make(chan []*data.User, 1)
	errorChan := make(chan error, 1)
	request := request{
		instruction: "getFollowing",
		userName:    userName,
		dataChan:    replyChan,
		errorChan:   errorChan,
	}
	cmds <- request
	select {
	case val := <-replyChan:
		result <- val
	case err := <-errorChan:
		errors <- err
	}
	close(replyChan)
	close(errorChan)
}

func CreateNewUser(userName string, resultChan chan struct{}, errorChan chan error) {
	request := request{
		instruction: "createNewUser",
		userName:    userName,
		resultChan:  resultChan,
		errorChan:   errorChan,
	}
	cmds <- request
	select {
	case <-resultChan:
		resultChan <- struct{}{}
	case err := <-errorChan:
		errorChan <- err
	}
}

func FollowUser(user string, target string, result chan struct{}, errors chan error) {
	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	request := request{
		instruction: "followUser",
		userName:    user,
		target:      target,
		resultChan:  resultChan,
		errorChan:   errorChan,
	}
	cmds <- request
	select {
	case <-resultChan:
		result <- struct{}{}
	case err := <-errorChan:
		errors <- err
	}
	close(resultChan)
	close(errorChan)
}

func UnfollowUser(user string, target string, result chan struct{}, errors chan error) {
	resultChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	request := request{
		instruction: "unfollowUser",
		userName:    user,
		target:      target,
		resultChan:  resultChan,
		errorChan:   errorChan,
	}
	cmds <- request
	select {
	case <-resultChan:
		result <- struct{}{}
	case err := <-errorChan:
		errors <- err
	}
	close(resultChan)
	close(errorChan)
}
