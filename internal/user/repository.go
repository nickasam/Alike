package user

import (
	"github.com/Alike/internal/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *Repository) FindNearby(lat, lng float64, radiusKm float64, limit int) ([]domain.User, error) {
	var users []domain.User
	
	// Simple distance calculation (for production, use PostGIS)
	// This is a simplified bounding box query
	latDelta := radiusKm / 111.0 // approximate km per degree latitude
	lngDelta := radiusKm / (111.0 * 0.866) // approximate km per degree longitude
	
	err := r.db.Where("location_lat BETWEEN ? AND ?", lat-latDelta, lat+latDelta).
		Where("location_lng BETWEEN ? AND ?", lng-lngDelta, lng+lngDelta).
		Where("is_active = ?", true).
		Limit(limit).
		Find(&users).Error
		
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

func (r *Repository) List(limit, offset int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
