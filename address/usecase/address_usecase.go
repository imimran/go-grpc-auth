package usecase

import (
	"context"
	"fmt"
	"strings"

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

func (u *AddressUsecase) List(ctx context.Context, page, limit int) ([]domain.Address, int64, error) {
	// 1. Business Logic: Prevent extreme limits
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	// 2. Call Repository to get data and total count
	addresses, total, err := u.repo.List(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return addresses, total, nil
}

func (u *AddressUsecase) Create(ctx context.Context, addr *domain.Address) error {
	// 1. Generate a new valid UUID (Fixes the 0000... error)
	addr.ID = uuid.New()

	// 2. Implement requirement: normalized_address is raw_address
	// We lowercase and trim it to make the uniqueIndex effective
	addr.NormalizedAddress = strings.ToLower(strings.TrimSpace(addr.RawAddress))

	// 3. Set the PostGIS Geom field (Long then Lat)
	addr.Geom = fmt.Sprintf("SRID=4326;POINT(%f %f)",
		addr.Coordinates.Longitude,
		addr.Coordinates.Latitude,
	)

	return u.repo.Create(ctx, addr)
}
func (u *AddressUsecase) GetByID(ctx context.Context, id string) (*domain.Address, error) {
	return u.repo.FindByID(ctx, id)
}


func (u *AddressUsecase) Update(ctx context.Context, addr *domain.Address) error {
    // Re-normalize in case the RawAddress changed
    addr.NormalizedAddress = strings.ToLower(strings.TrimSpace(addr.RawAddress))

    // Re-calculate spatial point
    addr.Geom = fmt.Sprintf("SRID=4326;POINT(%f %f)", 
        addr.Coordinates.Longitude, 
        addr.Coordinates.Latitude,
    )

    return u.repo.Update(ctx, addr)
}

func (u *AddressUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
