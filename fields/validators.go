package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
)

type Validator interface {
	Validate(interface{}, BaseField) *errors.ValidationError
}