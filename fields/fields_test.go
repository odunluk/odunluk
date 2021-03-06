package fields

import (
	"reflect"
	"testing"
)

func TestRoot(t *testing.T) {
	rootField := &Field{}

	var r BaseField = rootField

	childField := Field{
		Parent: r,
	}

	rootOutput := childField.GetRoot()

	if !reflect.DeepEqual(rootField, rootOutput) {
		t.Fail()
	}
}
