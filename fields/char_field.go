package fields

import (
	"fmt"
	"github.com/odunluk/odunluk/v1/errors"
	"reflect"
	"strings"
)

type CharField struct {
	Field

	AllowBlank     bool
	TrimWhiteSpace bool
	MaxLength      *int
	MinLength      *int
}

func (f *CharField) Init(args *InitArgs) {
	f.Field.Init(args)

	defaultArgs := InitArgs{
		"AllowBlank":     false,
		"TrimWhiteSpace": true,
		"MaxLength":      nil,
		"MinLength":      nil,
	}

	_args := *args

	for k, v := range defaultArgs {
		_args[k] = v
	}

	f.AllowBlank = _args["AllowBlank"].(bool)
	f.TrimWhiteSpace = _args["TrimWhiteSpace"].(bool)
	f.MaxLength = _args["MaxLength"].(*int)
	f.MinLength = _args["MinLength"].(*int)
}

func (f *CharField) GetDefaultMessages() *Messages {
	return &Messages{
		"invalid": "Invalid string value",
	}
}

func (f *CharField) ToInternalValue(data interface{}) (interface{}, *errors.ValidationError) {
	if reflect.ValueOf(data).Kind() == reflect.Bool {
		return nil, f.Fail("invalid")
	}

	var str interface{} = fmt.Sprintf("%v", data)

	if f.TrimWhiteSpace {
		str = strings.TrimSpace(str.(string))
	}

	return str, nil
}
