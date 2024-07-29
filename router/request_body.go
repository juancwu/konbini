package router

// reqBodyValidationMsgs is a collection of all request body validation error messages.
var reqBodyValidationMsgs = map[string]string{
	// signupReqBody messages
	"signupReqBody.Email.required":    "Missing email in request body",
	"signupReqBody.Email.email":       "Invalid email provided",
	"signupReqBody.Password.required": "Missing password in request body",
	"signupReqBody.Password.password": "Invalid password format. A password must of at least 12 characters and combines at least one special character, uppercase and lowercase letter, and number",
	"signupReqBody.Name.required":     "Missing name in request body",
	"signupReqBody.Name.min":          "Name must be at least 3 characters long",
	"signupReqBody.Name.max":          "Name must not be longer than 30 characters",
	"signupReqBody.Name.alpha":        "Name must only include alphabet characters",

	// signinReqBody messages
	"signinReqBody.Email.required":    "Missing email in request body",
	"signinReqBody.Email.email":       "Invalid email provided",
	"signinReqBody.Password.required": "Missing password in request body",

	// resendVerificationEmailReqBody messages
	"resendVerificationEmailReqBody.Email.required": "Missing email in request body",
	"resendVerificationEmailReqBody.Email.email":    "Invalid email provided",

	// new bento messages
	"newBentoReqBody.Name.required":               "Missing name in request body",
	"newBentoReqBody.Name.min":                    "Bento body name must longer than 3 and shorter than 15 characters",
	"newBentoReqBody.Name.max":                    "Bento body name must longer than 3 and shorter than 15 characters",
	"newBentoReqBody.Name.ascii":                  "Bento name can only contain printable ASCII characters",
	"newBentoReqBody.PubKey.required":             "Missing pub_key missing in request body",
	"newBentoReqBody.Ingridients.Name.required":   "Missing ingridient name.",
	"newBentoReqBody.Ingridients.Name.printascii": "Ingridient name can only contain printable ascii characters.",
	"newBentoReqBody.Ingridients.Value.required":  "Missing ingridient value.",

	// add ingridients messages
	"addIngridientsReqBody.Ingridients.required":        "Missing ingridients in request body.",
	"addIngridientsReqBody.Ingridients.gt":              "You must provide at least one Ingridient.",
	"addIngridientsReqBody.BentoId.required":            "Missing bento id in request body.",
	"addIngridientsReqBody.BentoId.uuid":                "Invalid uuid for bento id.",
	"addIngridientsReqBody.Challenge.required":          "Missing challenge in request body.",
	"addIngridientsReqBody.Signature.required":          "Missing signature in request body.",
	"addIngridientsReqBody.Ingridients.Name.required":   "Missing ingridient name.",
	"addIngridientsReqBody.Ingridients.Name.printascii": "Ingridient name can only contain printable ascii characters.",
	"addIngridientsReqBody.Ingridients.Value.required":  "Missing ingridient value.",

	// rename bento messages
	"renameBentoReqBody.BentoId.required": "Missing required 'bento_id' in request body.",
	"renameBentoReqBody.BentoId.uuid4":    "Invalid bento id. It must be a UUID v4.",
	"renameBentoReqBody.NewName.required": "Missing required 'new_name' in request body.",
	"renameBentoReqBody.NewName.min":      "New name must be at least 3 characters long.",
	"renameBentoReqBody.NewName.max":      "New name must not exceed 50 characters.",
	"renameBentoReqBody.NewName.ascii":    "New name must only contain ascii characters.",
}

// signupReqBody represents the request body that is expected when handling a signup request.
type signupReqBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Name     string `json:"name" validate:"required,min=3,max=30,alpha"`
}

// signinReqBody represents the request body that is expected when handling a signin request.
type signinReqBody struct {
	// Email is the username that is used to signin
	Email string `json:"email" validate:"required,email"`
	// Password no password validation here because there is no need to gives hints about it when signing in.
	Password string `json:"password" validate:"required"`
}

// resendVerificationEmailReqBody represents the request body that is expected when handling a resend ev request.
type resendVerificationEmailReqBody struct {
	Email string `json:"email" validate:"required,email"`
}

// newBentoReqBody represents the request body that is expected when handling a new bento requets.
type newBentoReqBody struct {
	Name        string       `json:"name" validate:"required,min=3,max=50,ascii"`
	PubKey      string       `json:"pub_key" validate:"required"`
	Ingridients []Ingridient `json:"ingridients,omitempty" validate:"omitnil,dive"`
}

// addIngridientsReqBody represents the request body that is expected when handling adding a new Ingridient to a prepared bento.
type addIngridientsReqBody struct {
	BentoId     string       `json:"bento_id" validate:"required,uuid4"`
	Ingridients []Ingridient `json:"ingridients" validate:"required,gt=0,dive"`
	Challenge   string       `json:"challenge" validate:"required"`
	Signature   string       `json:"signature" validate:"required"`
}

// Ingridient is used in the addIngridientsReqBody.
type Ingridient struct {
	Name  string `json:"name" validate:"required,printascii"`
	Value string `json:"value" validate:"required"`
}

// renameBentoReqBody represents the request body that is expected when handling rename bento requests.
type renameBentoReqBody struct {
	BentoId string `json:"bento_id" validate:"required,uuid4"`
	NewName string `json:"new_name" validate:"required,min=3,max=50,ascii"`
}
