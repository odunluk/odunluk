package fields

import (
	"reflect"
	"testing"
)

func TestRoot(t *testing.T) {
	rootField := &Field{}

	var r BaseField = rootField

	childField := Field{
		parent: r,
	}

	rootOutput := childField.Root()

	if !reflect.DeepEqual(rootField, rootOutput) {
		t.Fail()
	}
}
