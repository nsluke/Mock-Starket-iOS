import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useAlerts, useCreateAlert, useDeleteAlert } from '../use-alerts';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useAlerts', () => {
  it('fetches alerts and separates active/triggered', async () => {
    const { result } = renderHook(() => useAlerts(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toHaveLength(2);
    const active = result.current.data!.filter((a) => !a.triggered);
    const triggered = result.current.data!.filter((a) => a.triggered);
    expect(active).toHaveLength(1);
    expect(triggered).toHaveLength(1);
  });
});

describe('useCreateAlert', () => {
  it('creates a price alert', async () => {
    const { result } = renderHook(() => useCreateAlert(), { wrapper: createWrapper() });

    result.current.mutate({ ticker: 'PIPE', condition: 'above', targetPrice: '160.00' });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe('useDeleteAlert', () => {
  it('deletes an alert', async () => {
    const { result } = renderHook(() => useDeleteAlert(), { wrapper: createWrapper() });

    result.current.mutate('alert-1');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});
