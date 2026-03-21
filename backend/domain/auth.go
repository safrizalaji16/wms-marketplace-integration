package domain

import (
	"backend/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthLoginRequest) (dto.AuthResponse, error)
	Register(ctx context.Context, req dto.AuthRegisterRequest) (dto.UserData, error)
}
