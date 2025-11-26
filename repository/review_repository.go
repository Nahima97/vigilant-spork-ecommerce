package repository

import (
	"vigilant-spork/db"
	"vigilant-spork/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateReview(review *models.Review) error
	GetReviewsByProductID(productID uuid.UUID) ([]models.Review, error)
	GetReviewByUserForProduct(userID, productID uuid.UUID) (*models.Review, error)
	GetReviewByID(reviewID uuid.UUID) (*models.Review, error) 
	UpdateReview(review *models.Review) error
	DeleteReview(id uuid.UUID) error
	CalculateProductReviewAggregates(productID uuid.UUID) (avg float64, count int64, err error)
}

type ReviewRepo struct {
	Db *gorm.DB
}

// Constructor
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &ReviewRepo{Db: db}
}

func (r *ReviewRepo) CreateReview(review *models.Review) error {
	return r.Db.Create(review).Error
}

func (r *ReviewRepo) GetReviewsByProductID(productID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	err := db.Db.Where("product_id = ?", productID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepo) GetReviewByUserForProduct(userID, productID uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := db.Db.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// NEW method to fetch by review ID
func (r *ReviewRepo) GetReviewByID(reviewID uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := db.Db.First(&review, "id = ?", reviewID).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepo) UpdateReview(review *models.Review) error {
	return db.Db.Save(review).Error
}

func (r *ReviewRepo) DeleteReview(id uuid.UUID) error {
	return db.Db.Delete(&models.Review{}, id).Error
}

func (r *ReviewRepo) CalculateProductReviewAggregates(productID uuid.UUID) (float64, int64, error) {
	var result struct {
		AvgRating   float64 `gorm:"column:avg_rating"`
		ReviewCount int64   `gorm:"column:review_count"`
	}
	err := db.Db.Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("AVG(rating) AS avg_rating, COUNT(*) AS review_count").
		Scan(&result).Error
	return result.AvgRating, result.ReviewCount, err
}
