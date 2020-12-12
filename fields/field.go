package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
)

type BaseField interface {
	Init(args *InitArgs)
	GetRoot() BaseField
	ToInternalValue(data interface{}) (interface{}, *errors.ValidationError)
	ToRepresentation(interface{}) interface{}
	GetParent() BaseField
	IsPartial() bool
}

type Messages map[string]string

type Context map[string]interface{}

type DefaultCaller func() interface{}

type DefaultMessages interface {
	GetDefaultMessages() *Messages
}

type InitArgs map[string]interface{}

type Empty struct {
	IsEmpty	bool
}

var EmptyValue = Empty{
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

	Validators 		[]Validator
	Parent     		BaseField
	Context    		Context

	BaseFieldStruct 	BaseField
}

func (f *Field) Init(args *InitArgs) {
	defaultArgs := InitArgs{
		"ReadOnly":  false,
		"WriteOnly": false,
		"Required":  false,
		"FieldName": "",
		"AllowNull": false,
		"Label":     "",
		"Initial":   nil,
		"Context":   nil,
		"Default":   EmptyValue,
	}

	var _args InitArgs

	if args != nil {
		_args = *args
	} else {
		_args = InitArgs{}
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

	k, ok := interface{}(f).(DefaultMessages)

	if ok {
		d := k.GetDefaultMessages()

		for k, v := range *d {
			f.ErrorMessages[k] = v
		}
	}
}

func (f *Field) GetParent() BaseField {
	if f.Parent != nil {
		p := f.Parent.(BaseField)

		return p
	}

	return nil
}

func (f *Field) GetRoot() BaseField {
	var r BaseField = f

	for r.GetParent() != nil {
		r = r.GetParent()
	}

	return r
}

func (f *Field) Fail(key string) *errors.ValidationError {
	msg := f.ErrorMessages[key]

	err := errors.ValidationError{
		Detail: msg,
		Code:   key,
	}

	return &err
}

func (f *Field) ToInternalValue(data interface{}) (interface{}, *errors.ValidationError) {
	panic("internal value not implemented")
}

func (f *Field) ToRepresentation(data interface{}) interface{} {
	panic("representation not implemented")
}

func (f *Field) IsPartial() bool {
	return f.Partial
}

func (f *Field) ValidateEmptyValues(data interface{}) (bool, interface{}) {
	if f.ReadOnly {
		return true, f.GetDefault()
	}

	if _, ok := f.Default.(Empty); ok {
		if f.Partial == false {
			return false, errors.SkipField{}
		}

		if f.Required {
			return false, f.Fail("required")
		}

		return true, f.GetDefault()
	}

	if data == nil {
		if f.AllowNull == false {
			return false, f.Fail("null")
		}

		return true, nil
	}

	return false, data
}

func (f *Field) GetDefault() interface{} {
	if _, ok := f.Default.(Empty); ok || f.Parent.IsPartial() == false {
		return errors.SkipField{}
	}

	if d, ok := f.Default.(DefaultCaller); ok {
		return d()
	}

	return f.Default
}

func (f *Field) RunValidation(data interface{}) (interface{}, error) {
	isEmptyValue, data := f.ValidateEmptyValues(data)

	if isEmptyValue {
		return data, nil
	}

	value, err := f.ToInternalValue(data)

	if err != nil {
		return nil, err
	}

	if err := f.RunValidators(data); err != nil {
		return nil, err
	}

	return value, nil
}

func (f *Field) RunValidators(data interface{}) *errors.ValidationError {
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