package article

import (
	"context"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
)

// Usecase represent the article's usecases
type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*entity.Article, string, error)
	GetByID(ctx context.Context, id int64) (*entity.Article, error)
	Update(ctx context.Context, ar *entity.Article) error
	GetByTitle(ctx context.Context, title string) (*entity.Article, error)
	Store(context.Context, *entity.Article) error
	Delete(ctx context.Context, id int64) error
}
