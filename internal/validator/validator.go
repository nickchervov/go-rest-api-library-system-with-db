package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

// Регулярное выражение:
// \d — разрешает цифры
// \p{P} — разрешает любую пунктуацию (-, ?, !, @, #, $, %, etc.)
// \p{S} — разрешает математические и валютные символы (+, =, <, $, etc.)
// ^...+$ — строка должна состоять только из этих символы от начала до конца
var digitsAndSpecsRegex = regexp.MustCompile(`^[\d\p{P}\p{S}]+$`)

var alphaSpaceRegex = regexp.MustCompile(`^[\p{L}\s]+$`)

func AlphaSpaceValidator(fl validator.FieldLevel) bool {
	return alphaSpaceRegex.MatchString(fl.Field().String())
}

func DigitsAndSpecsValidator(fl validator.FieldLevel) bool {
	return digitsAndSpecsRegex.MatchString(fl.Field().String())
}

func init() {
	Validator = validator.New()

	Validator.RegisterValidation("alphaspace", AlphaSpaceValidator)
	Validator.RegisterValidation("numericspecs", DigitsAndSpecsValidator)
}
