package service

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/rcontrisha/api-url-shortener/internal/dto"
	"github.com/rcontrisha/api-url-shortener/internal/repository"
	"github.com/rcontrisha/api-url-shortener/internal/utils"
	"go.step.sm/crypto/randutil"
)

type LinkService struct {
	repo *repository.LinkRepo
}

func NewLinkService(repo *repository.LinkRepo) *LinkService {
	return &LinkService{
		repo: repo,
	}
}

func (l *LinkService) ShortenLinkService(ctx context.Context, userId int, payload dto.ShortenLink) error {
	var finalSlug string
	var err error

	if payload.Slug != "" {
		if err := utils.ValidateCustomSlug(payload.Slug); err != nil {
			return err
		}
		finalSlug = payload.Slug
	} else {
		finalSlug, err = randutil.Alphanumeric(6)
		if err != nil {
			return err
		}
	}

	_, err = l.repo.GetOriginalUrl(ctx, finalSlug)
	if err == nil {
		return errors.New("slug already in use, please choose another one")
	}

	err = l.repo.SetShortenLink(ctx, userId, payload.OriginalUrl, finalSlug)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkService) GetOriginalUrl(ctx context.Context, slug string) (string, error) {
	link, err := l.repo.GetOriginalUrl(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("URL Not Found.")
		}
		return "", err
	}

	return link, nil
}

func (s *LinkService) DeleteShortenedLink(ctx context.Context, id int, userID int) error {
	return s.repo.DeleteShortenedLink(ctx, id, userID)
}

func (s *LinkService) GetUserLinks(ctx context.Context, userID int, search string, page int, limit int) (dto.GetUserLinksResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalActive, err := s.repo.CountUserLinks(ctx, userID)
	if err != nil {
		return dto.GetUserLinksResponse{}, err
	}

	links, err := s.repo.GetUserLinks(ctx, userID, search, page, limit)
	if err != nil {
		return dto.GetUserLinksResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalActive) / float64(limit)))

	return dto.GetUserLinksResponse{
		Data: links,
		Meta: dto.PaginationMeta{
			TotalActive: totalActive,
			Page:        page,
			Limit:       limit,
			TotalPages:  totalPages,
		},
	}, nil
}
