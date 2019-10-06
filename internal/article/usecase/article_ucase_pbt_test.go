package usecase_test

import (
	"context"
	"testing"
	"testing/quick"
	"time"

	mock_article "github.com/caudaganesh/go-pbt-demo/internal/article/mocks"
	"github.com/caudaganesh/go-pbt-demo/internal/article/testgenerator"
	"github.com/caudaganesh/go-pbt-demo/internal/article/usecase"
	"github.com/caudaganesh/go-pbt-demo/internal/author/mocks"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/caudaganesh/go-pbt-demo/internal/generator"
	"github.com/golang/mock/gomock"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/arbitrary"
)

func TestFetchPBTTestingQuick(t *testing.T) {
	// t.Skip()
	ctrl := gomock.NewController(t)
	checkFetch := func(identifier testgenerator.FetchArticleUseCaseGenerator) bool {
		mockArticleRepo := mock_article.NewMockRepository(ctrl)
		mockArticleRepo.EXPECT().
			Fetch(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(identifier.FetchArticleRepoResp.Response,
				identifier.FetchArticleRepoResp.NextCursor,
				identifier.FetchArticleRepoResp.Err)
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		mockAuthorRepo.EXPECT().
			GetByID(gomock.Any(), gomock.Any()).
			Return(identifier.GetByIDAuthorRepoResp.Response,
				identifier.GetByIDAuthorRepoResp.Err).AnyTimes()

		uc := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second+2)
		articles, cursor, err := uc.Fetch(context.TODO(), "12", int64(identifier.NumArgs))

		return checkTestFetch(articles, cursor, err, &identifier)
	}

	if err := quick.CheckEqual(checkFetch, checkFetch, generator.GenerateQuickCheckConfig(15)); err != nil {
		t.Error("Error TestFetchPBT CheckEqual", err)
	}

	if err := quick.Check(checkFetch, generator.GenerateQuickCheckConfig(150)); err != nil {
		t.Error("Error TestFetchPBT", err)
	}
}

func TestFetchPBTGopter(t *testing.T) {
	// t.Skip()
	ctrl := gomock.NewController(t)
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234)
	parameters.MaxSize = 100
	arbitraries := arbitrary.DefaultArbitraries()
	arbitraries.RegisterGen(testgenerator.GenerateFetchArticleUseCase())
	properties := gopter.NewProperties(parameters)

	properties.Property("Should execute fetch articles", arbitraries.ForAll(
		func(identifier *testgenerator.FetchArticleUseCaseGenerator) bool {
			errFetchArtilce := generator.GenerateRandomError()
			identifier.FetchArticleRepoResp.Err = errFetchArtilce
			errGetByIDAuthorRepoResp := generator.GenerateRandomError()
			identifier.GetByIDAuthorRepoResp.Err = errGetByIDAuthorRepoResp

			mockArticleRepo := mock_article.NewMockRepository(ctrl)
			mockArticleRepo.EXPECT().
				Fetch(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(identifier.FetchArticleRepoResp.Response,
					identifier.FetchArticleRepoResp.NextCursor,
					identifier.FetchArticleRepoResp.Err)
			mockAuthorRepo := mock_author.NewMockRepository(ctrl)
			mockAuthorRepo.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(identifier.GetByIDAuthorRepoResp.Response,
					identifier.GetByIDAuthorRepoResp.Err).AnyTimes()

			uc := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second+2)
			articles, cursor, err := uc.Fetch(context.TODO(), "12", int64(identifier.NumArgs))

			return checkTestFetch(articles, cursor, err, identifier)
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(true))
}

func checkTestFetch(articles []*entity.Article, cursor string, err error, identifier *testgenerator.FetchArticleUseCaseGenerator) bool {
	notErrorWhenFetchingArticleRepoError := identifier.FetchArticleRepoResp.Err != nil && err == nil
	notErrorWhenArticlesExistsGetAuthorError := identifier.GetByIDAuthorRepoResp.Err != nil && len(identifier.FetchArticleRepoResp.Response) > 0 && err == nil
	fetchedDataLengthFromRepoIsDifferentWithTheResult := err == nil && len(identifier.FetchArticleRepoResp.Response) != len(articles)

	if notErrorWhenFetchingArticleRepoError ||
		notErrorWhenArticlesExistsGetAuthorError ||
		fetchedDataLengthFromRepoIsDifferentWithTheResult {
		return false
	}
	return true
}
