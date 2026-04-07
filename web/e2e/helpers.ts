import { type Page } from '@playwright/test';

// Mock data matching what the backend would return
const mockStocks = [
  {
    ticker: 'PIPE', name: 'Piper Industries', sector: 'Technology', asset_type: 'stock',
    base_price: '150.00', current_price: '155.50', day_open: '150.00', day_high: '158.00',
    day_low: '149.00', prev_close: '149.50', volume: 1250000, volatility: '0.02',
    description: 'Leading tech company',
  },
  {
    ticker: 'BREW', name: 'BrewCraft Holdings', sector: 'Consumer Goods', asset_type: 'stock',
    base_price: '45.00', current_price: '43.20', day_open: '45.00', day_high: '45.50',
    day_low: '42.80', prev_close: '44.90', volume: 800000, volatility: '0.03',
  },
];

const mockUser = {
  id: 'user-1', firebase_uid: 'guest_test123', display_name: 'E2ETrader',
  avatar_url: null, is_guest: true, created_at: '2026-01-01T00:00:00Z',
  login_streak: 3, longest_streak: 7,
};

const mockPortfolio = {
  portfolio: { id: 'port-1', user_id: 'user-1', cash: '85000.00', net_worth: '105000.00', created_at: '2026-01-01T00:00:00Z', updated_at: '2026-04-01T00:00:00Z' },
  positions: [
    { id: 'pos-1', portfolio_id: 'port-1', ticker: 'PIPE', shares: 100, avg_cost: '150.00', current_price: '155.50', market_value: '15550.00', pnl: '550.00', pnl_pct: '3.67' },
  ],
  net_worth: '105000.00',
  invested: '15000.00',
};

const mockLeaderboard = [
  { id: 1, user_id: 'user-1', display_name: 'TraderJoe', net_worth: '120000.00', total_return: '20.00', rank: 1, period: 'alltime' },
  { id: 2, user_id: 'user-2', display_name: 'StockPro', net_worth: '110000.00', total_return: '10.00', rank: 2, period: 'alltime' },
];

export async function mockAPI(page: Page) {
  // Intercept all API calls with a single handler for reliability
  await page.route('**/api/v1/**', (route) => {
    const url = route.request().url();
    const method = route.request().method();
    const path = new URL(url).pathname;

    // Auth
    if (path === '/api/v1/auth/guest') return route.fulfill({ json: mockUser });
    if (path === '/api/v1/auth/me' && method === 'GET') return route.fulfill({ json: mockUser });
    if (path === '/api/v1/auth/me' && method === 'PUT') return route.fulfill({ json: { status: 'updated' } });
    if (path === '/api/v1/auth/me' && method === 'DELETE') return route.fulfill({ json: { status: 'deleted' } });

    // Stocks
    if (path === '/api/v1/stocks/market-summary') return route.fulfill({ json: { index_value: '12500.00', index_change_pct: '1.25', total_stocks: 2, gainers: 1, losers: 1 } });
    if (path.match(/\/api\/v1\/stocks\/[A-Z]+\/history/)) return route.fulfill({ json: [] });
    if (path.match(/\/api\/v1\/stocks\/[A-Z]+\/holdings/)) return route.fulfill({ json: [] });
    if (path.match(/\/api\/v1\/stocks\/([A-Z]+)$/)) {
      const ticker = path.split('/').pop()!;
      const stock = mockStocks.find((s) => s.ticker === ticker);
      return route.fulfill({ json: stock || mockStocks[0] });
    }
    if (path === '/api/v1/stocks') return route.fulfill({ json: mockStocks });

    // Portfolio
    if (path.startsWith('/api/v1/portfolio/history')) return route.fulfill({ json: [] });
    if (path === '/api/v1/portfolio') return route.fulfill({ json: mockPortfolio });

    // Trades
    if (path === '/api/v1/trades' && method === 'POST') {
      return route.fulfill({ json: { id: 'trade-new', ticker: 'PIPE', side: 'buy', shares: 10, price: '155.50', total: '1555.00', created_at: new Date().toISOString() } });
    }
    if (path.startsWith('/api/v1/trades')) return route.fulfill({ json: [] });

    // Orders
    if (path.startsWith('/api/v1/orders') && method === 'DELETE') return route.fulfill({ json: { status: 'cancelled' } });
    if (path === '/api/v1/orders' && method === 'POST') return route.fulfill({ json: { id: 'order-new', status: 'open' } });
    if (path.startsWith('/api/v1/orders')) return route.fulfill({ json: [] });

    // Alerts
    if (path.startsWith('/api/v1/alerts') && method === 'DELETE') return route.fulfill({ json: { status: 'deleted' } });
    if (path === '/api/v1/alerts' && method === 'POST') return route.fulfill({ json: { id: 'alert-new' } });
    if (path.startsWith('/api/v1/alerts')) return route.fulfill({ json: [] });

    // Leaderboard
    if (path.startsWith('/api/v1/leaderboard')) return route.fulfill({ json: mockLeaderboard });

    // Achievements
    if (path === '/api/v1/achievements/me') return route.fulfill({ json: [] });
    if (path === '/api/v1/achievements') return route.fulfill({ json: [{ id: 'ach-1', name: 'First Trade', description: 'Execute your first trade', icon: '🎯', category: 'trading' }] });

    // Challenges
    if (path.startsWith('/api/v1/challenges')) return route.fulfill({ json: { challenge: null, progress: null } });

    // Watchlist
    if (path.startsWith('/api/v1/watchlist')) return route.fulfill({ json: [] });

    // System
    if (path === '/api/v1/system/health') return route.fulfill({ json: { status: 'ok' } });

    // Fallback
    return route.fulfill({ status: 404, json: { error: 'not found' } });
  });
}

export async function loginAsGuest(page: Page) {
  await page.goto('/login');
  await page.getByRole('button', { name: 'Continue as Guest' }).click();
  await page.waitForURL('/market');
}
