package author

import (
	"context"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
)

// Repository represent the author's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entity.Author, error)
}
