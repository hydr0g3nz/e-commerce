package dto

import "github.com/hydr0g3nz/e-commerce/internal/core/domain"

type RegisterResponse struct {
	User         domain.User          `json:"user"`
	TokenDetails domain.TokenResponse `json:"token_details"`
}
type LoginResponse struct {
	User         domain.User          `json:"user"`
	TokenDetails domain.TokenResponse `json:"token_details"`
}
