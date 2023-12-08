package models

type Address struct {
	State   string `json:"state" bson:"state"`
	City    string `json:"city" bson:"city"`
	Pincode int    `json:"pincode" bson:"pincode"`
}

type User struct {
	ID       string  `json:"id" bson:"user_id"`
	Password string  `json:"password" bson:"user_password"`
	Name     string  `json:"name" bson:"user_name"`
	Age      int     `json:"age" bson:"user_age"`
	Address  Address `json:"address" bson:"user_address"`
}

type GoogleUser struct {
	ID      string `json:"id" bson:"google_user_id"`
	Email   string `json:"email" bson:"google_user_email"`
	Name    string `json:"name" bson:"google_user_name"`
	Picture string `json:"picture" bson:"google_user_picture"`
}
