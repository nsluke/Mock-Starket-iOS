package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luke/mockstarket/internal/middleware"
	"github.com/luke/mockstarket/internal/service"
	"github.com/shopspring/decimal"
)

// ---- Options Handlers ----

// GetOptionChain returns the option chain for a stock.
func (h *Handler) GetOptionChain(w http.ResponseWriter, r *http.Request) {
	ticker := tickerParam(r)

	var expiration *time.Time
	if expStr := r.URL.Query().Get("expiration"); expStr != "" {
		t, err := time.Parse(time.RFC3339, expStr)
		if err != nil {
			// Try date-only format
			t, err = time.Parse("2006-01-02", expStr)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid expiration format")
				return
			}
			t = time.Date(t.Year(), t.Month(), t.Day(), 16, 0, 0, 0, time.UTC)
		}
		expiration = &t
	}

	contracts, err := h.repo.GetOptionChain(r.Context(), ticker, expiration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load option chain")
		return
	}

	// Split into calls and puts
	type chainResponse struct {
		Ticker          string      `json:"ticker"`
		UnderlyingPrice string      `json:"underlying_price"`
		Calls           interface{} `json:"calls"`
		Puts            interface{} `json:"puts"`
	}

	calls := make([]interface{}, 0)
	puts := make([]interface{}, 0)
	for _, c := range contracts {
		if c.OptionType == "call" {
			calls = append(calls, c)
		} else {
			puts = append(puts, c)
		}
	}

	price, _ := h.engine.GetPrice(ticker)
	writeJSON(w, http.StatusOK, chainResponse{
		Ticker:          ticker,
		UnderlyingPrice: price.String(),
		Calls:           calls,
		Puts:            puts,
	})
}

// GetOptionExpirations returns available expiration dates for a stock's options.
func (h *Handler) GetOptionExpirations(w http.ResponseWriter, r *http.Request) {
	ticker := tickerParam(r)

	expirations, err := h.repo.GetOptionExpirations(r.Context(), ticker)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load expirations")
		return
	}

	writeJSON(w, http.StatusOK, expirations)
}

// GetOptionContractDetail returns a single option contract.
func (h *Handler) GetOptionContractDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid contract ID")
		return
	}

	contract, err := h.repo.GetOptionContract(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "contract not found")
		return
	}

	writeJSON(w, http.StatusOK, contract)
}

// ExecuteOptionsTrade handles an options trade execution.
func (h *Handler) ExecuteOptionsTrade(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	user, err := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req service.OptionsTradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	trade, err := h.optionsTradeSvc.ExecuteOptionsTrade(r.Context(), user.ID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, trade)
}

// GetOptionsTradeHistory returns paginated options trade history.
func (h *Handler) GetOptionsTradeHistory(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	user, err := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	trades, err := h.repo.GetOptionTradesByUserID(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load trades")
		return
	}

	writeJSON(w, http.StatusOK, trades)
}

// GetOptionsPositions returns the user's options positions with enriched data.
func (h *Handler) GetOptionsPositions(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	user, err := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	portfolio, err := h.repo.GetPortfolioByUserID(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "portfolio not found")
		return
	}

	positions, err := h.repo.GetOptionPositionsByPortfolio(r.Context(), portfolio.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load positions")
		return
	}

	type enrichedPosition struct {
		ID          uuid.UUID       `json:"id"`
		PortfolioID uuid.UUID       `json:"portfolio_id"`
		ContractID  uuid.UUID       `json:"contract_id"`
		Quantity    int             `json:"quantity"`
		AvgCost     decimal.Decimal `json:"avg_cost"`
		Collateral  decimal.Decimal `json:"collateral"`
		Contract    interface{}     `json:"contract"`
		MarketValue decimal.Decimal `json:"market_value"`
		PnL         decimal.Decimal `json:"pnl"`
		PnLPct      decimal.Decimal `json:"pnl_pct"`
		IsLong      bool            `json:"is_long"`
	}

	multiplier := decimal.NewFromInt(100)
	result := make([]enrichedPosition, 0, len(positions))
	for _, pos := range positions {
		contract, err := h.repo.GetOptionContract(r.Context(), pos.ContractID)
		if err != nil {
			continue
		}

		absQty := decimal.NewFromInt(int64(pos.Quantity)).Abs()
		marketValue := contract.MarkPrice.Mul(absQty).Mul(multiplier)

		var pnl decimal.Decimal
		if pos.Quantity > 0 {
			// Long: P&L = (mark - avg_cost) * qty * 100
			pnl = contract.MarkPrice.Sub(pos.AvgCost).Mul(absQty).Mul(multiplier)
		} else {
			// Short: P&L = (avg_cost - mark) * qty * 100
			pnl = pos.AvgCost.Sub(contract.MarkPrice).Mul(absQty).Mul(multiplier)
		}

		totalCost := pos.AvgCost.Mul(absQty).Mul(multiplier)
		pnlPct := decimal.Zero
		if !totalCost.IsZero() {
			pnlPct = pnl.Div(totalCost).Mul(decimal.NewFromInt(100))
		}

		result = append(result, enrichedPosition{
			ID:          pos.ID,
			PortfolioID: pos.PortfolioID,
			ContractID:  pos.ContractID,
			Quantity:    pos.Quantity,
			AvgCost:     pos.AvgCost,
			Collateral:  pos.Collateral,
			Contract:    contract,
			MarketValue: marketValue,
			PnL:         pnl,
			PnLPct:      pnlPct,
			IsLong:      pos.Quantity > 0,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

// CreateOptionsOrder creates a limit order for options.
func (h *Handler) CreateOptionsOrder(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	user, err := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req struct {
		ContractID uuid.UUID        `json:"contract_id"`
		Side       string           `json:"side"`
		OrderType  string           `json:"order_type"`
		Quantity   int              `json:"quantity"`
		LimitPrice *decimal.Decimal `json:"limit_price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "quantity must be positive")
		return
	}
	if req.OrderType == "limit" && req.LimitPrice == nil {
		writeError(w, http.StatusBadRequest, "limit_price required for limit orders")
		return
	}

	order, err := h.repo.CreateOptionOrder(r.Context(), user.ID, req.ContractID, req.Side, req.OrderType, req.Quantity, req.LimitPrice)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

// ListOptionsOrders returns the user's open options orders.
func (h *Handler) ListOptionsOrders(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	user, err := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	orders, err := h.repo.GetOpenOptionOrders(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load orders")
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

// CancelOptionsOrder cancels an open options order.
func (h *Handler) CancelOptionsOrder(w http.ResponseWriter, r *http.Request) {
	firebaseUID := middleware.GetFirebaseUID(r.Context())
	user, err := h.repo.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	if err := h.repo.CancelOptionOrder(r.Context(), id, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to cancel order")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}
