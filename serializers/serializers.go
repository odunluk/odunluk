package serializers

import (
	"github.com/odunluk/odunluk/v1/fields"
)

type Context map[string]interface{}

var (
	defaultArgs = fields.InitArgs{
		"instance": nil,
		"data":     fields.EmptyValue,
		"partial":  false,
		"context":  nil,
	}
)

type BaseSerializer interface {
	fields.BaseField

	IsPartial() bool
}

type Serializer struct {
	fields.Field
	BaseSerializer

	instance    interface{}
	initialData map[string]interface{}
	partial     bool
	context     Context
}

func (s *Serializer) InitialData() map[string]interface{} {
	return s.initialData
}

func (s *Serializer) SetInitialData(initialData map[string]interface{}) {
	s.initialData = initialData
}

func (s *Serializer) Instance() interface{} {
	return s.instance
}

func (s *Serializer) SetInstance(instance interface{}) {
	s.instance = instance
}

func (s *Serializer) Init(args *fields.InitArgs) {
	fields.InitField(s, args)

	var _args fields.InitArgs

	if args != nil {
		_args = *args
	} else {
		_args = fields.InitArgs{}
	}

	for k, v := range defaultArgs {
		_args[k] = v
	}

	if _args["instance"] != nil {
		s.instance = _args["instance"]
	}

	if _args["data"] != fields.EmptyValue {
		s.initialData = _args["data"].(map[string]interface{})
	}

	if _args["context"] != nil {
		s.context = _args["context"].(Context)
	}
}

func NewSerializer(serializer *Serializer, args *fields.InitArgs) BaseSerializer {
	serializer.Init(args)

	return serializer
}

func main() {
	NewSerializer(&Serializer{}, &fields.InitArgs{
		"Data": map[string]interface{}{
			"username": "kalaomer",
			"password": "password",
		},
	})
}
