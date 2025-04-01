package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"log"
)

var (
	ErrNotFound = errors.New("news with id not found")
)

func CreatePostgreStorage(dsn string) (Storage, error) {
	queryNews := `CREATE TABLE IF NOT EXISTS News (
  Id SERIAL PRIMARY KEY,
  Title TEXT NOT NULL,
  Content TEXT NOT NULL
);`

	queryCategories := `CREATE TABLE IF NOT EXISTS NewsCategories (
  NewsId BIGINT NOT NULL,
  CategoryId BIGINT NOT NULL,
  PRIMARY KEY (NewsId, CategoryId)
);`
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	reformDB := reform.NewDB(db, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))
	_, err = reformDB.Query(queryNews)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(queryCategories)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: reformDB}, nil
}

func (p *PostgresStorage) CreateNew(ctx context.Context, new New) error {
	_, err := p.db.QueryContext(ctx, "INSERT INTO News (Id, Title, Content) VALUES ($1, $2, $3)", new.Id, new.Title, new.Content)
	if err != nil {
		return err
	}

	for _, categoryId := range new.Categories {
		categoryQuery := `INSERT INTO NewsCategories (NewsId, CategoryId) VALUES ($1, $2)`
		_, err = p.db.Exec(categoryQuery, new.Id, categoryId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresStorage) GetNews(ctx context.Context) ([]New, error) {
	var result []New

	query := `SELECT id, title, content FROM News`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n New
		err = rows.Scan(&n.Id, &n.Title, &n.Content)
		if err != nil {
			return nil, err
		}

		categoriesQuery := `SELECT categoryid FROM NewsCategories WHERE newsid = $1`
		var categories []int
		catRows, err := p.db.QueryContext(ctx, categoriesQuery, n.Id)
		if err != nil {
			return nil, err
		}
		defer catRows.Close()

		for catRows.Next() {
			var categoryId int
			err = catRows.Scan(&categoryId)
			if err != nil {
				return nil, err
			}
			categories = append(categories, categoryId)
		}

		n.Categories = categories

		result = append(result, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostgresStorage) GetNewFromID(ctx context.Context, id int) (New, error) {
	var result New

	query := `SELECT Id, Title, Content FROM News WHERE Id = $1`

	err := p.db.QueryRowContext(ctx, query, id).Scan(&result.Id, &result.Title, &result.Content)

	if err != nil {
		if err == sql.ErrNoRows {
			return result, ErrNotFound
		}
		return result, err
	}

	categoriesQuery := `SELECT CategoryId FROM NewsCategories WHERE NewsId = $1`
	rows, err := p.db.QueryContext(ctx, categoriesQuery, id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId int
		if err := rows.Scan(&categoryId); err != nil {
			return result, err
		}
		result.Categories = append(result.Categories, categoryId)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (p *PostgresStorage) UpdateNewFromID(ctx context.Context, oldID int, new New) error {
	updateQuery := `
		UPDATE News
		SET Id = $1, Title = $2, Content = $3
		WHERE Id = $4
	`
	_, err := p.db.ExecContext(ctx, updateQuery, new.Id, new.Title, new.Content, oldID)
	if err != nil {
		return fmt.Errorf("failed to update news with ID %d: %w", oldID, err)
	}

	deleteCategoriesQuery := `
		DELETE FROM NewsCategories WHERE NewsId = $1
	`
	_, err = p.db.ExecContext(ctx, deleteCategoriesQuery, oldID)
	if err != nil {
		return err
	}

	for _, categoryID := range new.Categories {
		insertCategoryQuery := `
			INSERT INTO NewsCategories (NewsId, CategoryId)
			VALUES ($1, $2)
		`
		_, err = p.db.ExecContext(ctx, insertCategoryQuery, new.Id, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
