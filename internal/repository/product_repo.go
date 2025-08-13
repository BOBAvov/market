package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"market/internal/domain"
)

var ErrProductNotFound = errors.New("product not found")

type ProductFilter struct {
	Limit  int32
	Offset int32
	Query  string // optional name search
}

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f ProductFilter) ([]domain.Product, error)
}

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) ProductRepository {
	return &productRepo{pool: pool}
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO products (seller_id, name, description, price_cents, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, p.SellerID, p.Name, p.Description, p.PriceCents, p.Stock).
		Scan(&id, &p.CreatedAt, &p.UpdatedAt)
	return id, err
}

func (r *productRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, seller_id, name, description, price_cents, stock, created_at, updated_at
		FROM products WHERE id = $1
	`, id)
	var p domain.Product
	if err := row.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.PriceCents, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) Update(ctx context.Context, p *domain.Product) error {
	return r.pool.QueryRow(ctx, `
		UPDATE products
		SET name = $2, description = $3, price_cents = $4, stock = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, p.ID, p.Name, p.Description, p.PriceCents, p.Stock).
		Scan(&p.UpdatedAt)
}

func (r *productRepo) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *productRepo) List(ctx context.Context, f ProductFilter) ([]domain.Product, error) {
	q := `
		SELECT id, seller_id, name, description, price_cents, stock, created_at, updated_at
		FROM products
	`
	args := []any{}
	if f.Query != "" {
		q += " WHERE lower(name) LIKE lower($1)"
		args = append(args, "%"+f.Query+"%")
	}
	q += " ORDER BY id DESC"
	if f.Limit <= 0 {
		f.Limit = 50
	}
	q += " LIMIT $2 OFFSET $3"
	args = append(args, f.Limit, f.Offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.PriceCents, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}
