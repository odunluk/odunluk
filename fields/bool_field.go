package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
	"reflect"
)

var TRUE_VALUES = []interface{}{
	"t", "T", "y",
	"Y", "yes", "YES",
	"true", "True", "TRUE",
	"on", "On", "ON",
	"1", 1,
	true,
}

var FALSE_VALUES = []interface{}{
	"f", "F",
	"n", "N", "no", "NO",
	"false", "False", "FALSE",
	"off", "Off", "OFF",
	"0", 0, 0.0,
	false,
}

var NULL_VALUES = []interface{}{"null", "Null", "NULL", "", nil}

type BooleanField struct {
	Field
}

func (f *BooleanField) GetDefaultMessages() *Messages {
	return &Messages{
		"invalid": "Invalid boolean value",
	}
}

func (f *BooleanField) ToInternalValue(data interface{}) (interface{}, *errors.ValidationError) {
	for _, trueValue := range TRUE_VALUES {
		if reflect.DeepEqual(trueValue, data) {
			return true, nil
		}
	}

	for _, falseValue := range FALSE_VALUES {
		if reflect.DeepEqual(falseValue, data) {
			return false, nil
		}
	}

	for _, nullValue := range NULL_VALUES {
		if reflect.DeepEqual(nullValue, data) {
			return nil, nil
		}
	}

	return nil, f.Fail("invalid")
}

func (f *BooleanField) ToRepresentation(data interface{}) interface{} {
	for _, trueValue := range TRUE_VALUES {
		if reflect.DeepEqual(trueValue, data) {
			return true
		}
	}

	for _, falseValue := range FALSE_VALUES {
		if reflect.DeepEqual(falseValue, data) {
			return false
		}
	}

	for _, nullValue := range NULL_VALUES {
		if reflect.DeepEqual(nullValue, data) {
			return nil
		}
	}

	return data.(bool)
}