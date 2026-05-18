package services

import (
	"website-starter/models"

	"gorm.io/gorm"
)

type HomeService struct {
	db *gorm.DB
}

func NewHomeService(db *gorm.DB) *HomeService {
	return &HomeService{
		db: db,
	}
}

func (service *HomeService) GetWelcomeMessage() string {
	return "Welcome to the website starter! (Home)"
}

func (service *HomeService) GetNotes() ([]models.Note, error) {
	var notes []models.Note

	err := service.db.
		Order(
			"created_at DESC",
		).
		Find(
			&notes,
		).
		Error

	return notes, err
}

func (service *HomeService) CreateNote(
	title,
	description string,
) error {
	note := models.Note{
		Title:       title,
		Description: description,
	}

	return service.db.Create(
		&note,
	).Error
}

func (service *HomeService) DeleteNote(
	id string,
) error {
	return service.db.Delete(
		&models.Note{},
		id,
	).Error
}

func (service *HomeService) UpdateNote(id, title, description string) error {
	return service.db.Model(
		&models.Note{},
	).Where(
		"id = ?",
		id,
	).Updates(
		models.Note{
			Title:       title,
			Description: description,
		}).Error
}
