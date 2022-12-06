package app

import "github.com/deanfoley/netspeak-go-test/internal/data"

// Accepts a slice of pointers to Users
// Returns a slice of strings of their usernames
func ParseUserList(users []*data.User) []string {
	listedUsers := make([]string, 0)

	for _, user := range users {
		listedUsers = append(listedUsers, user.UserName)
	}

	return listedUsers
}
