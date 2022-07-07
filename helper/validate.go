package helper

import (
	"icenews/backend/interfaces"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ErrMsg(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "min":
		return "This field needs more character"
	case "uri":
		return "This field must in URI format"
	default:
		return "Error"
	}
}

func RequestValidation(Validator *validator.Validate, request interface{}) (interface{}, int) {
	errValidate := Validator.Struct(request)

	if errValidate != nil {
		res := interfaces.ResponseValidationFailed{
			Message: "Field(s) validation error",
		}

		var emptyFields []interfaces.FieldError

		for _, err := range errValidate.(validator.ValidationErrors) {
			toAdd := interfaces.FieldError{
				Name:  err.Field(),
				Error: ErrMsg(err.Tag()),
			}

			emptyFields = append(emptyFields, toAdd)
		}

		res.Field = emptyFields

		return res, http.StatusUnprocessableEntity
	}

	return nil, 0
}
