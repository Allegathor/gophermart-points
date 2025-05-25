package entity

type User struct {
	Id    int    `json:"id,omitempty"`
	Login string `json:"login"`
	Pwd   string `json:"password"`
}
