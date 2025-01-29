package models

type User struct {
	/*
		Не использовал ID, так как проект работает на мапах, где получаем и отправляем юзернеймы.
		Без использования реляционной БД это придется делать новой MAPой
	*/
	//ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Created  int64  `json:"created"`
	Updated  uint64 `json:"updated"`
	Bio      string `json:"bio"`
	// [username]
	Follows map[string]struct{}
}
