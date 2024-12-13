package repositories

import (
	"meta-node-ficam/internal/model"
	"time"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmailRepository interface {
	SaveEmailVerification(email, codeEmail string, expiredTime time.Time) error
	CheckEmailVerification(email, codeEmail string, nowTime time.Time) (bool, error)
	DeleteEmailVerification(email, codeEmail string) error
}

type emailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) EmailRepository {
	return &emailRepository{db}
}

func (repo *emailRepository) SaveEmailVerification(
	email, codeEmail string, expiredTime time.Time) error {
	return repo.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&model.EmailVerification{
		Email:       email,
		CodeEmail:   codeEmail,
		ExpiredTime: expiredTime}).Error
}

func (repo *emailRepository) CheckEmailVerification(
	email, codeEmail string, nowTime time.Time) (bool, error) {
	var history *model.EmailVerification
	result := repo.db.Model(&model.EmailVerification{}).
		Where("email = ? AND codeEmail = ? AND expiredTime > ?", email, codeEmail, nowTime).
		First(&history)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error    
	}
	updateResult := repo.db.Model(&history).Update("VerifyEmail", true)
	if updateResult.Error != nil {
		return false, updateResult.Error    
	}
	
	return updateResult.RowsAffected > 0, nil
}

func (repo *emailRepository) DeleteEmailVerification(email, codeEmail string) error {
	return repo.db.Unscoped().Where("email = ? AND codeEmail = ?", email, codeEmail).
		Delete(&model.EmailVerification{}).Error
}

