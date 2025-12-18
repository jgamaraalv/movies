package account

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/internal/domain/valueobject"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
	"github.com/jgamaraalv/movies.git/pkg/token"
)

type AuthenticateInput struct {
	Email    string
	Password string
}

type AuthenticateOutput struct {
	Success bool
	Message string
	JWT     string
}

type AuthenticateUseCase struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewAuthenticateUseCase(repo repository.UserRepository, log *logger.Logger) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *AuthenticateUseCase) Execute(input AuthenticateInput) (*AuthenticateOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if input.Password == "" {
		return nil, valueobject.ErrEmptyPassword
	}

	success, err := uc.userRepo.Authenticate(email.String(), input.Password)
	if err != nil {
		return nil, err
	}

	jwt := token.CreateJWT(
		models.User{Email: email.String()},
		*uc.logger,
	)

	uc.logger.Info("User authenticated successfully: " + email.String())

	return &AuthenticateOutput{
		Success: success,
		Message: "User authenticated successfully",
		JWT:     jwt,
	}, nil
}
