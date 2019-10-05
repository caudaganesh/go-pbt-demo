package usecase

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/caudaganesh/go-pbt-demo/internal/article"
	"github.com/caudaganesh/go-pbt-demo/internal/author"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/sirupsen/logrus"
)

type articleUsecase struct {
	articleRepo    article.Repository
	authorRepo     author.Repository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewArticleUsecase(a article.Repository, ar author.Repository, timeout time.Duration) article.Usecase {
	return &articleUsecase{
		articleRepo:    a,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *articleUsecase) fillAuthorDetails(c context.Context, data []*entity.Article) ([]*entity.Article, error) {

	g, ctx := errgroup.WithContext(c)

	// Get the author's id
	mapAuthors := map[int64]entity.Author{}

	for _, article := range data {
		mapAuthors[article.Author.ID] = entity.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan *entity.Author)
	for authorID := range mapAuthors {
		// authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	var err error
	go func() {
		err = g.Wait()
		defer close(chanAuthor)
		if err != nil {
			logrus.Error(err)
			return
		}
	}()

	for author := range chanAuthor {
		if author != nil {
			mapAuthors[author.ID] = *author
		}
	}

	if err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) ([]*entity.Article, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listArticle, nextCursor, err := a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(listArticle) > 0 {
		listArticle, err = a.fillAuthorDetails(ctx, listArticle)
		if err != nil {
			return nil, "", err
		}
	}

	return listArticle, nextCursor, nil
}

func (a *articleUsecase) GetByID(c context.Context, id int64) (*entity.Article, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor
	return res, nil
}

func (a *articleUsecase) Update(c context.Context, ar *entity.Article) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleUsecase) GetByTitle(c context.Context, title string) (*entity.Article, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor

	return res, nil
}

func (a *articleUsecase) Store(c context.Context, m *entity.Article) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, m.Title)
	if existedArticle != nil {
		return entity.ErrConflict
	}

	err := a.articleRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *articleUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existedArticle == nil {
		return entity.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
