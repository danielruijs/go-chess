package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"go-chess/internal/cache"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string
	DisplayName    string
	HashedPassword []byte
}

type UserStore struct {
	cache *cache.Cache[string, User] // username (lowercase) -> User
}

func NewUserStore() (*UserStore, error) {
	cache, err := cache.New[string](cache.Options[User]{})
	if err != nil {
		return nil, err
	}
	return &UserStore{cache: cache}, nil
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

func validateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return errors.New("username must be between 3 and 20 characters and contain only letters, numbers, hyphens, or underscores")
	}
	return nil
}

func validateDisplayName(displayName string) error {
	length := len([]rune(displayName))
	if length < 3 || length > 30 {
		return errors.New("display name must be between 3 and 30 characters")
	}
	return nil
}

func validatePassword(password string) error {
	length := len([]rune(password))
	if length < 6 || length > 72 {
		return errors.New("password must be between 6 and 72 characters")
	}
	hasUpper := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	return nil
}

func (s *UserStore) Register(username, password, displayName string) (User, error) {
	if err := validateUsername(username); err != nil {
		return User{}, err
	}
	if err := validateDisplayName(displayName); err != nil {
		return User{}, err
	}
	if err := validatePassword(password); err != nil {
		return User{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := User{
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashed,
	}

	lowerUsername := strings.ToLower(username)
	wasSet := s.cache.SetIfAbsent(lowerUsername, user)
	if !wasSet {
		return User{}, errors.New("username is already taken")
	}

	return user, nil
}

func (s *UserStore) Login(username, password string) (User, error) {
	lowerUsername := strings.ToLower(username)
	user, exists := s.cache.Get(lowerUsername)
	if !exists {
		return User{}, errors.New("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		return User{}, errors.New("invalid username or password")
	}

	return user, nil
}
