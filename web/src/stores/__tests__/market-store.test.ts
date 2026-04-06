import { describe, it, expect, beforeEach } from 'vitest';
import { useMarketStore } from '../market-store';

describe('useMarketStore', () => {
  beforeEach(() => {
    useMarketStore.setState({
      stocks: [],
      summary: null,
      isLoading: false,
      searchQuery: '',
    });
  });

  it('sets stocks', () => {
    const stocks = [
      { ticker: 'PIPE', name: 'Piper', sector: 'Tech', asset_type: 'stock', base_price: '100', current_price: '110', day_open: '100', day_high: '115', day_low: '99', prev_close: '99', volume: 1000, volatility: '0.02' },
    ];
    useMarketStore.getState().setStocks(stocks);
    expect(useMarketStore.getState().stocks).toEqual(stocks);
  });

  it('updates prices', () => {
    useMarketStore.setState({
      stocks: [
        { ticker: 'PIPE', name: 'Piper', sector: 'Tech', asset_type: 'stock', base_price: '100', current_price: '100', day_open: '100', day_high: '100', day_low: '100', prev_close: '99', volume: 1000, volatility: '0.02' },
      ],
    });

    useMarketStore.getState().updatePrices([
      { ticker: 'PIPE', price: '115.00', change: '15.00', change_pct: '15.00', volume: 2000, high: '120', low: '95' },
    ]);

    const stock = useMarketStore.getState().stocks[0];
    expect(stock.current_price).toBe('115.00');
    expect(stock.day_high).toBe('120');
    expect(stock.day_low).toBe('95');
    expect(stock.volume).toBe(2000);
  });

  it('filters stocks by search query', () => {
    useMarketStore.setState({
      stocks: [
        { ticker: 'PIPE', name: 'Piper Industries', sector: 'Tech', asset_type: 'stock', base_price: '100', current_price: '110', day_open: '100', day_high: '115', day_low: '99', prev_close: '99', volume: 1000, volatility: '0.02' },
        { ticker: 'BREW', name: 'BrewCraft', sector: 'Consumer', asset_type: 'stock', base_price: '50', current_price: '48', day_open: '50', day_high: '51', day_low: '47', prev_close: '49', volume: 500, volatility: '0.03' },
      ],
      searchQuery: 'pipe',
    });

    const filtered = useMarketStore.getState().filteredStocks();
    expect(filtered).toHaveLength(1);
    expect(filtered[0].ticker).toBe('PIPE');
  });

  it('returns all stocks when search is empty', () => {
    useMarketStore.setState({
      stocks: [
        { ticker: 'PIPE', name: 'Piper', sector: 'Tech', asset_type: 'stock', base_price: '100', current_price: '110', day_open: '100', day_high: '115', day_low: '99', prev_close: '99', volume: 1000, volatility: '0.02' },
        { ticker: 'BREW', name: 'BrewCraft', sector: 'Consumer', asset_type: 'stock', base_price: '50', current_price: '48', day_open: '50', day_high: '51', day_low: '47', prev_close: '49', volume: 500, volatility: '0.03' },
      ],
      searchQuery: '',
    });

    expect(useMarketStore.getState().filteredStocks()).toHaveLength(2);
  });

  it('searches by name too', () => {
    useMarketStore.setState({
      stocks: [
        { ticker: 'PIPE', name: 'Piper Industries', sector: 'Tech', asset_type: 'stock', base_price: '100', current_price: '110', day_open: '100', day_high: '115', day_low: '99', prev_close: '99', volume: 1000, volatility: '0.02' },
        { ticker: 'BREW', name: 'BrewCraft', sector: 'Consumer', asset_type: 'stock', base_price: '50', current_price: '48', day_open: '50', day_high: '51', day_low: '47', prev_close: '49', volume: 500, volatility: '0.03' },
      ],
      searchQuery: 'brew',
    });

    const filtered = useMarketStore.getState().filteredStocks();
    expect(filtered).toHaveLength(1);
    expect(filtered[0].ticker).toBe('BREW');
  });
});
