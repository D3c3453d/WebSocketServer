package entity

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username,omitempty" db:"username"`
	Password string `json:"password,omitempty" db:"password"`
}

type FriendCheck struct {
	User   string `json:"user,omitempty" db:"user_name"`
	Friend string `json:"friend,omitempty" db:"friend_name"`
}
