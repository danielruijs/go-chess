package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"go-chess/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string
	DisplayName    string
	HashedPassword []byte
}

type UserStore struct {
	queries *sqlc.Queries
}

func NewUserStore(queries *sqlc.Queries) (*UserStore, error) {
	return &UserStore{queries: queries}, nil
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

func validateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return NewUserRegistrationError("username must be between 3 and 20 characters and contain only letters, numbers, hyphens, or underscores")
	}
	return nil
}

func validateDisplayName(displayName string) error {
	length := len([]rune(displayName))
	if length < 3 || length > 30 {
		return NewUserRegistrationError("display name must be between 3 and 30 characters")
	}
	return nil
}

func validatePassword(password string) error {
	length := len([]rune(password))
	if length < 6 || length > 72 {
		return NewUserRegistrationError("password must be between 6 and 72 characters")
	}
	hasUpper := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return NewUserRegistrationError("password must contain at least one uppercase letter")
	}
	return nil
}

func (s *UserStore) Register(ctx context.Context, username, password, displayName string) (User, error) {
	if err := validateUsername(username); err != nil {
		return User{}, err
	}
	if err := validateDisplayName(displayName); err != nil {
		return User{}, err
	}
	if err := validatePassword(password); err != nil {
		return User{}, err
	}

	lowerUsername := strings.ToLower(username)
	exists, err := s.queries.UserExists(ctx, lowerUsername)
	if err != nil {
		return User{}, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if exists {
		return User{}, NewUserRegistrationError("username is already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	dbUser, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username:       lowerUsername,
		DisplayName:    displayName,
		HashedPassword: hashed,
	})
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return User{
		Username:       dbUser.Username,
		DisplayName:    dbUser.DisplayName,
		HashedPassword: dbUser.HashedPassword,
	}, nil
}

func (s *UserStore) Login(ctx context.Context, username, password string) (User, error) {
	lowerUsername := strings.ToLower(username)
	dbUser, err := s.queries.GetUserByUsername(ctx, lowerUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, fmt.Errorf("failed to query user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(dbUser.HashedPassword, []byte(password))
	if err != nil {
		return User{}, ErrInvalidCredentials
	}

	return User{
		Username:       dbUser.Username,
		DisplayName:    dbUser.DisplayName,
		HashedPassword: dbUser.HashedPassword,
	}, nil
}
