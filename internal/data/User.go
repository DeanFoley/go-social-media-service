package data

type User struct {
	UserName  string
	Followers []*User
	Following []*User
}

func (u *User) AddNewFollower(follower *User) {
	u.Followers = append(u.Followers, follower)
	follower.Following = append(follower.Following, u)
}

func (u *User) RemoveFollower(follower *User) {
	followerIndex := indexOf(follower.UserName, u.Followers)
	if followerIndex > -1 {
		u.Followers = append(u.Followers[:followerIndex], u.Followers[followerIndex+1:]...)
	}
	followingIndex := indexOf(u.UserName, follower.Following)
	if followingIndex > -1 {
		follower.Following = append(follower.Following[:followingIndex], follower.Following[followingIndex+1:]...)
	}
}

func indexOf(username string, list []*User) int {
	for k, v := range list {
		if username == v.UserName {
			return k
		}
	}
	return -1
}
