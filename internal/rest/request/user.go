package request

type User struct {
	Login                string `json:"login" validate:"required,min=3,max=70,login"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
