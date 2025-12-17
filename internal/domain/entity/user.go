package entity

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/internal/domain/valueobject"
)

const (
	CollectionFavorites = "favorite"
	CollectionWatchlist = "watchlist"
)

type User struct {
	id        int
	name      string
	email     valueobject.Email
	password  valueobject.Password
	favorites []int
	watchlist []int
}

func NewUser(name string, email valueobject.Email, password valueobject.Password) (*User, error) {
	if name == "" {
		return nil, repository.ErrNameRequired
	}

	return &User{
		name:      name,
		email:     email,
		password:  password,
		favorites: make([]int, 0),
		watchlist: make([]int, 0),
	}, nil
}

func ReconstructUser(id int, name string, email string, favorites, watchlist []int) (*User, error) {
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        id,
		name:      name,
		email:     emailVO,
		favorites: favorites,
		watchlist: watchlist,
	}, nil
}

func (u *User) ID() int                        { return u.id }
func (u *User) Name() string                   { return u.name }
func (u *User) Email() valueobject.Email       { return u.email }
func (u *User) EmailString() string            { return u.email.String() }
func (u *User) Password() valueobject.Password { return u.password }

func (u *User) Favorites() []int {
	result := make([]int, len(u.favorites))
	copy(result, u.favorites)
	return result
}

func (u *User) Watchlist() []int {
	result := make([]int, len(u.watchlist))
	copy(result, u.watchlist)
	return result
}

func (u *User) AddToFavorites(movieID int) error {
	if u.IsInFavorites(movieID) {
		return repository.ErrMovieAlreadyInFavorites
	}
	u.favorites = append(u.favorites, movieID)
	return nil
}

func (u *User) RemoveFromFavorites(movieID int) error {
	for i, id := range u.favorites {
		if id == movieID {
			u.favorites = append(u.favorites[:i], u.favorites[i+1:]...)
			return nil
		}
	}
	return repository.ErrMovieNotInFavorites
}

func (u *User) AddToWatchlist(movieID int) error {
	if u.IsInWatchlist(movieID) {
		return repository.ErrMovieAlreadyInWatchlist
	}
	u.watchlist = append(u.watchlist, movieID)
	return nil
}

func (u *User) RemoveFromWatchlist(movieID int) error {
	for i, id := range u.watchlist {
		if id == movieID {
			u.watchlist = append(u.watchlist[:i], u.watchlist[i+1:]...)
			return nil
		}
	}
	return repository.ErrMovieNotInWatchlist
}

func (u *User) IsInFavorites(movieID int) bool {
	for _, id := range u.favorites {
		if id == movieID {
			return true
		}
	}
	return false
}

func (u *User) IsInWatchlist(movieID int) bool {
	for _, id := range u.watchlist {
		if id == movieID {
			return true
		}
	}
	return false
}

func (u *User) IsInCollection(movieID int, collectionType string) bool {
	switch collectionType {
	case CollectionFavorites:
		return u.IsInFavorites(movieID)
	case CollectionWatchlist:
		return u.IsInWatchlist(movieID)
	default:
		return false
	}
}

func (u *User) AddToCollection(movieID int, collectionType string) error {
	switch collectionType {
	case CollectionFavorites:
		return u.AddToFavorites(movieID)
	case CollectionWatchlist:
		return u.AddToWatchlist(movieID)
	default:
		return repository.ErrInvalidCollectionType
	}
}

func (u *User) FavoritesCount() int  { return len(u.favorites) }
func (u *User) WatchlistCount() int  { return len(u.watchlist) }
func (u *User) HasCollections() bool { return len(u.favorites) > 0 || len(u.watchlist) > 0 }
