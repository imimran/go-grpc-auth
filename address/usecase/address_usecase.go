package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imimran/go-grpc-auth/address/domain"
	"github.com/imimran/go-grpc-auth/address/repository"
)

type AddressUsecase struct {
	repo repository.AddressRepository
}

func NewAddressUsecase(r repository.AddressRepository) *AddressUsecase {
	return &AddressUsecase{repo: r}
}

func (u *AddressUsecase) Create(ctx context.Context, addr *domain.Address) error {
	addr.ID = uuid.New()
	addr.Geom = fmt.Sprintf(
		"SRID=4326;POINT(%f %f)",
		addr.Longitude,
		addr.Latitude,
	)
	return u.repo.Create(ctx, addr)
}

func (u *AddressUsecase) GetByID(ctx context.Context, id string) (*domain.Address, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *AddressUsecase) Update(ctx context.Context, addr *domain.Address) error {
	addr.Geom = fmt.Sprintf(
		"SRID=4326;POINT(%f %f)",
		addr.Longitude,
		addr.Latitude,
	)
	return u.repo.Update(ctx, addr)
}

func (u *AddressUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
