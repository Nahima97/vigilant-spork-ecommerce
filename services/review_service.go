package services

import (
	"errors"
	"github.com/gofrs/uuid"
	"sync"
	"time"
	"vigilant-spork/models"
	"vigilant-spork/repository"
)

type ReviewService struct {
	ReviewRepo  repository.ReviewRepository
	ProductRepo repository.ProductRepository
	rateLimiter map[uuid.UUID][]time.Time
	mu          sync.Mutex
}

var (
	ErrInvalidRating     = errors.New("rating must be between 1 and 5")
	ErrRateLimitExceeded = errors.New("rate limit exceeded, try later")
	ErrReviewNotFound    = errors.New("review not found")
)

func NewReviewService(reviewRepo repository.ReviewRepository, productRepo repository.ProductRepository) *ReviewService {
	return &ReviewService{
		ReviewRepo:  reviewRepo,
		ProductRepo: productRepo,
		rateLimiter: make(map[uuid.UUID][]time.Time),
	}
}

func (s *ReviewService) SubmitReview(review *models.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		return ErrInvalidRating
	}

	s.mu.Lock() 
	now := time.Now()
	timestamps := s.rateLimiter[review.UserID]
	var recent []time.Time
	for _, t := range timestamps {
		if now.Sub(t) < time.Minute {
			recent = append(recent, t)
		}
	}
	
	if len(recent) >= 5 {
		s.mu.Unlock()
		return ErrRateLimitExceeded
	}
	recent = append(recent, now)
	s.rateLimiter[review.UserID] = recent
	s.mu.Unlock()

	
	existing, _ := s.ReviewRepo.GetReviewByUserForProduct(review.UserID, review.ProductID)
	if existing != nil {
		existing.Title = review.Title
		existing.Description = review.Description
		existing.Rating = review.Rating
		if err := s.ReviewRepo.UpdateReview(existing); err != nil {
			return err
		}
	} else {
		if err := s.ReviewRepo.CreateReview(review); err != nil {
			return err
		}
	}

	avg, count, err := s.ReviewRepo.CalculateProductReviewAggregates(review.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateAggregates(review.ProductID, avg, count)
}

func (s *ReviewService) UpdateReview(review *models.Review) error {
	existing, err := s.ReviewRepo.GetReviewByUserForProduct(review.UserID, review.ProductID)
	if err != nil || existing == nil {
		return ErrReviewNotFound
	}

	existing.Title = review.Title
	existing.Description = review.Description
	existing.Rating = review.Rating

	if err := s.ReviewRepo.UpdateReview(existing); err != nil {
		return err
	}

	avg, count, err := s.ReviewRepo.CalculateProductReviewAggregates(review.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateAggregates(review.ProductID, avg, count)
}

func (s *ReviewService) DeleteReview(reviewID uuid.UUID) error {
	review, err := s.ReviewRepo.GetReviewByID(reviewID)
	if err != nil || review == nil {
		return ErrReviewNotFound
	}

	if err := s.ReviewRepo.DeleteReview(reviewID); err != nil {
		return err
	}

	avg, count, err := s.ReviewRepo.CalculateProductReviewAggregates(review.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateAggregates(review.ProductID, avg, count)
}

func (s *ReviewService) GetReviewsForProduct(productID uuid.UUID) ([]models.Review, error) {
	return s.ReviewRepo.GetReviewsByProductID(productID)
}

func (s *ReviewService) GetReviewByUserForProduct(userID, productID uuid.UUID) (*models.Review, error) {
	return s.ReviewRepo.GetReviewByUserForProduct(userID, productID)
}
