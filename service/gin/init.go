package gin

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	validatorSelf "github.com/lopysso/server/validator"
)

func InitValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("gin error for validator.engine")
	}

	v.RegisterValidation("SnowflakeInt64", validatorSelf.SnowflakeInt64)
	return nil
}
