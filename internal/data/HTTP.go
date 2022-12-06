package data

type Request struct {
	Follower string
	Target   string
}

type Response struct {
	UserName    string   `json:"userName"`
	RequestType string   `json:"requestType"`
	Users       []string `json:"users"`
}
