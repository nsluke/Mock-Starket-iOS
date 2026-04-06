import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useOrders, useCreateOrder, useCancelOrder } from '../use-orders';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useOrders', () => {
  it('fetches open orders', async () => {
    const { result } = renderHook(() => useOrders(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].ticker).toBe('PIPE');
    expect(result.current.data![0].order_type).toBe('limit');
  });
});

describe('useCreateOrder', () => {
  it('creates an order', async () => {
    const { result } = renderHook(() => useCreateOrder(), { wrapper: createWrapper() });

    result.current.mutate({
      ticker: 'PIPE',
      side: 'buy',
      orderType: 'limit',
      shares: 50,
      limitPrice: '145.00',
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe('useCancelOrder', () => {
  it('cancels an order', async () => {
    const { result } = renderHook(() => useCancelOrder(), { wrapper: createWrapper() });

    result.current.mutate('order-1');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});
