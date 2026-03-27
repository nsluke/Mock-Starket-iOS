import { create } from 'zustand';
import type { Stock, PriceUpdate, MarketSummary } from '@/types/stock';

interface MarketState {
  stocks: Stock[];
  summary: MarketSummary | null;
  isLoading: boolean;
  searchQuery: string;
  setStocks: (stocks: Stock[]) => void;
  setSummary: (summary: MarketSummary) => void;
  setLoading: (loading: boolean) => void;
  setSearchQuery: (query: string) => void;
  updatePrices: (updates: PriceUpdate[]) => void;
  filteredStocks: () => Stock[];
}

export const useMarketStore = create<MarketState>((set, get) => ({
  stocks: [],
  summary: null,
  isLoading: false,
  searchQuery: '',

  setStocks: (stocks) => set({ stocks }),
  setSummary: (summary) => set({ summary }),
  setLoading: (isLoading) => set({ isLoading }),
  setSearchQuery: (searchQuery) => set({ searchQuery }),

  updatePrices: (updates) => {
    set((state) => ({
      stocks: state.stocks.map((stock) => {
        const update = updates.find((u) => u.ticker === stock.ticker);
        if (!update) return stock;
        return {
          ...stock,
          current_price: update.price,
          day_high: update.high,
          day_low: update.low,
          volume: update.volume,
        };
      }),
    }));
  },

  filteredStocks: () => {
    const { stocks, searchQuery } = get();
    if (!searchQuery) return stocks;
    const q = searchQuery.toLowerCase();
    return stocks.filter(
      (s) =>
        s.ticker.toLowerCase().includes(q) ||
        s.name.toLowerCase().includes(q)
    );
  },
}));
