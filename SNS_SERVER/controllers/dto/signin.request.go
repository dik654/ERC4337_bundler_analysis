package dto

type SignInRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
