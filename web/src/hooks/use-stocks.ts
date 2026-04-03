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
      setStocks(data);
      return data;
    },
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
    queryFn: () => apiClient.getStockHistory(ticker, interval),
    enabled: !!ticker,
  });
}

export function useETFHoldings(ticker: string, assetType?: string) {
  return useQuery<ETFHolding[]>({
    queryKey: ['etf-holdings', ticker],
    queryFn: () => apiClient.getETFHoldings(ticker),
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
