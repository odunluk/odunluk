package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
	"github.com/odunluk/odunluk/v1/validators"
)

type BaseField interface {
	init(args *initArgs)
	getRoot() BaseField
	toInternalValue(data interface{}) (interface{}, *errors.ValidationError)
	toRepresentation(interface{}) interface{}
	getParent() BaseField
	isPartial() bool
}

type Messages map[string]string

type Context map[string]interface{}

type DefaultCaller func() interface{}

type DefaultMessages interface {
	getDefaultMessages() *Messages
}

type initArgs map[string]interface{}

type Empty struct {
	IsEmpty	bool
}

var empty = Empty{
	IsEmpty: true,
}

type Field struct {
	BaseField

	ReadOnly      	bool
	WriteOnly     	bool
	Required      	bool
	Partial			bool
	FieldName     	string
	AllowNull     	bool
	Label         	string
	ErrorMessages 	Messages

	Initial       	interface{}
	Default       	interface{}

	Validators 		[]validators.Validator
	Parent     		BaseField
	Context    		Context

	BaseFieldStruct 	BaseField
}

func (f *Field) init(args *initArgs) {
	defaultArgs := initArgs{
		"ReadOnly":  false,
		"WriteOnly": false,
		"Required":  false,
		"FieldName": "",
		"AllowNull": false,
		"Label": "",
		"Initial": nil,
		"Context": nil,
		"Default": empty,
	}

	var _args initArgs

	if args != nil {
		_args = *args
	} else {
		_args = initArgs{}
	}

	for k, v := range defaultArgs {
		_args[k] = v
	}

	f.ReadOnly = _args["ReadOnly"].(bool)
	f.WriteOnly = _args["WriteOnly"].(bool)
	f.Required = _args["Required"].(bool)
	f.FieldName = _args["FieldName"].(string)
	f.AllowNull = _args["AllowNull"].(bool)
	f.Label = _args["Label"].(string)

	if _args["Initial"] != nil {
		f.Initial = _args["Initial"]
	}
	if _args["Default"] != nil {
		f.Default = _args["Default"]
	}
	if _args["Context"] != nil {
		f.Context = _args["Context"].(Context)
	}

	if f.ReadOnly && f.WriteOnly {
		panic("Read only, write only same time error")
	}

	if f.ErrorMessages == nil {
		f.ErrorMessages = map[string]string{}
	}

	k, ok := f.BaseFieldStruct.(DefaultMessages)

	if ok {
		d := k.getDefaultMessages()

		for k, v := range *d {
			f.ErrorMessages[k] = v
		}
	}
}

func (f *Field) getParent() BaseField {
	if f.Parent != nil {
		p := f.Parent.(BaseField)

		return p
	}

	return nil
}

func (f *Field) getRoot() BaseField {
	var r BaseField = f

	for r.getParent() != nil {
		r = r.getParent()
	}

	return r
}

func (f *Field) fail(key string) *errors.ValidationError {
	msg := f.ErrorMessages[key]

	err := errors.ValidationError{
		Detail: msg,
		Code:   key,
	}

	return &err
}

func (f *Field) toInternalValue(data interface{}) (interface{}, *errors.ValidationError) {
	panic("internal value not implemented")
}

func (f *Field) toRepresentation(data interface{}) interface{} {
	panic("representation not implemented")
}

func (f *Field) isPartial() bool {
	return f.Partial
}

func (f *Field) validateEmptyValues(data interface{}) (bool, interface{}) {
	if f.ReadOnly {
		return true, f.getDefault()
	}

	if _, ok := f.Default.(Empty); ok {
		if f.Partial == false {
			return false, errors.SkipField{}
		}

		if f.Required {
			return false, f.fail("required")
		}

		return true, f.getDefault()
	}

	if data == nil {
		if f.AllowNull == false {
			return false, f.fail("null")
		}

		return true, nil
	}

	return false, data
}

func (f *Field) getDefault() interface{} {
	if _, ok := f.Default.(Empty); ok || f.Parent.isPartial() == false {
		return errors.SkipField{}
	}

	if d, ok := f.Default.(DefaultCaller); ok {
		return d()
	}

	return f.Default
}

func (f *Field) runValidation(data interface{}) (interface{}, error) {
	isEmptyValue, data := f.validateEmptyValues(data)

	if isEmptyValue {
		return data, nil
	}

	value, err := f.toInternalValue(data)

	if err != nil {
		return nil, err
	}

	if err := f.runValidators(data); err != nil {
		return nil, err
	}

	return value, nil
}

func (f *Field) runValidators(data interface{}) *errors.ValidationError {
	var errs []interface{}

	for _, v := range f.Validators {
		err := v.Validate(data, f)

		if err != nil {
			errs = append(errs, err.Detail)
		}
	}

	if len(errs) > 0 {
		return &errors.ValidationError{
			Detail: errs,
		}
	}

	return nil
}