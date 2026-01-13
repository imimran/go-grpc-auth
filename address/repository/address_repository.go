package repository

import (
	"context"

	"github.com/imimran/go-grpc-auth/address/domain"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(ctx context.Context, addr *domain.Address) error
	FindByID(ctx context.Context, id string) (*domain.Address, error)
	Update(ctx context.Context, addr *domain.Address) error
	Delete(ctx context.Context, id string) error
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) AddressRepository {
	return &addressRepo{db: db}
}

// Create address
func (r *addressRepo) Create(ctx context.Context, addr *domain.Address) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

// Find address by ID
func (r *addressRepo) FindByID(ctx context.Context, id string) (*domain.Address, error) {
	var addr domain.Address
	err := r.db.WithContext(ctx).
		First(&addr, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

// Update address
func (r *addressRepo) Update(ctx context.Context, addr *domain.Address) error {
	return r.db.WithContext(ctx).Save(addr).Error
}

// Delete address
func (r *addressRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Delete(&domain.Address{}, "id = ?", id).Error
}
