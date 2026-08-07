package request

type RegisterUser struct {
	Login                string `json:"login" validate:"required,min=3,max=70,login"`
	Password             string `json:"password" validate:"required,min=6,max=65"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginUser struct {
	Login    string `json:"login" validate:"required,min=3,max=70,login"`
	Password string `json:"password" validate:"required,min=6,max=65"`
}
