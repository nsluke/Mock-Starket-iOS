const API_URL = process.env.NEXT_PUBLIC_API_URL || '';

class APIClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string = API_URL) {
    this.baseURL = baseURL;
  }

  setToken(token: string | null) {
    this.token = token;
  }

  async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${this.baseURL}${path}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`HTTP ${response.status}: ${error}`);
    }

    return response.json();
  }

  // Auth
  register(displayName: string, isGuest: boolean) {
    return this.request<any>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify({ display_name: displayName, is_guest: isGuest }),
    });
  }

  createGuest() {
    return this.request<any>('/api/v1/auth/guest', { method: 'POST' });
  }

  getMe() {
    return this.request<any>('/api/v1/auth/me');
  }

  // Stocks
  getStocks() {
    return this.request<any[]>('/api/v1/stocks');
  }

  getStock(ticker: string) {
    return this.request<any>(`/api/v1/stocks/${ticker}`);
  }

  getStockHistory(ticker: string, interval: string = '1m') {
    return this.request<any[]>(`/api/v1/stocks/${ticker}/history?interval=${interval}`);
  }

  getETFHoldings(ticker: string) {
    return this.request<any[]>(`/api/v1/stocks/${ticker}/holdings`);
  }

  getMarketSummary() {
    return this.request<any>('/api/v1/stocks/market-summary');
  }

  // Trading
  executeTrade(ticker: string, side: string, shares: number) {
    return this.request<any>('/api/v1/trades', {
      method: 'POST',
      body: JSON.stringify({ ticker, side, shares }),
    });
  }

  getTradeHistory(limit = 50, offset = 0) {
    return this.request<any[]>(`/api/v1/trades?limit=${limit}&offset=${offset}`);
  }

  // Portfolio
  getPortfolio() {
    return this.request<any>('/api/v1/portfolio');
  }

  // Leaderboard
  getLeaderboard(period: string = 'alltime') {
    return this.request<any[]>(`/api/v1/leaderboard?period=${period}`);
  }

  // Orders
  createOrder(ticker: string, side: string, orderType: string, shares: number, limitPrice?: string, stopPrice?: string) {
    return this.request<any>('/api/v1/orders', {
      method: 'POST',
      body: JSON.stringify({ ticker, side, order_type: orderType, shares, limit_price: limitPrice, stop_price: stopPrice }),
    });
  }

  getOrders() {
    return this.request<any[]>('/api/v1/orders');
  }

  cancelOrder(id: string) {
    return this.request<any>(`/api/v1/orders/${id}`, { method: 'DELETE' });
  }

  // Alerts
  createAlert(ticker: string, condition: string, targetPrice: string) {
    return this.request<any>('/api/v1/alerts', {
      method: 'POST',
      body: JSON.stringify({ ticker, condition, target_price: targetPrice }),
    });
  }

  getAlerts() {
    return this.request<any[]>('/api/v1/alerts');
  }

  deleteAlert(id: string) {
    return this.request<any>(`/api/v1/alerts/${id}`, { method: 'DELETE' });
  }

  // Achievements
  getAchievements() {
    return this.request<any[]>('/api/v1/achievements');
  }

  getMyAchievements() {
    return this.request<any[]>('/api/v1/achievements/me');
  }

  // Watchlist
  getWatchlist() {
    return this.request<string[]>('/api/v1/watchlist');
  }

  addToWatchlist(ticker: string) {
    return this.request<any>('/api/v1/watchlist', {
      method: 'POST',
      body: JSON.stringify({ ticker }),
    });
  }

  removeFromWatchlist(ticker: string) {
    return this.request<any>(`/api/v1/watchlist/${ticker}`, { method: 'DELETE' });
  }

  // Challenges
  getTodaysChallenge() {
    return this.request<any>('/api/v1/challenges/today');
  }

  checkChallenge() {
    return this.request<any>('/api/v1/challenges/check', { method: 'POST' });
  }

  claimChallenge(id: string) {
    return this.request<any>(`/api/v1/challenges/${id}/claim`, { method: 'POST' });
  }

  // Portfolio History
  getPortfolioHistory(limit = 100) {
    return this.request<any[]>(`/api/v1/portfolio/history?limit=${limit}`);
  }

  // Options
  getOptionChain(ticker: string, expiration?: string) {
    const params = expiration ? `?expiration=${expiration}` : '';
    return this.request<any>(`/api/v1/stocks/${ticker}/options${params}`);
  }

  getOptionExpirations(ticker: string) {
    return this.request<string[]>(`/api/v1/stocks/${ticker}/options/expirations`);
  }

  getOptionContract(id: string) {
    return this.request<any>(`/api/v1/options/${id}`);
  }

  executeOptionsTrade(contractId: string, side: string, quantity: number) {
    return this.request<any>('/api/v1/options/trades', {
      method: 'POST',
      body: JSON.stringify({ contract_id: contractId, side, quantity }),
    });
  }

  getOptionsTradeHistory(limit = 50, offset = 0) {
    return this.request<any[]>(`/api/v1/options/trades?limit=${limit}&offset=${offset}`);
  }

  getOptionsPositions() {
    return this.request<any[]>('/api/v1/options/positions');
  }

  createOptionsOrder(contractId: string, side: string, orderType: string, quantity: number, limitPrice?: string) {
    return this.request<any>('/api/v1/options/orders', {
      method: 'POST',
      body: JSON.stringify({ contract_id: contractId, side, order_type: orderType, quantity, limit_price: limitPrice }),
    });
  }

  getOptionsOrders() {
    return this.request<any[]>('/api/v1/options/orders');
  }

  cancelOptionsOrder(id: string) {
    return this.request<any>(`/api/v1/options/orders/${id}`, { method: 'DELETE' });
  }
}

export const apiClient = new APIClient();
