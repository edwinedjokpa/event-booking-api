package user

import "gorm.io/gorm"

type UserRepository interface {
	Create(user User) error
	FindOneByID(userId string) (*User, error)
	Update(user User) error
	FindOneByEmail(email string) (*User, error)
	UpdatePassword(userID, hashedPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) Create(user User) error {
	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}
	return nil

}

func (repo *userRepository) FindOneByID(userID string) (*User, error) {
	var user User
	if err := repo.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) Update(user User) error {
	if err := repo.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) FindOneByEmail(email string) (*User, error) {
	var user User
	if err := repo.db.Model(&User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) UpdatePassword(userID string, hashedPassword string) error {
	result := repo.db.Model(&User{}).Where("id = ?", userID).Update("password", hashedPassword)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
