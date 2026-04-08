import { describe, it, expect } from 'vitest';
import { displayTicker } from '@/types/stock';

describe('displayTicker', () => {
  it('returns normal stock tickers unchanged', () => {
    expect(displayTicker('AAPL')).toBe('AAPL');
    expect(displayTicker('MSFT')).toBe('MSFT');
    expect(displayTicker('GOOGL')).toBe('GOOGL');
  });

  it('strips X: prefix and USD suffix from crypto tickers', () => {
    expect(displayTicker('X:BTCUSD')).toBe('BTC');
    expect(displayTicker('X:ETHUSD')).toBe('ETH');
    expect(displayTicker('X:SOLUSD')).toBe('SOL');
    expect(displayTicker('X:DOGEUSD')).toBe('DOGE');
  });

  it('does not strip partial matches', () => {
    // Has X: but not USD suffix
    expect(displayTicker('X:BTCEUR')).toBe('X:BTCEUR');
    // Has USD suffix but no X: prefix
    expect(displayTicker('BTCUSD')).toBe('BTCUSD');
  });

  it('handles tickers with dots (BRK.B)', () => {
    expect(displayTicker('BRK.B')).toBe('BRK.B');
  });

  it('handles empty string', () => {
    expect(displayTicker('')).toBe('');
  });
});
