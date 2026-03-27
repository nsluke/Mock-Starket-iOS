export interface Portfolio {
  id: string;
  user_id: string;
  cash: string;
  net_worth: string;
  created_at: string;
  updated_at: string;
}

export interface Position {
  id: string;
  portfolio_id: string;
  ticker: string;
  shares: number;
  avg_cost: string;
  current_price: string;
  market_value: string;
  pnl: string;
  pnl_pct: string;
}

export interface PortfolioResponse {
  portfolio: Portfolio;
  positions: Position[];
  net_worth: string;
  invested: string;
}

export interface Trade {
  id: string;
  user_id: string;
  ticker: string;
  side: 'buy' | 'sell';
  shares: number;
  price: string;
  total: string;
  created_at: string;
}

export interface Order {
  id: string;
  ticker: string;
  side: 'buy' | 'sell';
  order_type: 'limit' | 'stop' | 'stop_limit';
  shares: number;
  limit_price?: string;
  stop_price?: string;
  status: string;
  created_at: string;
}
