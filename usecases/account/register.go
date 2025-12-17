package account

import (
	"github.com/jgamaraalv/movies.git/domain/entity"
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
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	password, err := valueobject.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Use domain entity to validate and create user
	user, err := entity.NewUser(input.Name, email, password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := user.Password().Hash()
	if err != nil {
		uc.logger.Error("Failed to hash password", err)
		return nil, valueobject.ErrPasswordMismatch
	}

	success, err := uc.userRepo.Register(user.Name(), user.EmailString(), hashedPassword)
	if err != nil {
		return nil, err
	}

	jwt := token.CreateJWT(
		models.User{Email: user.EmailString(), Name: user.Name()},
		*uc.logger,
	)

	uc.logger.Info("User registered successfully: " + user.EmailString())

	return &RegisterOutput{
		Success: success,
		Message: "User registered successfully",
		JWT:     jwt,
	}, nil
}
