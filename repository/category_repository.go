package repository

import (
	"belajar-golang-restful-api/models/entity"
	"context"
	"database/sql"
)

type CategoryRepository interface {
	Create(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category
	Update(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category
	Delete(ctx context.Context, tx *sql.Tx, category entity.Category)
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error)
	GetAll(ctx context.Context, tx *sql.Tx) []entity.Category
}
