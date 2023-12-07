package dto

type SignInRequest struct {
	ID       string `json:"id" bson:"user_id"`
	Password string `json:"password" bson:"user_password"`
}
