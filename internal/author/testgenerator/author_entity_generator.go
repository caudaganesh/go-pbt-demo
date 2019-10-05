package testgenerator

import (
	"reflect"

	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/caudaganesh/go-pbt-demo/internal/generator"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

func GenerateRandomAuthor() *entity.Author {
	return &entity.Author{
		CreatedAt: generator.GenerateRandomTimestamp().Format("2018-10-28"),
		ID:        int64(generator.GenerateRandomInt()),
		Name:      generator.GenerateRandomString(),
		UpdatedAt: generator.GenerateRandomTimestamp().Format("2018-10-28"),
	}
}

func GenerateAuthor() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&entity.Author{}), map[string]gopter.Gen{
		"ID":        gen.Int64(),
		"Name":      gen.AnyString(),
		"CreatedAt": gen.AnyString(),
		"UpdatedAt": gen.AnyString(),
	})
}
