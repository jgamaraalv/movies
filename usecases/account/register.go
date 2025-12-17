package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/domain/valueobject"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/token"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterOutput struct {
	Success bool
	Message string
	JWT     string
}

type RegisterUseCase struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewRegisterUseCase(repo repository.UserRepository, log *logger.Logger) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *RegisterUseCase) Execute(input RegisterInput) (*RegisterOutput, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	password, err := valueobject.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		uc.logger.Error("Failed to hash password", err)
		return nil, errors.New("failed to process password")
	}

	success, err := uc.userRepo.Register(input.Name, email.String(), hashedPassword)
	if err != nil {
		return nil, err
	}

	jwt := token.CreateJWT(
		models.User{Email: email.String(), Name: input.Name},
		*uc.logger,
	)

	uc.logger.Info("User registered successfully: " + email.String())

	return &RegisterOutput{
		Success: success,
		Message: "User registered successfully",
		JWT:     jwt,
	}, nil
}
