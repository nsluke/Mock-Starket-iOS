import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { usePortfolio, useTradeHistory, useExecuteTrade } from '../use-portfolio';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('usePortfolio', () => {
  it('fetches portfolio with positions', async () => {
    const { result } = renderHook(() => usePortfolio(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.net_worth).toBe('105000.00');
    expect(result.current.data?.portfolio.cash).toBe('85000.00');
    expect(result.current.data?.positions).toHaveLength(1);
    expect(result.current.data?.positions[0].ticker).toBe('PIPE');
  });
});

describe('useTradeHistory', () => {
  it('fetches trade history', async () => {
    const { result } = renderHook(() => useTradeHistory(50, 0), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].ticker).toBe('PIPE');
    expect(result.current.data![0].side).toBe('buy');
  });
});

describe('useExecuteTrade', () => {
  it('executes a trade successfully', async () => {
    const { result } = renderHook(() => useExecuteTrade(), { wrapper: createWrapper() });

    result.current.mutate({ ticker: 'PIPE', side: 'buy', shares: 10 });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data.ticker).toBe('PIPE');
    expect(result.current.data.shares).toBe(10);
  });
});
