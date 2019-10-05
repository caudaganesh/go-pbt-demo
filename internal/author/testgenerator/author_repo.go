package testgenerator

import (
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
)

type GetByIDAuthorRepoResp struct {
	Response *entity.Author
	Err      error
}
