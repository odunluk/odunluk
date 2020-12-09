package serializers

import (
	"github.com/odunluk/odunluk/v1/fields"
)

type initArgs struct {
	data interface{}
}


type BaseSerializer struct {
	fields.Field
}

func (s *BaseSerializer) init(arg initArgs) {

}

func newSerializer(serializer BaseSerializer, args initArgs) BaseSerializer {
	serializer.init(args)

	return serializer
}

func main() {
	newSerializer(BaseSerializer{}, initArgs{
		data: map[string]interface{} {
			"username": "kalaomer",
			"password": "password",
		},
	})
}