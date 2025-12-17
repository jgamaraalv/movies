package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
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
	accountStorage providers.AccountStorage
	logger         *logger.Logger
}

func NewRegisterUseCase(storage providers.AccountStorage, log *logger.Logger) *RegisterUseCase {
	return &RegisterUseCase{
		accountStorage: storage,
		logger:         log,
	}
}

func (uc *RegisterUseCase) Execute(input RegisterInput) (*RegisterOutput, error) {
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	success, err := uc.accountStorage.Register(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	jwt := token.CreateJWT(
		models.User{Email: input.Email, Name: input.Name},
		*uc.logger,
	)

	uc.logger.Info("User registered successfully: " + input.Email)

	return &RegisterOutput{
		Success: success,
		Message: "User registered successfully",
		JWT:     jwt,
	}, nil
}

func (uc *RegisterUseCase) validateInput(input RegisterInput) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.Email == "" {
		return errors.New("email is required")
	}
	if input.Password == "" {
		return errors.New("password is required")
	}
	if len(input.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
