package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"app/hasher"
	"app/user/queries"
	"app/validator"
)

var (
	ErrInvalidEmail           = errors.New("invalid email")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrInvalidFullName        = errors.New("invalid name")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrNoSuchUser             = errors.New("such user does not exist")
	ErrInvalidCredential      = errors.New("invalid credentials")
	ErrInvalidCurrentPassword = errors.New("invalid current password")
	ErrInvalidNewPassword     = errors.New("invalid new password")
)

// userRepository the dao interface for user service
type userRepository interface {
	// RegisterUser will create a user in database
	RegisterUser(ctx context.Context, param queries.RegisterUserParams) error

	// GetUser will get an user by email
	GetUser(ctx context.Context, email string) (queries.User, error)

	// UpdateUserPassword will update the user's password by their email and new password hash
	UpdateUserPassword(ctx context.Context, arg queries.UpdateUserPasswordParams) error
}

// Service will handle all the business logic in user domain
type Service struct {
	Repo   userRepository // Repo the dao interface
	Logger *zap.Logger
}

func NewService(r userRepository, l *zap.Logger) Service {
	return Service{
		Repo:   r,
		Logger: l,
	}
}

func (s Service) RegisterUser(req queries.RegisterUserRequest) error {
	if validator.IsValidEmail(req.Email) {
		return ErrInvalidEmail
	}

	if !validator.IsValidPassword(req.Password) {
		return fmt.Errorf("%w: the password must contain at-lest one uppercase, lowercase,digit,special character \"@$!%%*#?&^_-\" and must be of at-least 8 characters", ErrInvalidPassword)
	}

	const minFullName = 3

	if len(req.FullName) < minFullName {
		return fmt.Errorf("%w: the fullname must contain at-least %d characters", ErrInvalidFullName, minFullName)
	}

	ctx := context.TODO()

	_, err := s.Repo.GetUser(ctx, req.Email)

	if err == nil {
		return fmt.Errorf("%w: user with given email already exists", ErrUserAlreadyExists)
	}

	hashedPass, err := hasher.HashAndSalt(req.Password)
	if err != nil {
		return err
	}

	q := req.ToRegisterUserParams(hashedPass)

	return s.Repo.RegisterUser(ctx, q)
}

func (s Service) GetUserByEmail(email string) (queries.User, error) {
	user, err := s.Repo.GetUser(context.TODO(), email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return user, ErrNoSuchUser
	default:
		return user, err
	}
}

func (s Service) GetUser(email, password string) (queries.User, error) {
	user, err := s.Repo.GetUser(context.TODO(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.User{}, ErrInvalidCredential
		}
		return queries.User{}, err
	}

	if err = hasher.MatchHashedPassword(user.PasswordHash, password); err != nil {
		return queries.User{}, ErrInvalidCredential
	}

	return user, nil
}

func (s Service) ResetPassword(email, currentPass, newPass, confirmPass string) error {
	user, err := s.Repo.GetUser(context.TODO(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidCredential
		}

		return err
	}
	if err = hasher.MatchHashedPassword(user.PasswordHash, currentPass); err != nil {
		return ErrInvalidCurrentPassword
	}

	if currentPass == newPass {
		return fmt.Errorf("%w: new password should not be same as current password", ErrInvalidNewPassword)
	}

	if !validator.IsValidPassword(newPass) {
		return fmt.Errorf("%w: the new password must contain at-lest one uppercase, lowercase,digit,special character \"@$!%%*#?&^_-\" and must be of at-least 8 characters", ErrInvalidNewPassword)
	}

	if newPass != confirmPass {
		return fmt.Errorf("%w: confirm password should be same as new password ", ErrInvalidNewPassword)
	}

	newPassHash, err := hasher.HashAndSalt(newPass)
	if err != nil {
		return err
	}

	return s.Repo.UpdateUserPassword(context.TODO(), queries.UpdateUserPasswordParams{
		Email:        email,
		PasswordHash: newPassHash,
	})
}
