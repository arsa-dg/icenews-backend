package helper

import (
	"icenews/backend/model"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func ErrMsg(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "min":
		return "This field needs more character"
	case "max":
		return "This field needs less character"
	case "uri":
		return "This field must in URI format"
	default:
		return "Error"
	}
}

func RequestValidation(Validator *validator.Validate, request interface{}) (interface{}, int) {
	errValidate := Validator.Struct(request)

	if errValidate != nil {
		log.Error().Err(errValidate).Msg("Incorrect request fields format")

		res := model.ResponseValidationFailed{
			Message: "Field(s) validation error",
		}

		var emptyFields []model.FieldError

		for _, err := range errValidate.(validator.ValidationErrors) {
			toAdd := model.FieldError{
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
