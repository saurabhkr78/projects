package book

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, book *Book) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO books
		(title, author, isbn, published_at)
		VALUES ($1, $2, $3, $4)
		`,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublishedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*Book, error) {
	var book Book

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			title,
			author,
			isbn,
			published_at,
			created_at,
			updated_at
		FROM books
		WHERE id = $1
		`,
		id,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.PublishedAt,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Book, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			title,
			author,
			isbn,
			published_at,
			created_at,
			updated_at
		FROM books
		ORDER BY id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublishedAt,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *PostgresRepository) Update(ctx context.Context, book *Book) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE books
		SET
			title = $1,
			author = $2,
			isbn = $3,
			published_at = $4,
			updated_at = NOW()
		WHERE id = $5
		`,
		book.Title,
		book.Author,
		book.ISBN,
		book.PublishedAt,
		book.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM books
		WHERE id = $1
		`,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
