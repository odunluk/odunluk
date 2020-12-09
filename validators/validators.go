package validators

import (
	"github.com/odunluk/odunluk/v1/errors"
	"github.com/odunluk/odunluk/v1/fields"
)

type Validator interface {
	Validate(interface{}, fields.BaseField) *errors.ValidationError
}