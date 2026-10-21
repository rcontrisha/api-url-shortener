package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rcontrisha/api-url-shortener/internal/dto"
	"github.com/rcontrisha/api-url-shortener/internal/model"
	"github.com/rcontrisha/api-url-shortener/internal/repository"
	"github.com/rcontrisha/api-url-shortener/pkg"
)

type AuthService struct {
	repo *repository.AuthRepo
}

func NewAuthService(repo *repository.AuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) LoginService(ctx context.Context, payload dto.LoginRequest) (string, error) {
	result, err := a.repo.FindAccount(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("Incorrect Email or Password.	")
		}
		return "", err
	}

	if err := pkg.Compare(payload.Password, result.Password); err != nil {
		log.Println(err)
		return "", err
	}

	claims := pkg.NewJWTClaims(result.Id)
	token, err := claims.GenToken()

	return token, nil
}

func (a *AuthService) RegisterService(ctx context.Context, payload dto.LoginRequest) error {
	_, err := a.repo.FindAccount(ctx, payload.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if err == nil {
		return errors.New("User Already Exist.")
	}

	hashedPwd := pkg.NewRecommendedHashConfig().GenHash(payload.Password)

	var user = model.User{
		Email:    payload.Email,
		Password: hashedPwd,
	}

	return a.repo.CreateAccount(ctx, user)
}

func (a *AuthService) Logout(ctx context.Context, jti string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	return a.repo.BlacklistToken(ctx, jti, ttl)
}