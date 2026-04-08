import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { useMarketStore } from '@/stores/market-store';
import type { Stock, PricePoint, ETFHolding, MarketSummary } from '@/types/stock';

export function useStocks() {
  const setStocks = useMarketStore((s) => s.setStocks);

  return useQuery<Stock[]>({
    queryKey: ['stocks'],
    queryFn: async () => {
      const data = await apiClient.getStocks();
      const stocks = data ?? [];
      setStocks(stocks);
      return stocks;
    },
    refetchInterval: 30_000, // Poll every 30s to keep prices fresh
  });
}

export function useStock(ticker: string) {
  return useQuery<Stock>({
    queryKey: ['stock', ticker],
    queryFn: () => apiClient.getStock(ticker),
    enabled: !!ticker,
  });
}

export function useStockHistory(ticker: string, interval: string) {
  return useQuery<PricePoint[]>({
    queryKey: ['stock-history', ticker, interval],
    queryFn: async () => {
      const data = await apiClient.getStockHistory(ticker, interval);
      return data ?? [];
    },
    enabled: !!ticker,
  });
}

export function useETFHoldings(ticker: string, assetType?: string) {
  return useQuery<ETFHolding[]>({
    queryKey: ['etf-holdings', ticker],
    queryFn: async () => {
      const data = await apiClient.getETFHoldings(ticker);
      return data ?? [];
    },
    enabled: !!ticker && assetType === 'etf',
  });
}

export function useMarketSummary() {
  const setSummary = useMarketStore((s) => s.setSummary);

  return useQuery<MarketSummary>({
    queryKey: ['market-summary'],
    queryFn: async () => {
      const data = await apiClient.getMarketSummary();
      setSummary(data);
      return data;
    },
  });
}

interface MarketStatusResponse {
  is_open: boolean;
  session: string;
  next_open: string;
  next_close: string;
}

export function useMarketStatus() {
  return useQuery<MarketStatusResponse>({
    queryKey: ['market-status'],
    queryFn: () => apiClient.request<MarketStatusResponse>('/api/v1/market/status'),
    refetchInterval: 60_000, // refresh every minute
  });
}
