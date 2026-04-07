import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useCurrentUser, useUpdateProfile } from '../use-user';
import { createWrapper } from '@/test/test-utils';
import { apiClient } from '@/lib/api-client';
import { mockUser } from '@/test/mocks/handlers';

beforeEach(() => {
  apiClient.setToken('test-token');
});

describe('useCurrentUser', () => {
  it('fetches current user profile', async () => {
    const { result } = renderHook(() => useCurrentUser(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.display_name).toBe(mockUser.display_name);
    expect(result.current.data?.is_guest).toBe(true);
    expect(result.current.data?.login_streak).toBe(5);
  });
});

describe('useUpdateProfile', () => {
  it('updates display name', async () => {
    const { result } = renderHook(() => useUpdateProfile(), { wrapper: createWrapper() });

    result.current.mutate('NewName');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});
