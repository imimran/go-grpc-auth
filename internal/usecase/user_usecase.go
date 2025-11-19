package usecase

import (
	"errors"
	"time"

	"github.com/imimran/go-grpc-auth/internal/domain"
	"github.com/imimran/go-grpc-auth/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type UserUsecase interface {
	Register(email, password, fullName string) (*domain.User, error)
	Login(email, password string) (string, error)
	Get(id int64) (*domain.User, error)
	Update(id int64, email, password, fullName string) (*domain.User, error)
	Delete(id int64) error
	List() ([]*domain.User, error)
}

type userUsecase struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(repo repository.UserRepository, jwtSecret []byte) UserUsecase {
	return &userUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (u *userUsecase) Register(email, password, fullName string) (*domain.User, error) {
	_, err := u.repo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	user, err := domain.NewUser(email, password, fullName)
	if err != nil {
		return nil, err
	}

	err = u.repo.Create(user)
	return user, err
}

func (u *userUsecase) Login(email, password string) (string, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if !user.CheckPassword(password) {
		return "", errors.New("invalid email or password")
	}

	token, err := u.generateJWT(user)
	return token, err
}

func (u *userUsecase) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}

func (u *userUsecase) Get(id int64) (*domain.User, error) {
	return u.repo.GetByID(id)
}

func (u *userUsecase) Update(id int64, email, password, fullName string) (*domain.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := user.Update(email, password, fullName); err != nil {
		return nil, err
	}
	if err := u.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}

func (u *userUsecase) List() ([]*domain.User, error) {
	return u.repo.List()
}
