import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useLeaderboard } from '../use-leaderboard';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useLeaderboard', () => {
  it('fetches leaderboard entries', async () => {
    const { result } = renderHook(() => useLeaderboard('alltime'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toHaveLength(2);
    expect(result.current.data![0].rank).toBe(1);
    expect(result.current.data![0].display_name).toBe('TraderJoe');
  });
});
