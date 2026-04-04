// user service

package service

import (
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) List() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateRole(id string, role models.Role) error {
	return s.userRepo.UpdateRole(id, role)
}

func (s *UserService) UpdateStatus(id string, active bool) error {
	return s.userRepo.UpdateStatus(id, active)
}
