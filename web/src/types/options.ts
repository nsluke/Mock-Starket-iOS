export interface OptionContract {
  id: string;
  ticker: string;
  option_type: 'call' | 'put';
  strike_price: string;
  expiration: string;
  contract_symbol: string;
  bid_price: string;
  ask_price: string;
  last_price: string;
  mark_price: string;
  open_interest: number;
  volume: number;
  implied_vol: string;
  delta: string;
  gamma: string;
  theta: string;
  vega: string;
  rho: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface OptionChainResponse {
  ticker: string;
  underlying_price: string;
  calls: OptionContract[];
  puts: OptionContract[];
}

export interface OptionPosition {
  id: string;
  portfolio_id: string;
  contract_id: string;
  quantity: number;
  avg_cost: string;
  collateral: string;
  contract: OptionContract;
  market_value: string;
  pnl: string;
  pnl_pct: string;
  is_long: boolean;
}

export interface OptionTrade {
  id: string;
  user_id: string;
  contract_id: string;
  side: 'buy_to_open' | 'buy_to_close' | 'sell_to_open' | 'sell_to_close';
  quantity: number;
  price: string;
  total: string;
  created_at: string;
}

export interface OptionOrder {
  id: string;
  user_id: string;
  contract_id: string;
  side: string;
  order_type: string;
  quantity: number;
  limit_price?: string;
  status: string;
  filled_price?: string;
  filled_at?: string;
  created_at: string;
  updated_at: string;
}
