package testgenerator

import (
	"math/rand"
	"reflect"

	authorgenerator "github.com/caudaganesh/go-pbt-demo/internal/author/testgenerator"
	"github.com/caudaganesh/go-pbt-demo/internal/generator"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

/*
testing/quick

Generate is a function that will be executed for each iteration,
it will generate values based on anything we implement here
*/

func (FetchArticleUseCaseGenerator) Generate(rand *rand.Rand, size int) reflect.Value {
	resp := FetchArticleUseCaseGenerator{
		FetchArticleRepoResp: FetchArticleRepoResp{
			Err:        generator.GenerateRandomError(),
			NextCursor: generator.GenerateRandomString(),
			Response:   GenerateRandomArticles(),
		},
		GetByIDAuthorRepoResp: authorgenerator.GetByIDAuthorRepoResp{
			Err:      generator.GenerateRandomError(),
			Response: authorgenerator.GenerateRandomAuthor(),
		},
		NumArgs: rand.Intn(100),
	}
	return reflect.ValueOf(FetchArticleUseCaseGenerator(resp))
}

/*
leanovate/gopter
*/
func GenerateFetchArticleUseCase() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&FetchArticleUseCaseGenerator{}), map[string]gopter.Gen{
		"FetchArticleRepoResp ": GenerateFetchArticleRepoResp(),
		"GetByIDAuthorRepoResp": authorgenerator.GenerateGetByIDAuthorRepoResp(),
		"NumArgs":               gen.IntRange(0, 0),
	})
}

func GenerateFetchArticleRepoResp() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&FetchArticleRepoResp{}), map[string]gopter.Gen{
		"Response":   gen.SliceOf(gen.PtrOf(GenerateArticle())),
		"NextCursor": gen.AnyString(),
	})
}
