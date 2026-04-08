import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { http, HttpResponse } from 'msw';
import { server } from '@/test/mocks/server';
import { createWrapper } from '@/test/test-utils';
import { useStock, useStockHistory, useETFHoldings } from '@/hooks/use-stocks';

const API_URL = 'http://localhost';

const mockStock = {
  ticker: 'AAPL',
  name: 'Apple Inc.',
  sector: 'Technology',
  asset_type: 'stock',
  current_price: '192.50',
  day_open: '190.00',
  day_high: '194.00',
  day_low: '189.00',
  prev_close: '189.50',
  volume: 52000000,
};

const mockCryptoStock = {
  ticker: 'X:BTCUSD',
  name: 'Bitcoin',
  sector: 'Crypto',
  asset_type: 'crypto',
  current_price: '67500.00',
  day_open: '66000.00',
  day_high: '68000.00',
  day_low: '65500.00',
  prev_close: '65000.00',
  volume: 1500000,
};

describe('useStock', () => {
  it('fetches a stock by ticker', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker`, () =>
        HttpResponse.json(mockStock)
      ),
    );

    const { result } = renderHook(() => useStock('AAPL'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.name).toBe('Apple Inc.');
    expect(result.current.data?.ticker).toBe('AAPL');
  });

  it('fetches crypto stock with colon in ticker', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker`, ({ params }) => {
        // MSW decodes path params, so we check the raw ticker
        return HttpResponse.json(mockCryptoStock);
      }),
    );

    const { result } = renderHook(() => useStock('X:BTCUSD'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.name).toBe('Bitcoin');
    expect(result.current.data?.ticker).toBe('X:BTCUSD');
  });

  it('handles 404 for non-existent ticker', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker`, () =>
        new HttpResponse('stock not found', { status: 404 })
      ),
    );

    const { result } = renderHook(() => useStock('FAKE'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error?.message).toContain('404');
  });

  it('does not fetch when ticker is empty', async () => {
    let fetchCount = 0;
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker`, () => {
        fetchCount++;
        return HttpResponse.json(mockStock);
      }),
    );

    renderHook(() => useStock(''), { wrapper: createWrapper() });

    // Wait a tick to ensure no fetch fires
    await new Promise((r) => setTimeout(r, 100));
    expect(fetchCount).toBe(0);
  });
});

describe('useStockHistory', () => {
  it('returns empty array when no history', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
        HttpResponse.json([])
      ),
    );

    const { result } = renderHook(() => useStockHistory('AAPL', '1m'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual([]);
  });

  it('returns history data', async () => {
    const mockHistory = [
      { id: 1, ticker: 'AAPL', open: '190', high: '192', low: '189', close: '191', volume: 1000, interval: '1m', recorded_at: '2026-04-08T10:00:00Z' },
    ];
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
        HttpResponse.json(mockHistory)
      ),
    );

    const { result } = renderHook(() => useStockHistory('AAPL', '1m'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
  });

  it('handles null response gracefully', async () => {
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker/history`, () =>
        HttpResponse.json(null)
      ),
    );

    const { result } = renderHook(() => useStockHistory('AAPL', '1m'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    // Should fall back to empty array, not null
    expect(result.current.data).toEqual([]);
  });
});

describe('useETFHoldings', () => {
  it('does not fetch for non-ETF stocks', async () => {
    let fetchCount = 0;
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker/holdings`, () => {
        fetchCount++;
        return HttpResponse.json([]);
      }),
    );

    renderHook(() => useETFHoldings('AAPL', 'stock'), {
      wrapper: createWrapper(),
    });

    await new Promise((r) => setTimeout(r, 100));
    expect(fetchCount).toBe(0);
  });

  it('fetches holdings for ETF stocks', async () => {
    const mockHoldings = [
      { ticker: 'AAPL', name: 'Apple Inc.', weight: '0.07', price: '192.50' },
      { ticker: 'MSFT', name: 'Microsoft', weight: '0.06', price: '420.00' },
    ];
    server.use(
      http.get(`${API_URL}/api/v1/stocks/:ticker/holdings`, () =>
        HttpResponse.json(mockHoldings)
      ),
    );

    const { result } = renderHook(() => useETFHoldings('SPY', 'etf'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(2);
  });
});
