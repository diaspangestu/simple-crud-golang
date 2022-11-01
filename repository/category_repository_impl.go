package repository

import (
	"belajar-golang-restful-api/helper"
	"belajar-golang-restful-api/models/entity"
	"context"
	"database/sql"
	"errors"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (c *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category {
	query := "INSERT INTO category(name) VALUES (?)"
	result, err := tx.ExecContext(ctx, query, category.Name)
	helper.PanicError(err)

	id, err := result.LastInsertId()
	helper.PanicError(err)

	category.Id = int(id)
	return category
}

func (c *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category {
	query := "UPDATE category SET name = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, category.Name, category.Id)
	helper.PanicError(err)

	return category
}

func (c *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category entity.Category) {
	query := "DELETE FROM category WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, category.Id)
	helper.PanicError(err)
}

func (c *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error) {
	query := "SELECT id, name FROM category WHERE id = ?"
	rows, err := tx.QueryContext(ctx, query, categoryId)
	helper.PanicError(err)
	defer rows.Close()

	category := entity.Category{}
	if rows.Next() {
		rows.Scan(&category.Id, &category.Name)
		helper.PanicError(err)
		return category, nil
	} else {
		return category, errors.New("Category is not found")
	}
}

func (c CategoryRepositoryImpl) GetAll(ctx context.Context, tx *sql.Tx) []entity.Category {
	query := "SELECT id, name FROM category"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicError(err)
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		category := entity.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		helper.PanicError(err)

		categories = append(categories, category)
	}
	return categories
}
