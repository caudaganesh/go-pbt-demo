package testgenerator

import (
	"reflect"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

func GenerateGetByIDAuthorRepoResp() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&GetByIDAuthorRepoResp{}), map[string]gopter.Gen{
		"Response": gen.PtrOf(GenerateAuthor()),
	})
}
