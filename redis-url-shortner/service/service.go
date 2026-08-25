package service

//in this layer both repo and redis will be injected so this is the layer where we do cache aside business logic

import (
	"context"
	"log"
	"redis-url-shortner/client"
	urlModel "redis-url-shortner/model"
	"redis-url-shortner/repository"
	"redis-url-shortner/utils"

	"github.com/redis/go-redis/v9"
)

type URLService struct {
	Repo  repository.URLRepository // no pointer to the interface, because we want to use the interface methods directly
	Redis *client.Client
}
type Service interface {
	// SaveURL saves the original URL and returns the generated short ID
	SaveURL(ctx context.Context, originalURL string) (string, error)

	// GetOriginalURL retrieves the original URL based on the provided short ID
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

func NewURLService(repo repository.URLRepository, redc *client.Client) *URLService {
	return &URLService{
		Repo:  repo,
		Redis: redc,
	}
}

func (s *URLService) SaveURL(
	ctx context.Context,
	originalURL string,
) (string, error) {

	for {
		shortID := utils.GenerateShortID()

		// Check whether this ID already exists.
		_, err := s.Repo.GetOriginalURL(ctx, shortID)

		if err == nil {
			// Collision → generate another ID.
			continue
		}

		// For now, assume any error means "not found".
		// We'll improve this with a proper ErrNotFound later.
		url := urlModel.URL{
			ShortID:     shortID,
			OriginalURL: originalURL,
		}

		if err := s.Repo.SaveURL(ctx, url); err != nil {
			return "", err
		}

		return "http://localhost:8080/" + shortID, nil
	}
}

func (s *URLService) GetOriginalURL(
	ctx context.Context,
	shortID string,
) (string, error) {

	// Redis GET
	log.Println("Checking Redis...")
	originalURL, err := s.Redis.RedisClient.Get(ctx, shortID).Result()

	// Cache HIT
	if err == nil {
		log.Println("CACHE HIT")
		return originalURL, nil
	}

	// Redis error other than cache miss
	if err != redis.Nil {
		return "", err
	}

	// Cache MISS → Repository
	log.Println("CACHE MISS → checking DB...")
	originalURL, err = s.Repo.GetOriginalURL(ctx, shortID)
	if err != nil {
		return "", err
	}

	// Populate Redis
	err = s.Redis.RedisClient.Set(
		ctx,
		shortID,
		originalURL,
		0,
	).Err()

	if err != nil {
		return "", err
	}

	return originalURL, nil
}
