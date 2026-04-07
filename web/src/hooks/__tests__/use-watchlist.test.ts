import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useWatchlist, useAddToWatchlist, useRemoveFromWatchlist } from '../use-watchlist';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useWatchlist', () => {
  it('fetches watchlist tickers', async () => {
    const { result } = renderHook(() => useWatchlist(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(['PIPE', 'BREW']);
  });
});

describe('useAddToWatchlist', () => {
  it('adds a stock to watchlist', async () => {
    const { result } = renderHook(() => useAddToWatchlist(), { wrapper: createWrapper() });

    result.current.mutate('GLDX');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe('useRemoveFromWatchlist', () => {
  it('removes a stock from watchlist', async () => {
    const { result } = renderHook(() => useRemoveFromWatchlist(), { wrapper: createWrapper() });

    result.current.mutate('PIPE');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});
