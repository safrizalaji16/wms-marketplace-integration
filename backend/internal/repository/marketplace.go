package repository

import (
	"bytes"
	"backend/domain"
	"backend/dto"
	"backend/internal/config"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"time"
)

type MarketplaceRepository struct {
	conf      *config.Config
	client    *http.Client
	tokenRepo domain.MarketplaceTokenRepository
}

const (
	maxMarketplaceAttempts = 4
	initialBackoff         = 500 * time.Millisecond
)

func NewMarketplace(cnf *config.Config, tokenRepo domain.MarketplaceTokenRepository) domain.MarketplaceRepository {
	return &MarketplaceRepository{
		conf:      cnf,
		client:    &http.Client{Timeout: 15 * time.Second},
		tokenRepo: tokenRepo,
	}
}

func (m *MarketplaceRepository) Bootstrap(ctx context.Context, state string) (dto.MarketplaceTokenData, error) {
	slog.Info("marketplace bootstrap started", "shop_id", m.conf.Marketplace.ShopID, "state", state)
	authorizeData, err := m.Authorize(ctx, state)
	if err != nil {
		slog.Error("marketplace bootstrap failed during authorize", "error", err.Error())
		return dto.MarketplaceTokenData{}, err
	}

	tokenData, err := m.ExchangeToken(ctx, authorizeData.Code)
	if err != nil {
		slog.Error("marketplace bootstrap failed during token exchange", "error", err.Error())
		return dto.MarketplaceTokenData{}, err
	}

	slog.Info("marketplace bootstrap succeeded", "shop_id", m.conf.Marketplace.ShopID)
	return tokenData, nil
}

func (m *MarketplaceRepository) Authorize(ctx context.Context, state string) (dto.MarketplaceAuthorizeData, error) {
	if m.conf.Marketplace.PartnerKey == "" {
		return dto.MarketplaceAuthorizeData{}, fmt.Errorf("marketplace partner key is required")
	}
	if state == "" {
		state = "pm"
	}

	timestamp := time.Now().Unix()
	apiPath := "/oauth/authorize"
	base := m.conf.Marketplace.PartnerID + apiPath + strconv.FormatInt(timestamp, 10) + m.conf.Marketplace.ShopID
	sign := signHMACSHA256(m.conf.Marketplace.PartnerKey, base)

	u, err := url.Parse(strings.TrimRight(m.conf.Marketplace.BaseURL, "/") + apiPath)
	if err != nil {
		return dto.MarketplaceAuthorizeData{}, err
	}

	query := u.Query()
	query.Set("shop_id", m.conf.Marketplace.ShopID)
	query.Set("state", state)
	query.Set("partner_id", m.conf.Marketplace.PartnerID)
	query.Set("timestamp", strconv.FormatInt(timestamp, 10))
	query.Set("sign", sign)
	query.Set("redirect", m.conf.Marketplace.RedirectURL)
	u.RawQuery = query.Encode()

	var response struct {
		Data dto.MarketplaceAuthorizeData `json:"data"`
	}
	if err := m.doRetriableJSON(ctx, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		return req, nil
	}, &response); err != nil {
		return dto.MarketplaceAuthorizeData{}, err
	}

	if _, err := m.tokenRepo.Upsert(ctx, &domain.MarketplaceToken{
		ShopID:   m.conf.Marketplace.ShopID,
		AuthCode: sqlNullString(response.Data.Code),
	}); err != nil {
		return dto.MarketplaceAuthorizeData{}, err
	}

	slog.Info("marketplace authorize succeeded", "shop_id", m.conf.Marketplace.ShopID)

	return response.Data, nil
}

func (m *MarketplaceRepository) ExchangeToken(ctx context.Context, code string) (dto.MarketplaceTokenData, error) {
	if code == "" {
		storedToken, err := m.tokenRepo.GetByShopID(ctx, m.conf.Marketplace.ShopID)
		if err != nil {
			return dto.MarketplaceTokenData{}, err
		}
		if storedToken.AuthCode.Valid {
			code = storedToken.AuthCode.String
		}
	}
	if code == "" {
		return dto.MarketplaceTokenData{}, fmt.Errorf("marketplace auth code is required")
	}

	tokenData, err := m.requestToken(ctx, map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}, code)
	if err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	if _, err := m.tokenRepo.Upsert(ctx, &domain.MarketplaceToken{
		ShopID:       m.conf.Marketplace.ShopID,
		AuthCode:     sqlNullString(code),
		AccessToken:  sqlNullString(tokenData.AccessToken),
		RefreshToken: sqlNullString(tokenData.RefreshToken),
		ExpiresAt:    sqlNullTimeFromNow(tokenData.ExpiresIn),
	}); err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	slog.Info("marketplace token exchange succeeded", "shop_id", m.conf.Marketplace.ShopID, "expires_in", tokenData.ExpiresIn)

	return tokenData, nil
}

func (m *MarketplaceRepository) RefreshToken(ctx context.Context, refreshToken string) (dto.MarketplaceTokenData, error) {
	if refreshToken == "" {
		storedToken, err := m.tokenRepo.GetByShopID(ctx, m.conf.Marketplace.ShopID)
		if err != nil {
			return dto.MarketplaceTokenData{}, err
		}
		if storedToken.RefreshToken.Valid {
			refreshToken = storedToken.RefreshToken.String
		}
	}
	if refreshToken == "" {
		return dto.MarketplaceTokenData{}, fmt.Errorf("marketplace refresh token is required")
	}

	tokenData, err := m.requestToken(ctx, map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}, refreshToken)
	if err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	if _, err := m.tokenRepo.Upsert(ctx, &domain.MarketplaceToken{
		ShopID:       m.conf.Marketplace.ShopID,
		AccessToken:  sqlNullString(tokenData.AccessToken),
		RefreshToken: sqlNullString(tokenData.RefreshToken),
		ExpiresAt:    sqlNullTimeFromNow(tokenData.ExpiresIn),
	}); err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	slog.Info("marketplace refresh token succeeded", "shop_id", m.conf.Marketplace.ShopID, "expires_in", tokenData.ExpiresIn)

	return tokenData, nil
}

func (m *MarketplaceRepository) GetOrders(ctx context.Context) ([]domain.Order, error) {
	if m.conf.Marketplace.BaseURL == "" {
		return nil, fmt.Errorf("marketplace base url is required")
	}

	var wrapped struct {
		Data []dto.MarketplaceOrderData `json:"data"`
	}
	if err := m.doAuthenticatedJSON(ctx, http.MethodGet, "/order/list", nil, nil, &wrapped); err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, item := range wrapped.Data {
		order, err := m.mapMarketplaceOrder(item)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (m *MarketplaceRepository) GetOrderBySN(ctx context.Context, orderSN string) (domain.Order, error) {
	if m.conf.Marketplace.BaseURL == "" {
		return domain.Order{}, fmt.Errorf("marketplace base url is required")
	}
	if strings.TrimSpace(orderSN) == "" {
		return domain.Order{}, fmt.Errorf("order_sn is required")
	}

	u, err := url.Parse(strings.TrimRight(m.conf.Marketplace.BaseURL, "/") + "/order/detail")
	if err != nil {
		return domain.Order{}, err
	}

	query := u.Query()
	query.Set("order_sn", orderSN)
	u.RawQuery = query.Encode()

	var wrapped struct {
		Data dto.MarketplaceOrderData `json:"data"`
	}
	if err := m.doAuthenticatedJSON(ctx, http.MethodGet, "/order/detail", query, nil, &wrapped); err != nil {
		return domain.Order{}, err
	}

	return m.mapMarketplaceOrder(wrapped.Data)
}

func (m *MarketplaceRepository) GetLogisticChannels(ctx context.Context) ([]dto.MarketplaceLogisticChannelData, error) {
	if m.conf.Marketplace.BaseURL == "" {
		return nil, fmt.Errorf("marketplace base url is required")
	}

	var response struct {
		Data []dto.MarketplaceLogisticChannelData `json:"data"`
	}
	if err := m.doAuthenticatedJSON(ctx, http.MethodGet, "/logistic/channels", nil, nil, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (m *MarketplaceRepository) mapMarketplaceOrder(item dto.MarketplaceOrderData) (domain.Order, error) {
	rawPayload, err := json.Marshal(item)
	if err != nil {
		return domain.Order{}, err
	}

	var orderItems []domain.OrderItem
	for _, marketplaceItem := range item.Items {
		orderItems = append(orderItems, domain.OrderItem{
			SKU:      marketplaceItem.SKU,
			Quantity: marketplaceItem.Quantity,
			Price:    marketplaceItem.Price,
		})
	}

	return domain.Order{
		OrderSN:               item.OrderSN,
		ShopID:                item.ShopID,
		MarketplaceStatus:     item.Status,
		ShippingStatus:        item.ShippingStatus,
		WMSStatus:             "READY_TO_PICK",
		TrackingNumber:        item.TrackingNumber,
		TotalAmount:           item.TotalAmount,
		RawMarketplacePayload: rawPayload,
		CreatedAt:             sqlNullTime(item.CreatedAt),
		Items:                 orderItems,
	}, nil
}

func (m *MarketplaceRepository) ShipOrder(ctx context.Context, orderSN, channelID string) (dto.MarketplaceShipData, error) {
	if m.conf.Marketplace.BaseURL == "" {
		return dto.MarketplaceShipData{}, fmt.Errorf("marketplace base url is required")
	}
	if strings.TrimSpace(channelID) == "" {
		return dto.MarketplaceShipData{}, fmt.Errorf("channel_id is required")
	}

	payload := map[string]string{
		"order_sn":   orderSN,
		"channel_id": channelID,
	}

	var response struct {
		Data dto.MarketplaceShipData `json:"data"`
	}
	if err := m.doAuthenticatedJSON(ctx, http.MethodPost, "/logistic/ship", nil, payload, &response); err != nil {
		return dto.MarketplaceShipData{}, err
	}

	return response.Data, nil
}

func (m *MarketplaceRepository) ensureValidToken(ctx context.Context) (dto.MarketplaceTokenData, error) {
	storedToken, err := m.tokenRepo.GetByShopID(ctx, m.conf.Marketplace.ShopID)
	if err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	if !storedToken.AccessToken.Valid || storedToken.AccessToken.String == "" {
		return m.Bootstrap(ctx, "pm")
	}

	// Refresh a little early to avoid using a token that expires mid-request.
	if storedToken.ExpiresAt.Valid && time.Now().Before(storedToken.ExpiresAt.Time.Add(-1*time.Minute)) {
		return dto.MarketplaceTokenData{
			AccessToken:  storedToken.AccessToken.String,
			RefreshToken: storedToken.RefreshToken.String,
		}, nil
	}

	if storedToken.RefreshToken.Valid && storedToken.RefreshToken.String != "" {
		tokenData, err := m.RefreshToken(ctx, storedToken.RefreshToken.String)
		if err == nil {
			return tokenData, nil
		}

		// If the refresh token is stale or revoked, recover by bootstrapping a new token set.
		slog.Warn("marketplace refresh token failed, falling back to bootstrap", "shop_id", m.conf.Marketplace.ShopID, "error", err.Error())
		return m.Bootstrap(ctx, "pm")
	}

	return m.Bootstrap(ctx, "pm")
}

func (m *MarketplaceRepository) requestToken(ctx context.Context, payload map[string]string, suffix string) (dto.MarketplaceTokenData, error) {
	if m.conf.Marketplace.PartnerKey == "" {
		return dto.MarketplaceTokenData{}, fmt.Errorf("marketplace partner key is required")
	}

	timestamp := time.Now().Unix()
	apiPath := "/oauth/token"
	base := m.conf.Marketplace.PartnerID + apiPath + strconv.FormatInt(timestamp, 10) + suffix
	sign := signHMACSHA256(m.conf.Marketplace.PartnerKey, base)

	u, err := url.Parse(strings.TrimRight(m.conf.Marketplace.BaseURL, "/") + apiPath)
	if err != nil {
		return dto.MarketplaceTokenData{}, err
	}
	query := u.Query()
	query.Set("partner_id", m.conf.Marketplace.PartnerID)
	query.Set("timestamp", strconv.FormatInt(timestamp, 10))
	query.Set("sign", sign)
	u.RawQuery = query.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	var response struct {
		Data dto.MarketplaceTokenData `json:"data"`
	}
	if err := m.doRetriableJSON(ctx, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}, &response); err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	return response.Data, nil
}

func (m *MarketplaceRepository) doAuthenticatedJSON(ctx context.Context, method, path string, query url.Values, payload any, out any) error {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	for attempt := 1; attempt <= maxMarketplaceAttempts; attempt++ {
		tokenData, err := m.ensureValidToken(ctx)
		if err != nil {
			return err
		}

		fullURL := strings.TrimRight(m.conf.Marketplace.BaseURL, "/") + path
		if query != nil {
			u, err := url.Parse(fullURL)
			if err != nil {
				return err
			}
			u.RawQuery = query.Encode()
			fullURL = u.String()
		}

		var bodyReader io.Reader
		if body != nil {
			bodyReader = bytes.NewReader(body)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		statusCode, err := m.executeJSON(req, out)
		if err == nil {
			return nil
		}

		if statusCode == http.StatusUnauthorized && attempt < maxMarketplaceAttempts {
			slog.Warn("marketplace protected request got 401, retrying with recovered token", "path", path, "attempt", attempt, "error", err.Error())
			if _, recoverErr := m.recoverAccessToken(ctx); recoverErr != nil {
				return recoverErr
			}
			continue
		}

		if shouldRetryMarketplaceStatus(statusCode) && attempt < maxMarketplaceAttempts {
			slog.Warn("marketplace protected request retry scheduled", "path", path, "attempt", attempt, "status_code", statusCode, "error", err.Error())
			if sleepErr := sleepWithBackoff(ctx, attempt); sleepErr != nil {
				return sleepErr
			}
			continue
		}

		return err
	}

	return fmt.Errorf("marketplace request failed after retries")
}

func (m *MarketplaceRepository) doRetriableJSON(ctx context.Context, reqFactory func() (*http.Request, error), out any) error {
	for attempt := 1; attempt <= maxMarketplaceAttempts; attempt++ {
		req, err := reqFactory()
		if err != nil {
			return err
		}

		statusCode, err := m.executeJSON(req, out)
		if err == nil {
			return nil
		}

		if shouldRetryMarketplaceStatus(statusCode) && attempt < maxMarketplaceAttempts {
			slog.Warn("marketplace request retry scheduled", "path", req.URL.Path, "attempt", attempt, "status_code", statusCode, "error", err.Error())
			if sleepErr := sleepWithBackoff(ctx, attempt); sleepErr != nil {
				return sleepErr
			}
			continue
		}

		return err
	}

	return fmt.Errorf("marketplace request failed after retries")
}

func (m *MarketplaceRepository) executeJSON(req *http.Request, out any) (int, error) {
	res, err := m.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res.StatusCode, fmt.Errorf("marketplace request failed: %s", string(body))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return res.StatusCode, err
	}

	return res.StatusCode, nil
}

func (m *MarketplaceRepository) recoverAccessToken(ctx context.Context) (dto.MarketplaceTokenData, error) {
	storedToken, err := m.tokenRepo.GetByShopID(ctx, m.conf.Marketplace.ShopID)
	if err != nil {
		return dto.MarketplaceTokenData{}, err
	}

	if storedToken.RefreshToken.Valid && storedToken.RefreshToken.String != "" {
		tokenData, err := m.RefreshToken(ctx, storedToken.RefreshToken.String)
		if err == nil {
			return tokenData, nil
		}

		slog.Warn("marketplace token recovery via refresh failed, bootstrapping new token", "shop_id", m.conf.Marketplace.ShopID, "error", err.Error())
	}

	return m.Bootstrap(ctx, "pm")
}

func shouldRetryMarketplaceStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func sleepWithBackoff(ctx context.Context, attempt int) error {
	backoff := initialBackoff * time.Duration(1<<(attempt-1))
	timer := time.NewTimer(backoff)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func signHMACSHA256(key, value string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

func sqlNullTimeFromNow(expiresIn int) sql.NullTime {
	if expiresIn <= 0 {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  time.Now().Add(time.Duration(expiresIn) * time.Second),
		Valid: true,
	}
}

func sqlNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func sqlNullTime(value string) sql.NullTime {
	if value == "" {
		return sql.NullTime{}
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  parsed,
		Valid: true,
	}
}
