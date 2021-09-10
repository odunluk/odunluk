package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
	"reflect"
	"testing"
)

func TestCharField_ToInternalValue(t *testing.T) {
	f := &CharField{}

	v, _ := f.ToInternalValue("foo")

	if ! reflect.DeepEqual(v.(string), "foo") {
		t.Fail()
	}

	nilValue, err := f.ToInternalValue(true)

	if nilValue != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(reflect.TypeOf(err), reflect.TypeOf(&errors.ValidationError{})) {
		t.Fail()
	}
}
