package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64 `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null;type:varchar(255)"`
	Password string `gorm:"not null"` // Stored as bcrypt hash
	FullName string `gorm:"type:varchar(255)"`
}

func NewUser(email, password, fullName string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    email,
		Password: hashed,
		FullName: fullName,
	}, nil
}

func (u *User) Update(email, password, fullName string) error {
	if email != "" {
		u.Email = email
	}
	if password != "" {
		hashed, err := hashPassword(password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	if fullName != "" {
		u.FullName = fullName
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}
