package services

import (
	"errors"
	"sync"
	"time"
	"vigilant-spork/models"
	"vigilant-spork/repository"

	"github.com/gofrs/uuid"
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

// Creates a new ReviewService
func NewReviewService(reviewRepo repository.ReviewRepository, productRepo repository.ProductRepository) *ReviewService {
	return &ReviewService{
		ReviewRepo:  reviewRepo,
		ProductRepo: productRepo,
		rateLimiter: make(map[uuid.UUID][]time.Time),
	}
}

// SubmitReview creates or updates a user's review for a product
func (s *ReviewService) SubmitReview(review *models.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		return ErrInvalidRating
	}

	// Rate limiting: max 5 reviews per user per minute
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

	// Check if user already reviewed this product
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

	// Update product rating aggregates
	avg, count, err := s.ReviewRepo.CalculateProductReviewAggregates(review.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateAggregates(review.ProductID, avg, count)
}

// UpdateReview allows a user to update their existing review
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

	// Update aggregates
	avg, count, err := s.ReviewRepo.CalculateProductReviewAggregates(review.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateAggregates(review.ProductID, avg, count)
}

// DeleteReview allows a user to delete their review
func (s *ReviewService) DeleteReview(reviewID uuid.UUID) error {
	// Get review to find productID
	review, err := s.ReviewRepo.GetReviewByID(reviewID)
	if err != nil || review == nil {
		return ErrReviewNotFound
	}

	if err := s.ReviewRepo.DeleteReview(reviewID); err != nil {
		return err
	}

	// Update aggregates
	avg, count, err := s.ReviewRepo.CalculateProductReviewAggregates(review.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateAggregates(review.ProductID, avg, count)
}

// GetReviewsForProduct fetches all reviews for a product
func (s *ReviewService) GetReviewsForProduct(productID uuid.UUID) ([]models.Review, error) {
	return s.ReviewRepo.GetReviewsByProductID(productID)
}

// GetReviewByUserForProduct fetches a review by a user for a specific product
func (s *ReviewService) GetReviewByUserForProduct(userID, productID uuid.UUID) (*models.Review, error) {
	return s.ReviewRepo.GetReviewByUserForProduct(userID, productID)
}