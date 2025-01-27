package models

type User struct {
	//ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Created   int64  `json:"created"`
	Updated   uint64 `json:"updated"`
	ProfileID uint64 `json:"profileID"`
}
