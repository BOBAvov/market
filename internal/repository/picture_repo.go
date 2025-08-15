package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"market/internal/domain"
)

var ErrPictureNotFound = errors.New("picture not found")
var ErrNotAttached = errors.New("picture is not attached to product")

type PictureRepository interface {
	Create(ctx context.Context, data []byte, mime string) (int64, error)
	AttachAutoPosition(ctx context.Context, productID, pictureID int64) (int, error)
	ListByProduct(ctx context.Context, productID int64) ([]domain.Picture, error)
	GetData(ctx context.Context, pictureID int64) ([]byte, string, error)
	Detach(ctx context.Context, productID, pictureID int64) error
	DeletePicture(ctx context.Context, pictureID int64) error
	SetCoverIfAttached(ctx context.Context, productID, pictureID int64) error
}

type pictureRepo struct {
	pool *pgxpool.Pool
}

func NewPictureRepository(pool *pgxpool.Pool) PictureRepository {
	return &pictureRepo{pool: pool}
}

func (r *pictureRepo) Create(ctx context.Context, data []byte, mime string) (int64, error) {
	var id int64
	_, err := r.pool.Exec(ctx, "SET LOCAL bytea_output='hex'") // no-op safety
	_ = err
	err = r.pool.QueryRow(ctx, `
		INSERT INTO pictures (data, mime_type, size_bytes)
		VALUES ($1, $2, $3) RETURNING id
	`, data, mime, int64(len(data))).Scan(&id)
	return id, err
}

func (r *pictureRepo) AttachAutoPosition(ctx context.Context, productID, pictureID int64) (int, error) {
	var pos int
	err := r.pool.QueryRow(ctx, `
		WITH next_pos AS (
		  SELECT COALESCE(MAX(position), 0) + 1 AS pos
		  FROM product_pictures
		  WHERE product_id = $1
		)
		INSERT INTO product_pictures (product_id, picture_id, position)
		SELECT $1, $2, pos FROM next_pos
		RETURNING position
	`, productID, pictureID).Scan(&pos)
	return pos, err
}

func (r *pictureRepo) ListByProduct(ctx context.Context, productID int64) ([]domain.Picture, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.mime_type, p.size_bytes, p.created_at, pp.position
		FROM product_pictures pp
		JOIN pictures p ON p.id = pp.picture_id
		WHERE pp.product_id = $1
		ORDER BY pp.position
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Picture
	for rows.Next() {
		var pic domain.Picture
		if err := rows.Scan(&pic.ID, &pic.MIMEType, &pic.SizeBytes, &pic.CreatedAt, &pic.Position); err != nil {
			return nil, err
		}
		out = append(out, pic)
	}
	return out, rows.Err()
}

func (r *pictureRepo) GetData(ctx context.Context, pictureID int64) ([]byte, string, error) {
	var data []byte
	var mime string
	err := r.pool.QueryRow(ctx, `
		SELECT data, mime_type FROM pictures WHERE id = $1
	`, pictureID).Scan(&data, &mime)
	if err != nil {
		return nil, "", err
	}
	return data, mime, nil
}

func (r *pictureRepo) Detach(ctx context.Context, productID, pictureID int64) error {
	ct, err := r.pool.Exec(ctx, `
		DELETE FROM product_pictures WHERE product_id = $1 AND picture_id = $2
	`, productID, pictureID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotAttached
	}
	// если удалили обложку — обнулим её
	_, _ = r.pool.Exec(ctx, `
		UPDATE products SET cover_picture_id = NULL
		WHERE id = $1 AND cover_picture_id = $2
	`, productID, pictureID)
	return nil
}

func (r *pictureRepo) DeletePicture(ctx context.Context, pictureID int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM pictures WHERE id = $1`, pictureID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrPictureNotFound
	}
	return nil
}

func (r *pictureRepo) SetCoverIfAttached(ctx context.Context, productID, pictureID int64) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE products p
		SET cover_picture_id = $2
		WHERE p.id = $1
		  AND EXISTS (SELECT 1 FROM product_pictures pp WHERE pp.product_id = $1 AND pp.picture_id = $2)
	`, productID, pictureID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotAttached
	}
	return nil
}
