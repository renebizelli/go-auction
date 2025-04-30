package validations

import (
	"encoding/json"
	"errors"
	"renebizelli/go-leilao/configuration/rest_err"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validator_en "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate = validator.New()
	transl   ut.Translator
)

func init() {

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		enTransl := ut.New(en, en)

		transl, _ = enTransl.GetTranslator("en")

		validator_en.RegisterDefaultTranslations(val, transl)
	}

}

func ValidateErr(validation_err error) *rest_err.RestErr {

	var jsonErr *json.UnmarshalTypeError
	var jsonValidation validator.ValidationErrors

	if ok := errors.As(validation_err, &jsonErr); ok {
		return rest_err.NewBadRequestError("Invalid type error")
	} else if ok := errors.As(validation_err, &jsonValidation); ok {
		errosCauses := []rest_err.Causes{}

		for _, e := range validation_err.(validator.ValidationErrors) {
			erro := rest_err.Causes{
				Field:   e.Field(),
				Message: e.Translate(transl),
			}
			errosCauses = append(errosCauses, erro)
		}

		return rest_err.NewBadRequestError("Validation error", errosCauses...)
	} else {
		return rest_err.InternalServerError("Internal server error")
	}
}
