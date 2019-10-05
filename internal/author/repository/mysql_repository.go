package repository

import (
	"context"
	"database/sql"

	"github.com/caudaganesh/go-pbt-demo/internal/author"
	"github.com/caudaganesh/go-pbt-demo/internal/entity"
	"github.com/sirupsen/logrus"
)

type mysqlAuthorRepo struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMysqlAuthorRepository(db *sql.DB) author.Repository {
	return &mysqlAuthorRepo{
		DB: db,
	}
}

func (m *mysqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (*entity.Author, error) {

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	a := &entity.Author{}

	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return a, nil
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id int64) (*entity.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
