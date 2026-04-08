package simulation

import (
	"fmt"
	"math"
	"time"
)

const (
	// RiskFreeRate is the simulated risk-free interest rate.
	RiskFreeRate = 0.05
	// ContractMultiplier is the number of shares per contract (standard US options).
	ContractMultiplier = 100
)

// Greeks holds the option greeks for a single contract.
type Greeks struct {
	Delta float64
	Gamma float64
	Theta float64 // per day
	Vega  float64 // per 1% vol change
	Rho   float64
}

// BlackScholes computes the theoretical option price.
// s = underlying price, k = strike, t = time to expiry (years),
// r = risk-free rate, sigma = implied volatility, isCall = true for calls.
func BlackScholes(s, k, t, r, sigma float64, isCall bool) float64 {
	if t <= 0 {
		// At or past expiration: intrinsic value only.
		if isCall {
			return math.Max(s-k, 0)
		}
		return math.Max(k-s, 0)
	}
	if sigma <= 0 {
		sigma = 0.001 // prevent division by zero
	}

	d1 := (math.Log(s/k) + (r+0.5*sigma*sigma)*t) / (sigma * math.Sqrt(t))
	d2 := d1 - sigma*math.Sqrt(t)

	if isCall {
		return s*cnd(d1) - k*math.Exp(-r*t)*cnd(d2)
	}
	return k*math.Exp(-r*t)*cnd(-d2) - s*cnd(-d1)
}

// CalculateGreeks computes all greeks for the given parameters.
func CalculateGreeks(s, k, t, r, sigma float64, isCall bool) Greeks {
	if t <= 0 || sigma <= 0 {
		d := 0.0
		if isCall {
			if s > k {
				d = 1.0
			}
		} else {
			if s < k {
				d = -1.0
			}
		}
		return Greeks{Delta: d}
	}

	sqrtT := math.Sqrt(t)
	d1 := (math.Log(s/k) + (r+0.5*sigma*sigma)*t) / (sigma * sqrtT)
	d2 := d1 - sigma*sqrtT

	nd1 := normalPDF(d1)
	discount := math.Exp(-r * t)

	var delta float64
	if isCall {
		delta = cnd(d1)
	} else {
		delta = cnd(d1) - 1
	}

	gamma := nd1 / (s * sigma * sqrtT)

	// Theta: per calendar day (divide annual theta by 365)
	var theta float64
	commonTheta := -(s * nd1 * sigma) / (2 * sqrtT)
	if isCall {
		theta = commonTheta - r*k*discount*cnd(d2)
	} else {
		theta = commonTheta + r*k*discount*cnd(-d2)
	}
	theta /= 365.0

	// Vega: per 1% change in volatility
	vega := s * sqrtT * nd1 / 100.0

	// Rho: per 1% change in interest rate
	var rho float64
	if isCall {
		rho = k * t * discount * cnd(d2) / 100.0
	} else {
		rho = -k * t * discount * cnd(-d2) / 100.0
	}

	return Greeks{
		Delta: delta,
		Gamma: gamma,
		Theta: theta,
		Vega:  vega,
		Rho:   rho,
	}
}

// SimulatedIV derives implied volatility from a stock's base volatility
// with a skew model: OTM puts get higher IV, deep OTM gets more.
func SimulatedIV(baseVolatility float64, s, k float64, isCall bool, t float64) float64 {
	// Convert per-tick volatility to annualized (approximate).
	// Base volatility is per-tick as fraction of price (e.g., 0.0001 = 0.01%/tick).
	// With 150 ticks/day, 252 trading days: annual vol ≈ basVol * sqrt(150*252)
	annualVol := baseVolatility * math.Sqrt(150*252)

	// Clamp to reasonable range
	if annualVol < 0.10 {
		annualVol = 0.10
	}
	if annualVol > 2.0 {
		annualVol = 2.0
	}

	// Moneyness ratio
	moneyness := k / s
	otmFactor := 0.0

	if isCall {
		// OTM calls (strike > spot): slight smile
		if moneyness > 1.0 {
			otmFactor = (moneyness - 1.0) * 0.3
		}
	} else {
		// OTM puts (strike < spot): steeper skew (volatility smile)
		if moneyness < 1.0 {
			otmFactor = (1.0 - moneyness) * 0.5
		}
	}

	iv := annualVol * (1.0 + otmFactor)

	// Time decay on skew: less skew for longer-dated options
	if t > 0.25 {
		iv -= otmFactor * 0.2
	}

	if iv < 0.05 {
		iv = 0.05
	}
	return iv
}

// GenerateStrikes returns strike prices centered around the current price.
// Increment: <$25 → $1, $25-$200 → $5, >$200 → $10.
// Returns ~10-15 strikes above and below current price.
func GenerateStrikes(currentPrice float64) []float64 {
	var increment float64
	switch {
	case currentPrice < 25:
		increment = 1.0
	case currentPrice < 200:
		increment = 5.0
	default:
		increment = 10.0
	}

	// Round current price to nearest increment
	center := math.Round(currentPrice/increment) * increment

	strikes := make([]float64, 0, 25)
	for i := -12; i <= 12; i++ {
		strike := center + float64(i)*increment
		if strike > 0 {
			strikes = append(strikes, math.Round(strike*100)/100)
		}
	}
	return strikes
}

// GenerateExpirations returns simulated expiration dates from the given time.
// Weekly (next 4), monthly (next 3), quarterly (next 2).
func GenerateExpirations(now time.Time) []time.Time {
	expirations := make([]time.Time, 0, 9)

	// Weekly: next 4 Fridays
	next := nextFriday(now)
	for i := 0; i < 4; i++ {
		expirations = append(expirations, next)
		next = next.AddDate(0, 0, 7)
	}

	// Monthly: next 3 months (3rd Friday of month)
	for i := 1; i <= 3; i++ {
		monthStart := time.Date(now.Year(), now.Month()+time.Month(i), 1, 16, 0, 0, 0, time.UTC)
		expirations = append(expirations, thirdFriday(monthStart))
	}

	// Quarterly: next 2 quarter ends
	for i := 1; i <= 2; i++ {
		q := quarterEnd(now, i)
		expirations = append(expirations, thirdFriday(q))
	}

	// Deduplicate and sort
	seen := map[string]bool{}
	unique := make([]time.Time, 0, len(expirations))
	for _, exp := range expirations {
		key := exp.Format("2006-01-02")
		if !seen[key] {
			seen[key] = true
			unique = append(unique, exp)
		}
	}
	return unique
}

// BuildContractSymbol creates an OCC-style symbol: TICKER + YYMMDD + C/P + strike*1000 (8 digits).
func BuildContractSymbol(ticker string, expiration time.Time, optionType string, strike float64) string {
	typeChar := "C"
	if optionType == "put" {
		typeChar = "P"
	}
	return fmt.Sprintf("%s%s%s%08d", ticker, expiration.Format("060102"), typeChar, int(strike*1000))
}

// --- helpers ---

// cnd computes the cumulative standard normal distribution.
func cnd(x float64) float64 {
	return 0.5 * math.Erfc(-x/math.Sqrt2)
}

// normalPDF returns the standard normal probability density function.
func normalPDF(x float64) float64 {
	return math.Exp(-0.5*x*x) / math.Sqrt(2*math.Pi)
}

func nextFriday(t time.Time) time.Time {
	daysUntil := (5 - int(t.Weekday()) + 7) % 7
	if daysUntil == 0 {
		daysUntil = 7
	}
	friday := t.AddDate(0, 0, daysUntil)
	return time.Date(friday.Year(), friday.Month(), friday.Day(), 16, 0, 0, 0, time.UTC)
}

func thirdFriday(monthStart time.Time) time.Time {
	first := time.Date(monthStart.Year(), monthStart.Month(), 1, 16, 0, 0, 0, time.UTC)
	// First Friday
	daysUntil := (5 - int(first.Weekday()) + 7) % 7
	firstFriday := first.AddDate(0, 0, daysUntil)
	// Third Friday = first Friday + 14 days
	return firstFriday.AddDate(0, 0, 14)
}

func quarterEnd(now time.Time, offset int) time.Time {
	// Quarter months: March, June, September, December
	currentQ := (int(now.Month()) - 1) / 3
	targetQ := currentQ + offset
	targetMonth := time.Month((targetQ%4)*3 + 3)
	targetYear := now.Year() + targetQ/4
	return time.Date(targetYear, targetMonth, 1, 16, 0, 0, 0, time.UTC)
}
