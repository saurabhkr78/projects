package book

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func (r *PostgresRepository) Createa(ctx context.Context, book *Book) error {

}
func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*Book, error) {

}
func (r *PostgresRepository) GetAll(ctx context.Context) ([]Book, error) {

}
func (r *PostgresRepository) Update(ctx context.Context, book *Book) error {

}
func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {

}
