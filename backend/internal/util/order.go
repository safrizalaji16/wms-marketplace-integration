package util

import (
	"database/sql"
	"strings"
	"time"
)

func NormalizeStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func IsSupportedMarketplaceWebhookStatus(status string) bool {
	switch status {
	case "processing", "paid", "shipping", "delivered", "cancelled":
		return true
	default:
		return false
	}
}

func IsTerminalMarketplaceStatus(status string) bool {
	switch status {
	case "delivered", "cancelled":
		return true
	default:
		return false
	}
}

func IsSupportedShippingWebhookStatus(status string) bool {
	switch status {
	case "label_created", "awaiting_pickup", "shipped", "delivered", "cancelled":
		return true
	default:
		return false
	}
}

func IsTerminalShippingStatus(status string) bool {
	switch status {
	case "delivered", "cancelled":
		return true
	default:
		return false
	}
}

func NullableString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

func NullableTime(value sql.NullTime) string {
	if value.Valid {
		return value.Time.Format(time.RFC3339)
	}
	return ""
}
