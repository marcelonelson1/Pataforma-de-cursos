package models

import (
	"time"
	"mime/multipart"
	"gorm.io/gorm"
)

// Activity representa una actividad o clase en el sistema
type Activity struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Day         string         `gorm:"size:50;not null" json:"day"`
	Time        string         `gorm:"size:50;not null" json:"time"`
	Duration    int            `gorm:"not null" json:"duration"`
	Instructor  string         `gorm:"size:100;not null" json:"instructor"`
	Category    string         `gorm:"size:100;not null" json:"category"`
	Capacity    int            `gorm:"not null" json:"capacity"`
	Enrolled    int            `gorm:"default:0" json:"enrolled"`
	ImageUrl    string         `gorm:"size:255" json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateActivityRequest representa una solicitud para crear una actividad
type CreateActivityRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Day         string `form:"day" binding:"required"`
	Time        string `form:"time" binding:"required"`
	Duration    int    `form:"duration" binding:"required"`
	Instructor  string `form:"instructor" binding:"required"`
	Category    string `form:"category" binding:"required"`
	Capacity    int    `form:"capacity" binding:"required"`
	Image      *multipart.FileHeader `form:"image"`
}

// UpdateActivityRequest representa una solicitud para actualizar una actividad
type UpdateActivityRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Day         string `form:"day"`	
	Time        string `form:"time"`
	Duration    int    `form:"duration"`
	Instructor  string `form:"instructor"`
	Category    string `form:"category"`
	Capacity    int    `form:"capacity"`
	Image      *multipart.FileHeader `form:"image"`
}

// EnrollmentRequest representa una solicitud para inscribirse en una actividad
type EnrollmentRequest struct {
	ActivityID uint `json:"activity_id" binding:"required"`
}