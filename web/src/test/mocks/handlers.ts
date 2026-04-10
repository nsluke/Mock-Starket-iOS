import { http, HttpResponse } from 'msw';

const API_URL = 'http://localhost';

// Sample data
export const mockStocks = [
  {
    ticker: 'PIPE',
    name: 'Piper Industries',
    sector: 'Technology',
    asset_type: 'stock',
    base_price: '150.00',
    current_price: '155.50',
    day_open: '150.00',
    day_high: '158.00',
    day_low: '149.00',
    prev_close: '149.50',
    volume: 1250000,
    volatility: '0.02',
    description: 'Leading tech company',
  },
  {
    ticker: 'BREW',
    name: 'BrewCraft Holdings',
    sector: 'Consumer Goods',
    asset_type: 'stock',
    base_price: '45.00',
    current_price: '43.20',
    day_open: '45.00',
    day_high: '45.50',
    day_low: '42.80',
    prev_close: '44.90',
    volume: 800000,
    volatility: '0.03',
  },
  {
    ticker: 'GLDX',
    name: 'Gold Standard ETF',
    sector: 'Commodities',
    asset_type: 'etf',
    base_price: '200.00',
    current_price: '205.00',
    day_open: '200.00',
    day_high: '206.00',
    day_low: '199.00',
    prev_close: '199.50',
    volume: 500000,
    volatility: '0.01',
  },
];

export const mockMarketSummary = {
  index_value: '12500.00',
  index_change_pct: '1.25',
  total_stocks: 3,
  gainers: 2,
  losers: 1,
};

export const mockPortfolio = {
  portfolio: {
    id: 'port-1',
    user_id: 'user-1',
    cash: '85000.00',
    net_worth: '105000.00',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-04-01T00:00:00Z',
  },
  positions: [
    {
      id: 'pos-1',
      portfolio_id: 'port-1',
      ticker: 'PIPE',
      shares: 100,
      avg_cost: '150.00',
      current_price: '155.50',
      market_value: '15550.00',
      pnl: '550.00',
      pnl_pct: '3.67',
    },
  ],
  net_worth: '105000.00',
  invested: '15000.00',
};

export const mockTrades = [
  {
    id: 'trade-1',
    user_id: 'user-1',
    ticker: 'PIPE',
    side: 'buy',
    shares: 100,
    price: '150.00',
    total: '15000.00',
    created_at: '2026-03-15T10:30:00Z',
  },
];

export const mockOrders = [
  {
    id: 'order-1',
    ticker: 'PIPE',
    side: 'buy',
    order_type: 'limit',
    shares: 50,
    limit_price: '145.00',
    stop_price: null,
    status: 'open',
    created_at: '2026-04-01T08:00:00Z',
  },
];

export const mockAlerts = [
  {
    id: 'alert-1',
    ticker: 'PIPE',
    condition: 'above',
    target_price: '160.00',
    triggered: false,
    triggered_at: null,
    created_at: '2026-04-01T09:00:00Z',
  },
  {
    id: 'alert-2',
    ticker: 'BREW',
    condition: 'below',
    target_price: '40.00',
    triggered: true,
    triggered_at: '2026-04-02T14:00:00Z',
    created_at: '2026-04-01T09:00:00Z',
  },
];

export const mockLeaderboard = [
  { id: 1, user_id: 'user-1', display_name: 'TraderJoe', net_worth: '120000.00', total_return: '20.00', rank: 1, period: 'alltime' },
  { id: 2, user_id: 'user-2', display_name: 'StockPro', net_worth: '110000.00', total_return: '10.00', rank: 2, period: 'alltime' },
];

export const mockUser = {
  id: 'user-1',
  firebase_uid: 'guest_abc123',
  display_name: 'TestTrader',
  avatar_url: null,
  is_guest: true,
  created_at: '2026-01-01T00:00:00Z',
  login_streak: 5,
  longest_streak: 12,
};

export const handlers = [
  // Auth
  http.post(`${API_URL}/api/v1/auth/guest`, () =>
    HttpResponse.json(mockUser)
  ),
  http.get(`${API_URL}/api/v1/auth/me`, () =>
    HttpResponse.json(mockUser)
  ),
  http.put(`${API_URL}/api/v1/auth/me`, () =>
    HttpResponse.json({ status: 'updated' })
  ),

  // Stocks
  http.get(`${API_URL}/api/v1/stocks`, () =>
    HttpResponse.json(mockStocks)
  ),
  http.get(`${API_URL}/api/v1/stocks/market-summary`, () =>
    HttpResponse.json(mockMarketSummary)
  ),
  http.get(`${API_URL}/api/v1/stocks/:ticker`, ({ params }) => {
    const stock = mockStocks.find((s) => s.ticker === params.ticker);
    if (!stock) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(stock);
  }),
  http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
    HttpResponse.json([])
  ),

  // Portfolio
  http.get(`${API_URL}/api/v1/portfolio`, () =>
    HttpResponse.json(mockPortfolio)
  ),
  http.get(`${API_URL}/api/v1/portfolio/history`, () =>
    HttpResponse.json([])
  ),

  // Trades
  http.post(`${API_URL}/api/v1/trades`, async ({ request }) => {
    const body = (await request.json()) as { ticker: string; side: string; shares: number };
    return HttpResponse.json({
      id: 'trade-new',
      user_id: 'user-1',
      ticker: body.ticker,
      side: body.side,
      shares: body.shares,
      price: '155.50',
      total: String(body.shares * 155.5),
      created_at: new Date().toISOString(),
    });
  }),
  http.get(`${API_URL}/api/v1/trades`, () =>
    HttpResponse.json(mockTrades)
  ),

  // Orders
  http.get(`${API_URL}/api/v1/orders`, () =>
    HttpResponse.json(mockOrders)
  ),
  http.post(`${API_URL}/api/v1/orders`, () =>
    HttpResponse.json({ id: 'order-new', status: 'open' })
  ),
  http.delete(`${API_URL}/api/v1/orders/:id`, () =>
    HttpResponse.json({ status: 'cancelled' })
  ),

  // Alerts
  http.get(`${API_URL}/api/v1/alerts`, () =>
    HttpResponse.json(mockAlerts)
  ),
  http.post(`${API_URL}/api/v1/alerts`, () =>
    HttpResponse.json({ id: 'alert-new', status: 'created' })
  ),
  http.delete(`${API_URL}/api/v1/alerts/:id`, () =>
    HttpResponse.json({ status: 'deleted' })
  ),

  // Leaderboard
  http.get(`${API_URL}/api/v1/leaderboard`, () =>
    HttpResponse.json(mockLeaderboard)
  ),

  // Achievements
  http.get(`${API_URL}/api/v1/achievements`, () =>
    HttpResponse.json([
      { id: 'ach-1', name: 'First Trade', description: 'Execute your first trade', icon: '🎯', category: 'trading' },
    ])
  ),
  http.get(`${API_URL}/api/v1/achievements/me`, () =>
    HttpResponse.json([{ id: 'ua-1', achievement_id: 'ach-1', earned_at: '2026-03-01T00:00:00Z' }])
  ),

  // Challenges
  http.get(`${API_URL}/api/v1/challenges/today`, () =>
    HttpResponse.json({
      challenge: {
        id: 'ch-1',
        date: '2026-04-06',
        challenge_type: 'trade_count',
        description: 'Execute 3 trades today',
        reward_cash: '500.00',
      },
      progress: { completed: false, completed_at: null, claimed: false },
    })
  ),
  http.post(`${API_URL}/api/v1/challenges/check`, () =>
    HttpResponse.json({ completed: false })
  ),

  // Market Status
  http.get(`${API_URL}/api/v1/market/status`, () =>
    HttpResponse.json({
      is_open: true,
      session: 'regular',
      next_open: '2026-04-09T13:30:00Z',
      next_close: '2026-04-08T20:00:00Z',
    })
  ),

  // Options
  http.get(`${API_URL}/api/v1/stocks/:ticker/options/expirations`, () =>
    HttpResponse.json([])
  ),
  http.get(`${API_URL}/api/v1/stocks/:ticker/options`, () =>
    HttpResponse.json({ calls: [], puts: [] })
  ),
  http.get(`${API_URL}/api/v1/options/positions`, () =>
    HttpResponse.json([])
  ),
  http.get(`${API_URL}/api/v1/options/trades`, () =>
    HttpResponse.json([])
  ),
  http.get(`${API_URL}/api/v1/options/orders`, () =>
    HttpResponse.json([])
  ),

  // Watchlist
  http.get(`${API_URL}/api/v1/watchlist`, () =>
    HttpResponse.json(['PIPE', 'BREW'])
  ),
  http.post(`${API_URL}/api/v1/watchlist`, () =>
    HttpResponse.json({ status: 'added' })
  ),
  http.delete(`${API_URL}/api/v1/watchlist/:ticker`, () =>
    HttpResponse.json({ status: 'removed' })
  ),
];
