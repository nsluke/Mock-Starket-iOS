import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useTodaysChallenge, useCheckChallenge } from '../use-challenges';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useTodaysChallenge', () => {
  it('fetches todays challenge with progress', async () => {
    const { result } = renderHook(() => useTodaysChallenge(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.challenge?.challenge_type).toBe('trade_count');
    expect(result.current.data?.challenge?.description).toBe('Execute 3 trades today');
    expect(result.current.data?.progress?.completed).toBe(false);
    expect(result.current.data?.progress?.claimed).toBe(false);
  });
});

describe('useCheckChallenge', () => {
  it('checks challenge completion', async () => {
    const { result } = renderHook(() => useCheckChallenge(), { wrapper: createWrapper() });

    result.current.mutate();

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data.completed).toBe(false);
  });
});
