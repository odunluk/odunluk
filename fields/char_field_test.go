package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
	"reflect"
	"testing"
)

func TestToInternalValue(t *testing.T) {
	f := &CharField{}

	v, _ := f.toInternalValue("foo")

	if ! reflect.DeepEqual(v.(string), "foo") {
		t.Fail()
	}

	nilValue, err := f.toInternalValue(true)

	if nilValue != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(reflect.TypeOf(err), reflect.TypeOf(&errors.ValidationError{})) {
		t.Fail()
	}
}