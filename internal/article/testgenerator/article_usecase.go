package testgenerator

import (
	authorgenerator "github.com/caudaganesh/go-pbt-demo/internal/author/testgenerator"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
)

type FetchArticleUseCaseGenerator struct {
	FetchArticleRepoResp  FetchArticleRepoResp
	GetByIDAuthorRepoResp authorgenerator.GetByIDAuthorRepoResp
	NumArgs               int
}

type FetchArticleRepoResp struct {
	Response   []*entity.Article
	NextCursor string
	Err        error
}
