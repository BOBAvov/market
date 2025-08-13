package service

import (
	"context"
	"errors"

	"market/internal/domain"
	"market/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

type ProductCreateInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"`
	Stock       int    `json:"stock"`
}

type ProductUpdateInput = ProductCreateInput

func (s *ProductService) Create(ctx context.Context, sellerID int64, in ProductCreateInput) (*domain.Product, error) {
	if in.Name == "" || in.PriceCents < 0 || in.Stock < 0 {
		return nil, errors.New("invalid product data")
	}
	p := &domain.Product{
		SellerID:    sellerID,
		Name:        in.Name,
		Description: in.Description,
		PriceCents:  in.PriceCents,
		Stock:       in.Stock,
	}
	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, sellerID, productID int64, in ProductUpdateInput) (*domain.Product, error) {
	if in.Name == "" || in.PriceCents < 0 || in.Stock < 0 {
		return nil, errors.New("invalid product data")
	}
	p, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if p.SellerID != sellerID {
		return nil, errors.New("forbidden: not owner")
	}
	p.Name = in.Name
	p.Description = in.Description
	p.PriceCents = in.PriceCents
	p.Stock = in.Stock
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, sellerID, productID int64) error {
	p, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if p.SellerID != sellerID {
		return errors.New("forbidden: not owner")
	}
	return s.repo.Delete(ctx, productID)
}

func (s *ProductService) Get(ctx context.Context, productID int64) (*domain.Product, error) {
	return s.repo.GetByID(ctx, productID)
}

func (s *ProductService) List(ctx context.Context, limit, offset int32, q string) ([]domain.Product, error) {
	return s.repo.List(ctx, repository.ProductFilter{
		Limit:  limit,
		Offset: offset,
		Query:  q,
	})
}
