package user

type UserService interface {
	Dashboard(userID string) (*User, error)
}

type userService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) UserService {
	return &userService{repository}
}

func (svc *userService) Dashboard(userID string) (*User, error) {
	user, err := svc.repository.FindOneByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
