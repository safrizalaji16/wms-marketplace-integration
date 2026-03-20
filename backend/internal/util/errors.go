package util

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized                 = errors.New("unauthorized")
	ErrUserAlreadyExists            = errors.New("user already exists")
	ErrOrderNotFound                = errors.New("order not found")
	ErrInvalidRequestBody           = errors.New("invalid request body")
	ErrInvalidStatusTransition      = errors.New("invalid status transition")
	ErrUnsupportedMarketplaceStatus = errors.New("unsupported marketplace status")
	ErrUnsupportedShippingStatus    = errors.New("unsupported shipping status")
	ErrOrderSNRequired              = errors.New("order_sn is required")
	ErrStatusRequired               = errors.New("status is required")
	ErrChannelIDRequired            = errors.New("channel_id is required")
	ErrMarketplaceNotConfigured     = errors.New("marketplace repository is not configured")
)

func HTTPStatusFromError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrInvalidRequestBody):
		return http.StatusUnprocessableEntity
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrUserAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, ErrOrderNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrOrderSNRequired),
		errors.Is(err, ErrStatusRequired),
		errors.Is(err, ErrChannelIDRequired),
		errors.Is(err, ErrUnsupportedMarketplaceStatus),
		errors.Is(err, ErrUnsupportedShippingStatus),
		errors.Is(err, ErrInvalidStatusTransition):
		return http.StatusBadRequest
	}

	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "marketplace auth code is required"),
		strings.Contains(message, "marketplace refresh token is required"),
		strings.Contains(message, "marketplace partner key is required"),
		strings.Contains(message, "marketplace base url is required"):
		return http.StatusBadRequest
	case strings.Contains(message, "invalid credentials"),
		strings.Contains(message, "invalid refresh token"),
		strings.Contains(message, "invalid access token"),
		strings.Contains(message, "unauthorized"):
		return http.StatusUnauthorized
	case strings.Contains(message, "marketplace request failed"),
		strings.Contains(message, "marketplace order list failed"):
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
