package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/market"
	"github.com/luke/mockstarket/internal/middleware"
	"github.com/luke/mockstarket/internal/model"
	"github.com/luke/mockstarket/internal/polygon"
	"github.com/luke/mockstarket/internal/repository"
	"github.com/luke/mockstarket/internal/service"
	ws "github.com/luke/mockstarket/internal/websocket"
	"github.com/shopspring/decimal"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	repo            *repository.Repo
	tradeSvc        *service.TradeService
	optionsTradeSvc *service.OptionsTradeService
	challengeSvc    *service.ChallengeService
	engine          market.PriceProvider
	hub             *ws.Hub
	startingCash    float64
}

// New creates a new Handler.
func New(repo *repository.Repo, tradeSvc *service.TradeService, engine market.PriceProvider, hub *ws.Hub, startingCash float64) *Handler {
	return &Handler{
		repo:         repo,
		tradeSvc:     tradeSvc,
		engine:       engine,
		hub:          hub,
		startingCash: startingCash,
	}
}

// SetChallengeService sets the challenge service (avoids circular deps during init).
func (h *Handler) SetChallengeService(svc *service.ChallengeService) {
	h.challengeSvc = svc
}

// SetOptionsTradeService sets the options trade service.
func (h *Handler) SetOptionsTradeService(svc *service.OptionsTradeService) {
	h.optionsTradeSvc = svc
}

// ---- Market Status ----

func (h *Handler) GetMarketStatus(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	session := polygon.GetMarketSession(now)
	isOpen := session == polygon.SessionRegular

	resp := map[string]interface{}{
		"is_open":    isOpen,
		"session":    string(session),
		"next_open":  polygon.NextMarketOpen(now).Format(time.RFC3339),
		"next_close": polygon.NextMarketClose(now).Format(time.RFC3339),
		"timestamp":  now.Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, resp)
}

// ---- Auth ----

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())

	var req struct {
		DisplayName string `json:"display_name"`
		IsGuest     bool   `json:"is_guest"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DisplayName == "" {
		req.DisplayName = "Trader"
	}

	user, err := h.repo.CreateUser(r.Context(), firebaseUID, req.DisplayName, req.IsGuest)
	if err != nil {
		// User already exists — return existing user (idempotent register)
		existing, getErr := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
		if getErr != nil {
			writeError(w, http.StatusConflict, "user already exists")
			return
		}
		writeJSON(w, http.StatusOK, existing)
		return
	}

	// Create portfolio with starting cash
	startingCash := decimal.NewFromFloat(h.startingCash)
	_, err = h.repo.CreatePortfolio(r.Context(), user.ID, startingCash)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create portfolio")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) CreateGuest(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())

	user, err := h.repo.CreateUser(r.Context(), firebaseUID, "Guest Trader", true)
	if err != nil {
		// Guest already exists — return existing user
		existing, getErr := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
		if getErr != nil {
			writeError(w, http.StatusConflict, "user already exists")
			return
		}
		writeJSON(w, http.StatusOK, existing)
		return
	}

	startingCash := decimal.NewFromFloat(h.startingCash)
	_, err = h.repo.CreatePortfolio(r.Context(), user.ID, startingCash)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create portfolio")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	// Update login streak on each auth check
	_ = h.repo.UpdateLoginStreak(r.Context(), user.ID)

	// Re-fetch to get updated streak values
	user, _ = h.repo.GetUserByID(r.Context(), user.ID)

	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	var req struct {
		DisplayName string  `json:"display_name"`
		AvatarURL   *string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DisplayName == "" {
		req.DisplayName = user.DisplayName
	}

	if err := h.repo.UpdateUser(r.Context(), user.ID, req.DisplayName, req.AvatarURL); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	if err := h.repo.DeleteUser(r.Context(), user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- Stocks ----

func (h *Handler) ListStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.repo.GetAllStocks(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch stocks")
		return
	}

	// Overlay live prices from simulation engine
	livePrices := h.engine.GetAllPrices()
	for i := range stocks {
		if price, ok := livePrices[stocks[i].Ticker]; ok {
			stocks[i].CurrentPrice = price
		}
	}

	writeJSON(w, http.StatusOK, stocks)
}

func (h *Handler) GetStock(w http.ResponseWriter, r *http.Request) {
	ticker := tickerParam(r)
	stock, err := h.repo.GetStockByTicker(r.Context(), ticker)
	if err != nil {
		writeError(w, http.StatusNotFound, "stock not found")
		return
	}

	if price, ok := h.engine.GetPrice(ticker); ok {
		stock.CurrentPrice = price
	}

	writeJSON(w, http.StatusOK, stock)
}

func (h *Handler) GetStockHistory(w http.ResponseWriter, r *http.Request) {
	ticker := tickerParam(r)
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1m"
	}

	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	history, err := h.repo.GetPriceHistory(r.Context(), ticker, interval, from, to, 500)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch history")
		return
	}

	writeJSON(w, http.StatusOK, history)
}

func (h *Handler) GetMarketSummary(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.repo.GetAllStocks(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch stocks")
		return
	}

	livePrices := h.engine.GetAllPrices()

	totalValue := decimal.Zero
	totalPrevClose := decimal.Zero
	var gainers, losers []model.Stock

	for i := range stocks {
		if price, ok := livePrices[stocks[i].Ticker]; ok {
			stocks[i].CurrentPrice = price
		}
		totalValue = totalValue.Add(stocks[i].CurrentPrice)
		totalPrevClose = totalPrevClose.Add(stocks[i].PrevClose)

		change := stocks[i].CurrentPrice.Sub(stocks[i].PrevClose)
		if change.IsPositive() {
			gainers = append(gainers, stocks[i])
		} else if change.IsNegative() {
			losers = append(losers, stocks[i])
		}
	}

	indexChange := decimal.Zero
	if !totalPrevClose.IsZero() {
		indexChange = totalValue.Sub(totalPrevClose).Div(totalPrevClose).Mul(decimal.NewFromInt(100)).Round(2)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"index_value":      totalValue.Div(decimal.NewFromInt(int64(len(stocks)))).Round(2),
		"index_change_pct": indexChange,
		"total_stocks":     len(stocks),
		"gainers":          len(gainers),
		"losers":           len(losers),
	})
}

func (h *Handler) GetETFHoldings(w http.ResponseWriter, r *http.Request) {
	ticker := tickerParam(r)

	holdings, err := h.repo.GetETFHoldings(r.Context(), ticker)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch ETF holdings")
		return
	}

	// Enrich with current prices and stock names
	livePrices := h.engine.GetAllPrices()
	type enrichedHolding struct {
		Ticker string  `json:"ticker"`
		Name   string  `json:"name"`
		Weight string  `json:"weight"`
		Price  string  `json:"price"`
	}

	result := make([]enrichedHolding, 0, len(holdings))
	for _, holding := range holdings {
		price := ""
		if p, ok := livePrices[holding.HoldingTicker]; ok {
			price = p.String()
		}
		name := holding.HoldingTicker
		if stock, err := h.repo.GetStockByTicker(r.Context(), holding.HoldingTicker); err == nil {
			name = stock.Name
		}
		result = append(result, enrichedHolding{
			Ticker: holding.HoldingTicker,
			Name:   name,
			Weight: holding.Weight.String(),
			Price:  price,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

// ---- Trading ----

func (h *Handler) ExecuteTrade(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req service.TradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	trade, err := h.tradeSvc.ExecuteTrade(r.Context(), user.ID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Notify via WebSocket
	tradeData, _ := json.Marshal(trade)
	h.hub.SendToUser(user.ID.String(), ws.Message{
		Type: "trade_executed",
		Data: tradeData,
	})

	writeJSON(w, http.StatusCreated, trade)
}

func (h *Handler) GetTradeHistory(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	limit := getQueryInt(r, "limit", 50)
	offset := getQueryInt(r, "offset", 0)

	trades, err := h.tradeSvc.GetTradeHistory(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch trades")
		return
	}

	writeJSON(w, http.StatusOK, trades)
}

// ---- Orders ----

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req struct {
		Ticker     string           `json:"ticker"`
		Side       string           `json:"side"`
		OrderType  string           `json:"order_type"`
		Shares     int              `json:"shares"`
		LimitPrice *decimal.Decimal `json:"limit_price,omitempty"`
		StopPrice  *decimal.Decimal `json:"stop_price,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Shares <= 0 {
		writeError(w, http.StatusBadRequest, "shares must be positive")
		return
	}
	if req.Side != "buy" && req.Side != "sell" {
		writeError(w, http.StatusBadRequest, "side must be 'buy' or 'sell'")
		return
	}

	// Validate order type and required price fields
	switch req.OrderType {
	case "limit":
		if req.LimitPrice == nil {
			writeError(w, http.StatusBadRequest, "limit_price is required for limit orders")
			return
		}
	case "stop":
		if req.StopPrice == nil {
			writeError(w, http.StatusBadRequest, "stop_price is required for stop orders")
			return
		}
	case "stop_limit":
		if req.LimitPrice == nil || req.StopPrice == nil {
			writeError(w, http.StatusBadRequest, "limit_price and stop_price are required for stop_limit orders")
			return
		}
	default:
		writeError(w, http.StatusBadRequest, "order_type must be 'limit', 'stop', or 'stop_limit'")
		return
	}

	// Validate ticker exists (check price provider, then DB)
	if _, ok := h.engine.GetPrice(req.Ticker); !ok {
		if _, err := h.repo.GetStockByTicker(r.Context(), req.Ticker); err != nil {
			writeError(w, http.StatusBadRequest, "stock "+req.Ticker+" not found")
			return
		}
	}

	order, err := h.repo.CreateOrder(r.Context(), user.ID, req.Ticker, req.Side, req.OrderType, req.Shares, req.LimitPrice, req.StopPrice)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	orders, err := h.repo.GetOpenOrders(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	if err := h.repo.CancelOrder(r.Context(), orderID, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to cancel order")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

// ---- Portfolio ----

func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	portfolio, err := h.repo.GetPortfolioByUserID(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "portfolio not found")
		return
	}

	holdings, err := h.repo.GetHoldingsByPortfolioID(r.Context(), portfolio.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch holdings")
		return
	}

	// Calculate net worth with live prices
	livePrices := h.engine.GetAllPrices()
	investedValue := decimal.Zero
	type positionResponse struct {
		model.Holding
		CurrentPrice decimal.Decimal `json:"current_price"`
		MarketValue  decimal.Decimal `json:"market_value"`
		PnL          decimal.Decimal `json:"pnl"`
		PnLPct       decimal.Decimal `json:"pnl_pct"`
	}

	positions := make([]positionResponse, 0, len(holdings))
	for _, h := range holdings {
		currentPrice := h.AvgCost
		if price, ok := livePrices[h.Ticker]; ok {
			currentPrice = price
		}
		marketValue := currentPrice.Mul(decimal.NewFromInt(int64(h.Shares)))
		costBasis := h.AvgCost.Mul(decimal.NewFromInt(int64(h.Shares)))
		pnl := marketValue.Sub(costBasis)
		pnlPct := decimal.Zero
		if !costBasis.IsZero() {
			pnlPct = pnl.Div(costBasis).Mul(decimal.NewFromInt(100)).Round(2)
		}
		investedValue = investedValue.Add(marketValue)

		positions = append(positions, positionResponse{
			Holding:      h,
			CurrentPrice: currentPrice,
			MarketValue:  marketValue,
			PnL:          pnl,
			PnLPct:       pnlPct,
		})
	}

	netWorth := portfolio.Cash.Add(investedValue)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"portfolio": portfolio,
		"positions": positions,
		"net_worth": netWorth,
		"invested":  investedValue,
	})
}

func (h *Handler) GetPortfolioHistory(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	limit := getQueryInt(r, "limit", 100)
	history, err := h.repo.GetPortfolioHistory(r.Context(), user.ID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch history")
		return
	}

	writeJSON(w, http.StatusOK, history)
}

// ---- Leaderboard ----

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "alltime"
	}
	limit := getQueryInt(r, "limit", 50)
	offset := getQueryInt(r, "offset", 0)

	entries, err := h.repo.GetLeaderboard(r.Context(), period, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch leaderboard")
		return
	}

	writeJSON(w, http.StatusOK, entries)
}

// ---- Alerts ----

func (h *Handler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req struct {
		Ticker      string          `json:"ticker"`
		Condition   string          `json:"condition"`
		TargetPrice decimal.Decimal `json:"target_price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	alert, err := h.repo.CreatePriceAlert(r.Context(), user.ID, req.Ticker, req.Condition, req.TargetPrice)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, alert)
}

func (h *Handler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	alerts, err := h.repo.GetAlertsByUserID(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch alerts")
		return
	}

	writeJSON(w, http.StatusOK, alerts)
}

func (h *Handler) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	alertID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid alert ID")
		return
	}

	if err := h.repo.DeleteAlert(r.Context(), alertID, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete alert")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- Achievements ----

func (h *Handler) ListAchievements(w http.ResponseWriter, r *http.Request) {
	achievements, err := h.repo.GetAllAchievements(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch achievements")
		return
	}
	writeJSON(w, http.StatusOK, achievements)
}

func (h *Handler) GetMyAchievements(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	achievements, err := h.repo.GetUserAchievements(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch achievements")
		return
	}
	writeJSON(w, http.StatusOK, achievements)
}

// ---- Watchlist ----

func (h *Handler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	tickers, err := h.repo.GetWatchlist(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch watchlist")
		return
	}
	writeJSON(w, http.StatusOK, tickers)
}

func (h *Handler) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req struct {
		Ticker string `json:"ticker"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.repo.AddToWatchlist(r.Context(), user.ID, req.Ticker); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add to watchlist")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "added"})
}

func (h *Handler) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	ticker := tickerParam(r)
	if err := h.repo.RemoveFromWatchlist(r.Context(), user.ID, ticker); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove from watchlist")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

// ---- Challenges ----

func (h *Handler) GetTodaysChallenge(w http.ResponseWriter, r *http.Request) {
	if h.challengeSvc == nil {
		writeError(w, http.StatusServiceUnavailable, "challenges not available")
		return
	}

	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	challenge, userChallenge, err := h.challengeSvc.GetTodaysChallenge(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "no challenge today")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"challenge": challenge,
		"progress":  userChallenge,
	})
}

func (h *Handler) CheckChallenge(w http.ResponseWriter, r *http.Request) {
	if h.challengeSvc == nil {
		writeError(w, http.StatusServiceUnavailable, "challenges not available")
		return
	}

	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	completed, err := h.challengeSvc.EvaluateForUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to evaluate challenge")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"completed": completed,
	})
}

func (h *Handler) ClaimChallengeReward(w http.ResponseWriter, r *http.Request) {
	if h.challengeSvc == nil {
		writeError(w, http.StatusServiceUnavailable, "challenges not available")
		return
	}

	user, err := h.getUserFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	challengeID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid challenge ID")
		return
	}

	if err := h.challengeSvc.ClaimReward(r.Context(), user.ID, challengeID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "claimed"})
}

// ---- System ----

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":     "ok",
		"time":       time.Now().UTC(),
		"ws_clients": h.hub.ClientCount(),
	})
}

// ---- Helpers ----

func (h *Handler) getUserFromContext(r *http.Request) (*model.User, error) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	return h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
}

// tickerParam extracts and URL-decodes the {ticker} path parameter.
// Handles encoded tickers like X%3ABTCUSD -> X:BTCUSD.
func tickerParam(r *http.Request) string {
	raw := chi.URLParam(r, "ticker")
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}
	return decoded
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func getQueryInt(r *http.Request, key string, fallback int) int {
	if val := r.URL.Query().Get(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
