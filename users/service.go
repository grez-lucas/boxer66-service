package users

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/grez-lucas/boxer66-service/internal/repository"
	"github.com/grez-lucas/boxer66-service/middleware"
	"golang.org/x/crypto/bcrypt"
)

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenLength = 5
)

type UserService struct {
	ctx           context.Context
	repository    *repository.Queries
	passwordCache sync.Map
}

func NewUserService(
	ctx context.Context,
	repository *repository.Queries,
) *UserService {
	return &UserService{
		ctx:           ctx,
		repository:    repository,
		passwordCache: sync.Map{},
	}
}

var (
	ErrUserDoesntExist = errors.New("user does not exist")
	ErrInvalidPassword = errors.New("password is invalid")
	ErrInvalidToken    = errors.New("token is invalid")
	ErrTokenIsExpired  = errors.New("token is expired")
)

func (s *UserService) GetUsers() ([]repository.User, error) {
	users, err := s.repository.GetAllUsers(s.ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUserByEmail(email string) (*repository.User, error) {
	user, err := s.repository.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserDoesntExist
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserService) CreateUser(email, password string) (*repository.User, error) {
	// Encrypt the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.CreateUser(s.ctx, repository.CreateUserParams{
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) Login(email string, requestPassword string) (*repository.User, string, error) {
	// Get the user
	user, err := s.repository.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrUserDoesntExist
		}
		return nil, "", err
	}

	// Compare his request's password vs the hashedpassword
	if err := comparePassword(user.Password, requestPassword); err != nil {
		return nil, "", ErrInvalidPassword
	}

	token, err := middleware.CreateJWT(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *UserService) Register(email, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	verificationToken, err := generateUniqueToken()
	if err != nil {
		return fmt.Errorf("failed to generate unique token: %w", err)
	}
	cacheKey := generateCacheKey(email)
	s.passwordCache.Store(cacheKey, hashedPassword)

	// Save token in DB
	token, err := s.repository.CreateEmailVerificationToken(s.ctx, repository.CreateEmailVerificationTokenParams{
		Email:                  email,
		VerificationToken:      verificationToken,
		HashedPasswordCacheKey: cacheKey,
		ExpiresAt:              time.Now().Add(1 * time.Hour), // TODO: Export this to a const
	})
	if err != nil {
		return fmt.Errorf("failed to create email verification token in db: %w", err)
	}

	// TODO: Send the token to the user via email service
	slog.Info("Sending token via email", slog.String("token", token.VerificationToken), slog.String("email", email))

	return nil
}

func (s *UserService) VerifyEmailToken(email, token string) (*repository.User, string, error) {
	// 1. Query the email_verification_tokens table for a matching email and verification_token_key
	dbToken, err := s.repository.GetEmailVerificationTokenByEmail(s.ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get email verification token for email %s: %w", email, err)
	}

	// 2.1 Check that the token is valid
	if token != dbToken.VerificationToken {
		slog.Error("Token is invalid", slog.String("user_token", token), slog.String("db_token", dbToken.VerificationToken))
		return nil, "", ErrInvalidToken
	}

	cacheKey := generateCacheKey(email)

	// 2.2 Verify that the token hasn't expired
	if dbToken.ExpiresAt.Before(time.Now()) {
		// Delete the temporary password
		s.passwordCache.Delete(cacheKey)
		return nil, "", ErrTokenIsExpired
	}

	// 3 Get the password from the cache for said email
	hashedPassword, ok := s.passwordCache.Load(cacheKey)
	if !ok {
		return nil, "", errors.New("password not found in cache")
	}
	hashedPasswordStr, ok := hashedPassword.(string)
	if !ok {
		return nil, "", errors.New("stored hashedPassword value is not of type string")
	}

	// 4. Create a user in the users table with the memory password
	user, err := s.repository.CreateUser(s.ctx, repository.CreateUserParams{
		Email:    email,
		Password: []byte(hashedPasswordStr),
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user in db: %w", err)
	}

	// 4. Delete the record from the email_verification_tokens table.
	if err := s.repository.DeleteEmailVerificationTokenByID(s.ctx, dbToken.ID); err != nil {
		return nil, "", fmt.Errorf("failed to delete email verification token from db: %w", err)
	}

	// 4.1 Delete the record from the in memory map
	s.passwordCache.Delete(cacheKey)

	jwt, err := middleware.CreateJWT(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, jwt, nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func comparePassword(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err
}

func generateUniqueToken() (string, error) {
	b := make([]byte, 5)
	for i := range b {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		b[i] = letterBytes[nBig.Int64()]
	}
	return string(b), nil
}

func generateCacheKey(email string) string {
	hasher := sha256.New()
	hasher.Write([]byte(email))
	hashBytes := hasher.Sum(nil) // Get the resulting hash as a byte slice
	return "verification_code:" + hex.EncodeToString(hashBytes)
}
