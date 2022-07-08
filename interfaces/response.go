package interfaces

type AuthLoginResponse struct {
	Token      string `json:"token"`
	Scheme     string `json:"scheme"`
	Expires_at string `json:"expires_at"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseBadRequest struct {
	Message string `json:"message"`
}

type ResponseUnauthorized struct {
	Message string `json:"message"`
}

type FieldError struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

type ResponseValidationFailed struct {
	Message string       `json:"message"`
	Field   []FieldError `json:"field"`
}

type ResponseInternalServerError struct {
	Message string `json:"message"`
}
