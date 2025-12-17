package repository

import "github.com/jgamaraalv/movies.git/models"

type UserRepository interface {
	Register(name string, email string, hashedPassword string) (bool, error)
	Authenticate(email string, password string) (bool, error)
	GetAccountDetails(email string) (models.User, error)
	SaveCollection(user models.User, movieID int, collectionType string) (bool, error)
}
