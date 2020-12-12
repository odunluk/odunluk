package fields

import (
	"fmt"
	"github.com/odunluk/odunluk/v1/errors"
	"reflect"
	"strings"
)

type CharField struct {
	Field

	allowBlank		bool
	trimWhiteSpace	bool
	maxLength		*int
	minLength		*int
}

func (f *CharField) Init(args *InitArgs) {
	f.Field.Init(args)

	defaultArgs := InitArgs{
		"allowBlank": false,
		"trimWhiteSpace": true,
		"maxLength": nil,
		"minLength": nil,
	}

	_args := *args

	for k, v := range defaultArgs {
		_args[k] = v
	}

	f.AllowNull = _args["AllowNull"].(bool)
	f.trimWhiteSpace = _args["trimWhiteSpace"].(bool)
	f.maxLength = _args["maxLength"].(*int)
	f.minLength = _args["minLength"].(*int)
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

	var str interface{}= fmt.Sprintf("%v", data)

	if f.trimWhiteSpace {
		str = strings.TrimSpace(str.(string))
	}

	return str, nil
}

