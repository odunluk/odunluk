package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
	"reflect"
	"testing"
)

func TestBoolFieldRoot(t *testing.T) {
	var rootField BaseField = &BooleanField{}

	childField := BooleanField{}

	childField.Parent = rootField

	rootOutput := childField.getRoot()

	if !reflect.DeepEqual(rootField, rootOutput) {
		t.Fail()
	}
}

func TestBoolFieldToInternalValue(t *testing.T) {
	field := &BooleanField{}

	v, _ := field.toInternalValue("true")

	if v.(bool) != true {
		t.Fail()
	}

	nilValue, err := field.toInternalValue("tT")

	if nilValue != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(reflect.TypeOf(err), reflect.TypeOf(&errors.ValidationError{})) {
		t.Fail()
	}
}

func TestBoolFieldToRepresentation(t *testing.T) {
	field := &BooleanField{}

	v := field.toRepresentation(true)

	if v.(bool) != true {
		t.Fail()
	}

	v = field.toRepresentation("false")

	if v.(bool) != false {
		t.Fail()
	}
}

func TestDefaultMessages(t *testing.T) {
	field := &BooleanField{}

	var bf BaseField = &BooleanField{}

	field.BaseFieldStruct = bf

	field.init(nil)

	messages := Messages{
		"invalid": "Invalid boolean value",
	}

	if !reflect.DeepEqual(&messages, field.getDefaultMessages()) {
		t.Fail()
	}
}
