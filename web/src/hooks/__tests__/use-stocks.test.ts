import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useStocks, useStock, useMarketSummary } from '../use-stocks';
import { createWrapper } from '@/test/test-utils';
import { mockStocks, mockMarketSummary } from '@/test/mocks/handlers';
import { apiClient } from '@/lib/api-client';
import { useMarketStore } from '@/stores/market-store';

beforeEach(() => {
  apiClient.setToken('test-token');
  useMarketStore.setState({ stocks: [], summary: null });
});

describe('useStocks', () => {
  it('fetches stocks and syncs to market store', async () => {
    const { result } = renderHook(() => useStocks(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toHaveLength(mockStocks.length);
    expect(result.current.data![0].ticker).toBe('PIPE');

    // Should sync to Zustand store
    expect(useMarketStore.getState().stocks).toHaveLength(mockStocks.length);
  });
});

describe('useStock', () => {
  it('fetches a single stock by ticker', async () => {
    const { result } = renderHook(() => useStock('PIPE'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.ticker).toBe('PIPE');
    expect(result.current.data?.name).toBe('Piper Industries');
  });
});

describe('useMarketSummary', () => {
  it('fetches market summary', async () => {
    const { result } = renderHook(() => useMarketSummary(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.total_stocks).toBe(mockMarketSummary.total_stocks);
    expect(result.current.data?.gainers).toBe(2);
    expect(result.current.data?.losers).toBe(1);
  });
});
