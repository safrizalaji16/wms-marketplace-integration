package service

import (
	"backend/domain"
	"backend/dto"
	"backend/internal/config"
	"backend/internal/util"
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	conf           *config.Config
	userRepository domain.UserRepository
}

func NewAuth(cnf *config.Config, userRepository domain.UserRepository) domain.AuthService {
	return authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}

func (a authService) Login(ctx context.Context, req dto.AuthLoginRequest) (dto.AuthResponse, error) {
	user, err := a.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("login failed while fetching user", "email", req.Email, "error", err.Error())
		return dto.AuthResponse{}, err
	}

	if user.ID == "" {
		slog.Warn("login rejected: user not found", "email", req.Email)
		return dto.AuthResponse{}, util.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		slog.Warn("login rejected: invalid password", "email", req.Email)
		return dto.AuthResponse{}, util.ErrUnauthorized
	}

	claim := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		slog.Error("login failed while signing token", "user_id", user.ID, "error", err.Error())
		return dto.AuthResponse{}, util.ErrUnauthorized
	}

	slog.Info("login succeeded", "user_id", user.ID, "email", user.Email, "role", user.Role)

	return dto.AuthResponse{
		Token: t,
	}, nil
}

func (a authService) Register(ctx context.Context, req dto.AuthRegisterRequest) (dto.UserData, error) {
	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	_, err := a.userRepository.GetByEmail(ctx, req.Email)
	if err == nil {
		slog.Warn("register rejected: user already exists", "email", req.Email)
		return dto.UserData{}, util.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("register failed while hashing password", "email", req.Email, "error", err.Error())
		return dto.UserData{}, err
	}
	user.Password = string(hashedPassword)

	createdUser, err := a.userRepository.Create(ctx, &user)
	if err != nil {
		slog.Error("register failed while creating user", "email", req.Email, "role", req.Role, "error", err.Error())
		return dto.UserData{}, err
	}

	slog.Info("register succeeded", "user_id", createdUser.ID, "email", createdUser.Email, "role", createdUser.Role)

	return dto.UserData{
		ID:       createdUser.ID,
		Username: createdUser.Email,
		Role:     createdUser.Role,
	}, nil
}
