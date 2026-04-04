// All business side logic are in service files

// Auth Service Code, contains business logic used regarding to authentication service

package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/repository"
	jwtpkg "github.com/suthar345Piyush/financedashboard/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrUserInactive = errors.New("account is inactive")
var ErrDuplicateEmail = repository.ErrDuplicateEmail

// struct for the auth service

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *jwtpkg.Manager
}

func NewAuthService(userRepo *repository.UserRepository, jwt *jwtpkg.Manager) *AuthService {

	return &AuthService{userRepo: userRepo, jwtManager: jwt}

}

// user registration

func (s *AuthService) Register(email, password string) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
		Role:     models.RoleViewer,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// login the user

func (s *AuthService) Login(email, password string) (string, *models.User, error) {

	user, err := s.userRepo.GetByEmail(email)

	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			return "", nil, ErrInvalidCredentials
		}

		return "", nil, err
	}

	if !user.IsActive {
		return "", nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))

	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
