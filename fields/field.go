package fields

import (
	"github.com/odunluk/odunluk/v1/errors"
	"github.com/odunluk/odunluk/v1/serializers"
)

type BaseField interface {
	Parent() serializers.BaseSerializer
	SetParent(parent serializers.BaseSerializer)

	Validators() []Validator
	SetValidators(validators []Validator)

	Default() interface{}
	SetDefault(_default interface{})

	Initial() interface{}
	SetInitial(initial interface{})

	ErrorMessages() Messages
	SetErrorMessages(errorMessages Messages)

	Label() string
	SetLabel(label string)

	AllowNull() bool
	SetAllowNull(allowNull bool)

	FieldName() string
	SetFieldName(fieldName string)

	Required() bool
	SetRequired(required bool)

	WriteOnly() bool
	SetWriteOnly(writeOnly bool)

	ReadOnly() bool
	SetReadOnly(readOnly bool)

	Init(args *InitArgs)
	SetErrorMessage(key string, message string)

	Root() serializers.BaseSerializer
	Fail(key string) *errors.ValidationError

	ToInternalValue(data interface{}) (interface{}, *errors.ValidationError)
	ToRepresentation(data interface{}) interface{}

	ValidateEmptyValues(data interface{}) (bool, interface{})
	GetDefault() interface{}
	RunValidation(data interface{}) (interface{}, error)
	RunValidators(data interface{}) *errors.ValidationError
}

type Messages map[string]string

type DefaultCaller func() interface{}

type DefaultMessages interface {
	GetDefaultMessages() *Messages
	SetErrorMessage(key string, message string)
}

type InitArgs map[string]interface{}

type Empty struct {
	IsEmpty	bool
}

var EmptyValue = Empty{
	IsEmpty: true,
}

func IsEmptyValue(value interface{}) bool {
	var v, ok = value.(Empty)

	if ok {
		return v.IsEmpty
	}

	return false
}

func InitField(field BaseField, args *InitArgs) BaseField {
	field.Init(args)

	k, ok := interface{}(field).(DefaultMessages)

	if ok {
		d := k.GetDefaultMessages()

		for key, message := range *d {
			k.SetErrorMessage(key, message)
		}
	}

	return field
}

type Field struct {
	BaseField

	readOnly      bool
	writeOnly     bool
	required      bool
	fieldName     string
	allowNull     bool
	label         string
	errorMessages Messages

	initial  interface{}
	_default interface{}

	validators []Validator
	parent     serializers.BaseSerializer
}

func (f *Field) Parent() serializers.BaseSerializer {
	return f.parent
}

func (f *Field) SetParent(parent serializers.BaseSerializer) {
	f.parent = parent
}

func (f *Field) Validators() []Validator {
	return f.validators
}

func (f *Field) SetValidators(validators []Validator) {
	f.validators = validators
}

func (f *Field) Default() interface{} {
	return f._default
}

func (f *Field) SetDefault(_default interface{}) {
	f._default = _default
}

func (f *Field) Initial() interface{} {
	return f.initial
}

func (f *Field) SetInitial(initial interface{}) {
	f.initial = initial
}

func (f *Field) ErrorMessages() Messages {
	return f.errorMessages
}

func (f *Field) SetErrorMessages(errorMessages Messages) {
	f.errorMessages = errorMessages
}

func (f *Field) Label() string {
	return f.label
}

func (f *Field) SetLabel(label string) {
	f.label = label
}

func (f *Field) AllowNull() bool {
	return f.allowNull
}

func (f *Field) SetAllowNull(allowNull bool) {
	f.allowNull = allowNull
}

func (f *Field) FieldName() string {
	return f.fieldName
}

func (f *Field) SetFieldName(fieldName string) {
	f.fieldName = fieldName
}

func (f *Field) Required() bool {
	return f.required
}

func (f *Field) SetRequired(required bool) {
	f.required = required
}

func (f *Field) WriteOnly() bool {
	return f.writeOnly
}

func (f *Field) SetWriteOnly(writeOnly bool) {
	f.writeOnly = writeOnly
}

func (f *Field) ReadOnly() bool {
	return f.readOnly
}

func (f *Field) SetReadOnly(readOnly bool) {
	f.readOnly = readOnly
}

func (f *Field) Init(args *InitArgs) {
	defaultArgs := InitArgs{
		"readOnly":  false,
		"writeOnly": false,
		"required":  false,
		"fieldName": "",
		"allowNull": false,
		"label":     "",
		"initial":   nil,
		"default":   EmptyValue,
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

	f.readOnly = _args["readOnly"].(bool)
	f.writeOnly = _args["writeOnly"].(bool)
	f.required = _args["required"].(bool)
	f.fieldName = _args["fieldName"].(string)
	f.allowNull = _args["allowNull"].(bool)
	f.label = _args["label"].(string)

	if _args["initial"] != nil {
		f.initial = _args["initial"]
	}
	if _args["default"] != nil {
		f._default = _args["default"]
	}

	if f.readOnly && f.writeOnly {
		panic("Read only, write only same time error")
	}

	if f.errorMessages == nil {
		f.errorMessages = map[string]string{}
	}
}

func (f *Field) SetErrorMessage(key string, message string)  {
	f.errorMessages[key] = message
}

func (f *Field) Root() serializers.BaseSerializer {
	var r serializers.BaseSerializer = f.Parent()

	for r != nil {
		r = r.Parent()
	}

	return r
}

func (f *Field) Fail(key string) *errors.ValidationError {
	msg := f.errorMessages[key]

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

func (f *Field) ValidateEmptyValues(data interface{}) (bool, interface{}) {
	if f.readOnly {
		return true, f.GetDefault()
	}

	if IsEmptyValue(f._default) {
		if f.Root().IsPartial() == true {
			return false, errors.SkipField{}
		}

		if f.required {
			return false, f.Fail("required")
		}

		return true, f.GetDefault()
	}

	if data == nil {
		if f.allowNull == false {
			return false, f.Fail("null")
		}

		return true, nil
	}

	return false, data
}

func (f *Field) GetDefault() interface{} {
	if IsEmptyValue(f._default) || f.parent.IsPartial() == false {
		return errors.SkipField{}
	}

	if d, ok := f._default.(DefaultCaller); ok {
		return d()
	}

	return f._default
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

	for _, v := range f.validators {
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