package handler

type basicResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type signUpBody struct {
	Email    string  `json:"email" validate:"required,email" errormsg:"required=Email is required,email=Invalid email address"`
	Password string  `json:"password" validate:"required,min=12" errormsg:"required=Password is required,min=Password length must be at least 12"`
	Nickname *string `json:"nickname"`
}

type signInBody struct {
	Email    string `json:"email" validate:"required,email" errormsg:"required=Email is required,email=Invalid email address"`
	Password string `json:"password" validate:"required" errormsg:"Password is required."`
}
