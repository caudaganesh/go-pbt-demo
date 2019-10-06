package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_article "github.com/caudaganesh/go-pbt-demo/internal/article/mocks"
	"github.com/caudaganesh/go-pbt-demo/internal/article/usecase"
	mock_author "github.com/caudaganesh/go-pbt-demo/internal/author/mocks"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	// t.Skip()
	ctrl := gomock.NewController(t)
	mockArticleRepo := mock_article.NewMockRepository(ctrl)
	mockArticle := &entity.Article{
		Title:   "Hello",
		Content: "Content",
		Author: entity.Author{
			ID: 1,
		},
	}

	mockListArtilce := make([]*entity.Article, 0)
	mockListArtilce = append(mockListArtilce, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockAuthor := &entity.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		expect := []*entity.Article{
			&entity.Article{
				Title:   "Hello",
				Content: "Content",
				Author: entity.Author{
					ID:   1,
					Name: "Iman Tumorang",
				},
			},
		}

		mockArticleRepo.EXPECT().
			Fetch(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockListArtilce, "next-cursor", nil)
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		mockAuthorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(mockAuthor, nil)

		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		list, nextCursor, err := u.Fetch(context.TODO(), "12", int64(1))

		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))
		assert.Equal(t, "next-cursor", nextCursor)
		assert.Equal(t, expect, list)

	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.EXPECT().
			Fetch(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockListArtilce, "next-cursor", errors.New("Unexpexted Error"))
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)

		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
	})

}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArticleRepo := mock_article.NewMockRepository(ctrl)
	mockArticle := entity.Article{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := &entity.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&mockArticle, nil)
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		mockAuthorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(mockAuthor, nil)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		expect := &entity.Article{
			Title:   "Hello",
			Content: "Content",
			Author: entity.Author{
				ID:   1,
				Name: "Iman Tumorang",
			},
		}
		assert.NoError(t, err)
		assert.NotNil(t, a)
		assert.Equal(t, expect, a)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Nil(t, a)

	})

}

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArticleRepo := mock_article.NewMockRepository(ctrl)
	mockArticle := entity.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(nil, entity.ErrNotFound)
		mockArticleRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(&existingArticle, nil)
		mockAuthor := &entity.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}

		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		mockAuthorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(mockAuthor, nil)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
	})

}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArticleRepo := mock_article.NewMockRepository(ctrl)
	mockArticle := entity.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&mockArticle, nil)
		mockArticleRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)

		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
	})
	t.Run("article-is-not-exist", func(t *testing.T) {
		mockArticleRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil)

		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockArticleRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
	})

}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArticleRepo := mock_article.NewMockRepository(ctrl)
	mockArticle := entity.Article{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockAuthorRepo := mock_author.NewMockRepository(ctrl)
		u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Update(context.TODO(), &mockArticle)
		assert.NoError(t, err)
	})
}
