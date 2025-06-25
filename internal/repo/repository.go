package repo

import (
	"errors"
	"fmt"

	"github.com/Koyo-os/form-crud-service/internal/entity"
	"github.com/Koyo-os/form-crud-service/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Create(form *entity.Form) error {
	if form == nil {
		return errors.New("form is nil")
	}
	if err := r.db.Create(form).Error; err != nil {
		r.logger.Error("failed to create form", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) Get(id string) (*entity.Form, error) {
	var form entity.Form
	if err := r.db.Where("id = ?", id).First(&form).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("failed to get form by id",
			zap.String("id", id),
			zap.Error(err))

		return nil, err
	}

	if err := r.db.Model(&form).Association("Questions").Find(&form.Questions); err != nil {
		r.logger.Warn("failed to load questions for form",
			zap.String("id", id),
			zap.Error(err))
	}
	return &form, nil
}

func (r *Repository) Update(id string, key string, value interface{}) error {
	updates := map[string]interface{}{key: value}
	if err := r.db.Model(&entity.Form{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		r.logger.Error("failed to update form",
			zap.String("key", key),
			zap.String("id", id),
			zap.Error(err))

		return err
	}
	return nil
}

func (r *Repository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Form{}).Error; err != nil {
		r.logger.Error("failed to delete form",
			zap.String("id", id),
			zap.Error(err))

		return err
	}
	return nil
}

func (r *Repository) GetMore(key string, value interface{}) ([]entity.Form, error) {
	var forms []entity.Form
	query := fmt.Sprintf("%s = ?", key)
	if err := r.db.Where(query, value).Find(&forms).Error; err != nil {
		r.logger.Error("failed to get forms by",
			zap.String("key", key))
			
		return nil, err
	}
	return forms, nil
}
