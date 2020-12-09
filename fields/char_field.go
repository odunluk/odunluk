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

func (f *CharField) init(args *initArgs) {
	f.Field.init(args)

	defaultArgs := initArgs{
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

func (f *CharField) getDefaultMessages() *Messages {
	return &Messages{
		"invalid": "Invalid string value",
	}
}

func (f *CharField) toInternalValue(data interface{}) (interface{}, *errors.ValidationError) {
	if reflect.ValueOf(data).Kind() == reflect.Bool {
		return nil, f.fail("invalid")
	}

	var str interface{}= fmt.Sprintf("%v", data)

	if f.trimWhiteSpace {
		str = strings.TrimSpace(str.(string))
	}

	return str, nil
}

