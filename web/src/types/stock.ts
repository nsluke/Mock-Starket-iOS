export interface Stock {
  ticker: string;
  name: string;
  sector: string;
  base_price: string;
  current_price: string;
  day_open: string;
  day_high: string;
  day_low: string;
  prev_close: string;
  volume: number;
  volatility: string;
  description?: string;
}

export interface PriceUpdate {
  ticker: string;
  price: string;
  change: string;
  change_pct: string;
  volume: number;
  high: string;
  low: string;
}

export interface PricePoint {
  id: number;
  ticker: string;
  price: string;
  open: string;
  high: string;
  low: string;
  close: string;
  volume: number;
  interval: string;
  recorded_at: string;
}

export interface MarketSummary {
  index_value: string;
  index_change_pct: string;
  total_stocks: number;
  gainers: number;
  losers: number;
}
