package service

import (
	"context"
	"errors"

	"market/internal/domain"
	"market/internal/repository"
)

type PictureService struct {
	products repository.ProductRepository
	pictures repository.PictureRepository
	maxSize  int64
}

func NewPictureService(products repository.ProductRepository, pictures repository.PictureRepository) *PictureService {
	return &PictureService{
		products: products,
		pictures: pictures,
		maxSize:  10 << 20, // 10 MiB
	}
}

func (s *PictureService) UploadAndAttach(ctx context.Context, sellerID, productID int64, data []byte, mime string) (*domain.Picture, error) {
	if int64(len(data)) == 0 || int64(len(data)) > s.maxSize {
		return nil, errors.New("invalid file size")
	}
	p, err := s.products.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if p.SellerID != sellerID {
		return nil, errors.New("forbidden: not owner")
	}
	picID, err := s.pictures.Create(ctx, data, mime)
	if err != nil {
		return nil, err
	}
	pos, err := s.pictures.AttachAutoPosition(ctx, productID, picID)
	if err != nil {
		return nil, err
	}
	return &domain.Picture{
		ID:        picID,
		MIMEType:  mime,
		SizeBytes: int64(len(data)),
		Position:  pos,
	}, nil
}

func (s *PictureService) List(ctx context.Context, productID int64) ([]domain.Picture, error) {
	return s.pictures.ListByProduct(ctx, productID)
}

func (s *PictureService) Download(ctx context.Context, pictureID int64) ([]byte, string, error) {
	return s.pictures.GetData(ctx, pictureID)
}

func (s *PictureService) Detach(ctx context.Context, sellerID, productID, pictureID int64, hardDelete bool) error {
	p, err := s.products.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if p.SellerID != sellerID {
		return errors.New("forbidden: not owner")
	}
	if err := s.pictures.Detach(ctx, productID, pictureID); err != nil {
		return err
	}
	if hardDelete {
		return s.pictures.DeletePicture(ctx, pictureID)
	}
	return nil
}

func (s *PictureService) SetCover(ctx context.Context, sellerID, productID, pictureID int64) error {
	p, err := s.products.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if p.SellerID != sellerID {
		return errors.New("forbidden: not owner")
	}
	return s.pictures.SetCoverIfAttached(ctx, productID, pictureID)
}
