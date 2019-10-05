package testgenerator

import (
	"math/rand"
	"reflect"

	"github.com/caudaganesh/go-pbt-demo/internal/author/testgenerator"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/caudaganesh/go-pbt-demo/internal/generator"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

func GenerateRandomArticle() *entity.Article {
	return &entity.Article{
		Author:    *testgenerator.GenerateRandomAuthor(),
		Content:   generator.GenerateRandomString(),
		CreatedAt: generator.GenerateRandomTimestamp(),
		ID:        int64(generator.GenerateRandomInt()),
		Title:     generator.GenerateRandomString(),
		UpdatedAt: generator.GenerateRandomTimestamp(),
	}
}

func GenerateRandomArticles() []*entity.Article {
	count := rand.Intn(100)
	articles := []*entity.Article{}
	for index := 0; index < count; index++ {
		articles = append(articles, GenerateRandomArticle())
	}

	return articles
}

// gopter
func GenerateArticle() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&entity.Article{}), map[string]gopter.Gen{
		"ID":        gen.Int64(),
		"Title":     gen.AnyString(),
		"Content":   gen.AnyString(),
		"Author":    testgenerator.GenerateAuthor(),
		"UpdatedAt": gen.AnyTime(),
		"CreatedAt": gen.AnyTime(),
	})
}
