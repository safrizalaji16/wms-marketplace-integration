package dto

type UserData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
