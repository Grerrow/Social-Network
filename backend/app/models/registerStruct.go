package models

type RegisterStruct struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	// avatar is an image
	Avatar    string `json:"avatar"`
	Username  string `json:"username"`
	AboutMe   string `json:"about_me"`
	IsPrivate bool   `json:"is_private"`
}
