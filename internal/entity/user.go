package entity

type User struct {
	ID    int    `json:"id,omitempty"`
	Login string `json:"login"`
	Pwd   string `json:"password"`
}
