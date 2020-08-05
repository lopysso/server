package validator

import (
	"regexp"

	validator2 "github.com/go-playground/validator/v10"
)

var SnowflakeInt64 validator2.Func = func(fl validator2.FieldLevel) bool {

	matched, err := regexp.MatchString(`^\d{19}$`, fl.Field().String())
	if err != nil {
		return false
	}

	return matched
}
