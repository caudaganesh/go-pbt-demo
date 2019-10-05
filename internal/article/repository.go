package article

import (
	"context"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
)

// Repository represent the article's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*entity.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*entity.Article, error)
	GetByTitle(ctx context.Context, title string) (*entity.Article, error)
	Update(ctx context.Context, ar *entity.Article) error
	Store(ctx context.Context, a *entity.Article) error
	Delete(ctx context.Context, id int64) error
}
