package http_test

import (
	"encoding/json"
	"github.com/caudaganesh/go-pbt-demo/internal/article/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	articleHttp "github.com/caudaganesh/go-pbt-demo/internal/article/delivery/http"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	var mockArticle entity.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockArticleUC := mock_article.NewMockUsecase(ctrl)
	mockListArticle := make([]*entity.Article, 0)
	mockListArticle = append(mockListArticle, &mockArticle)
	cursor := "2"
	mockArticleUC.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockListArticle, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := articleHttp.ArticleHandler{
		AUsecase: mockArticleUC,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArticleUC := mock_article.NewMockUsecase(ctrl)
	mockArticleUC.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "", entity.ErrInternalServerError)
	cursor := "2"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := articleHttp.ArticleHandler{
		AUsecase: mockArticleUC,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetByID(t *testing.T) {
	var mockArticle entity.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	mockArticleUC := mock_article.NewMockUsecase(ctrl)
	mockArticleUC.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&mockArticle, nil)

	num := int(mockArticle.ID)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := articleHttp.ArticleHandler{
		AUsecase: mockArticleUC,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestStore(t *testing.T) {
	mockArticle := entity.Article{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	ctrl := gomock.NewController(t)
	mockUCase := mock_article.NewMockUsecase(ctrl)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockUCase.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/article", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/article")

	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestDelete(t *testing.T) {
	var mockArticle entity.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	ctrl := gomock.NewController(t)
	mockUCase := mock_article.NewMockUsecase(ctrl)

	num := int(mockArticle.ID)
	mockUCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)

}
