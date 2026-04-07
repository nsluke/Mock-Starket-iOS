import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useAchievements } from '../use-achievements';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useAchievements', () => {
  it('fetches achievements and earned data together', async () => {
    const { result } = renderHook(() => useAchievements(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.achievements).toHaveLength(1);
    expect(result.current.data?.achievements[0].name).toBe('First Trade');
    expect(result.current.data?.earned).toHaveLength(1);
    expect(result.current.data?.earned[0].achievement_id).toBe('ach-1');
  });
});
