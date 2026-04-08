package polygon

import "time"

var nyLoc *time.Location

func init() {
	var err error
	nyLoc, err = time.LoadLocation("America/New_York")
	if err != nil {
		// Fallback: UTC-5 (EST without DST handling)
		nyLoc = time.FixedZone("EST", -5*60*60)
	}
}

// MarketSession describes the current trading session.
type MarketSession string

const (
	SessionPreMarket  MarketSession = "pre_market"
	SessionRegular    MarketSession = "regular"
	SessionAfterHours MarketSession = "after_hours"
	SessionClosed     MarketSession = "closed"
)

// GetMarketSession returns the current market session for US equities.
func GetMarketSession(t time.Time) MarketSession {
	ny := t.In(nyLoc)

	if ny.Weekday() == time.Saturday || ny.Weekday() == time.Sunday {
		return SessionClosed
	}

	if isUSMarketHoliday(ny) {
		return SessionClosed
	}

	hour, min := ny.Hour(), ny.Minute()
	minuteOfDay := hour*60 + min

	switch {
	case minuteOfDay >= 4*60 && minuteOfDay < 9*60+30:
		return SessionPreMarket
	case minuteOfDay >= 9*60+30 && minuteOfDay < 16*60:
		return SessionRegular
	case minuteOfDay >= 16*60 && minuteOfDay < 20*60:
		return SessionAfterHours
	default:
		return SessionClosed
	}
}

// IsMarketOpen returns true if US stock market regular session is active.
func IsMarketOpen(t time.Time) bool {
	return GetMarketSession(t) == SessionRegular
}

// NextMarketOpen returns the next regular session open time.
func NextMarketOpen(t time.Time) time.Time {
	ny := t.In(nyLoc)

	// Try today first if before open
	candidate := time.Date(ny.Year(), ny.Month(), ny.Day(), 9, 30, 0, 0, nyLoc)
	if candidate.After(t) && candidate.Weekday() != time.Saturday && candidate.Weekday() != time.Sunday && !isUSMarketHoliday(candidate) {
		return candidate
	}

	// Otherwise advance day by day
	for i := 1; i <= 7; i++ {
		next := ny.AddDate(0, 0, i)
		candidate = time.Date(next.Year(), next.Month(), next.Day(), 9, 30, 0, 0, nyLoc)
		if candidate.Weekday() != time.Saturday && candidate.Weekday() != time.Sunday && !isUSMarketHoliday(candidate) {
			return candidate
		}
	}

	// Shouldn't happen but return next Monday
	return candidate
}

// NextMarketClose returns the next regular session close time.
func NextMarketClose(t time.Time) time.Time {
	ny := t.In(nyLoc)
	candidate := time.Date(ny.Year(), ny.Month(), ny.Day(), 16, 0, 0, 0, nyLoc)
	if candidate.After(t) && IsMarketOpen(t) {
		return candidate
	}
	// Find next trading day
	open := NextMarketOpen(t)
	return time.Date(open.Year(), open.Month(), open.Day(), 16, 0, 0, 0, nyLoc)
}

// isUSMarketHoliday checks if the given date is a US stock market holiday.
// Covers 2025 and 2026 holidays.
func isUSMarketHoliday(t time.Time) bool {
	key := t.Format("01-02")
	year := t.Year()

	// Fixed-date holidays
	fixedHolidays := map[string]bool{
		"01-01": true, // New Year's Day
		"06-19": true, // Juneteenth
		"07-04": true, // Independence Day
		"12-25": true, // Christmas
	}
	if fixedHolidays[key] {
		return true
	}

	// Variable holidays by year
	type ymd struct {
		year  int
		month time.Month
		day   int
	}
	variableHolidays := []ymd{
		// 2025
		{2025, time.January, 20},  // MLK Day
		{2025, time.February, 17}, // Presidents' Day
		{2025, time.April, 18},    // Good Friday
		{2025, time.May, 26},      // Memorial Day
		{2025, time.September, 1}, // Labor Day
		{2025, time.November, 27}, // Thanksgiving
		// 2026
		{2026, time.January, 19},  // MLK Day
		{2026, time.February, 16}, // Presidents' Day
		{2026, time.April, 3},     // Good Friday
		{2026, time.May, 25},      // Memorial Day
		{2026, time.September, 7}, // Labor Day
		{2026, time.November, 26}, // Thanksgiving
	}

	for _, h := range variableHolidays {
		if year == h.year && t.Month() == h.month && t.Day() == h.day {
			return true
		}
	}

	return false
}
