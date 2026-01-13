package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/imimran/go-grpc-auth/address/domain"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(ctx context.Context, addr *domain.Address) error
	FindByID(ctx context.Context, id string) (*domain.Address, error)
	Update(ctx context.Context, addr *domain.Address) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]domain.Address, int64, error)
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) AddressRepository {
	return &addressRepo{db: db}
}

// Create persists a new address to the database
func (r *addressRepo) Create(ctx context.Context, addr *domain.Address) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *addressRepo) List(ctx context.Context, page, limit int) ([]domain.Address, int64, error) {
	var addresses []domain.Address
	var total int64

	// 1. Get total count using a fresh Model session
	// This ensures Count() doesn't interfere with the subsequent Find()
	if err := r.db.WithContext(ctx).Model(&domain.Address{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// If total is 0, don't even bother running the Find query
	if total == 0 {
		return []domain.Address{}, 0, nil
	}

	// 2. Fetch the records using a fresh session
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("id DESC"). // Use ID or created_at
		Find(&addresses).Error

	return addresses, total, err
}

// FindByID retrieves a single address by its UUID string
func (r *addressRepo) FindByID(ctx context.Context, id string) (*domain.Address, error) {
	if id == "" {
		return nil, errors.New("address id is required")
	}

	// Validate if the incoming string is a real UUID before querying DB
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	var addr domain.Address
	// Using .First returns gorm.ErrRecordNotFound if the ID doesn't exist
	err = r.db.WithContext(ctx).First(&addr, "id = ?", parsedID).Error
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

// Update modifies an existing address record using the struct's ID
func (r *addressRepo) Update(ctx context.Context, addr *domain.Address) error {
	// .Save() performs an UPSERT. If ID exists, it updates; if not, it inserts.
	result := r.db.WithContext(ctx).Save(addr)
	if result.Error != nil {
		return result.Error
	}

	// If no rows were affected, the ID provided didn't exist in the database
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete removes an address record by ID
func (r *addressRepo) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid format: %w", err)
	}

	// Delete requires a model type and the condition
	result := r.db.WithContext(ctx).Delete(&domain.Address{}, "id = ?", parsedID)
	if result.Error != nil {
		return result.Error
	}

	// Return Not Found error if the ID was not present in the DB
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
