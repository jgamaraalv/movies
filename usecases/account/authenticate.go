package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
	"github.com/jgamaraalv/movies.git/token"
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
	accountStorage providers.AccountStorage
	logger         *logger.Logger
}

func NewAuthenticateUseCase(storage providers.AccountStorage, log *logger.Logger) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		accountStorage: storage,
		logger:         log,
	}
}

func (uc *AuthenticateUseCase) Execute(input AuthenticateInput) (*AuthenticateOutput, error) {
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	success, err := uc.accountStorage.Authenticate(input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	jwt := token.CreateJWT(
		models.User{Email: input.Email},
		*uc.logger,
	)

	uc.logger.Info("User authenticated successfully: " + input.Email)

	return &AuthenticateOutput{
		Success: success,
		Message: "User authenticated successfully",
		JWT:     jwt,
	}, nil
}

func (uc *AuthenticateUseCase) validateInput(input AuthenticateInput) error {
	if input.Email == "" {
		return errors.New("email is required")
	}
	if input.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
