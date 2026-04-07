import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAuthStore } from '../auth-store';

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: vi.fn((key: string) => store[key] ?? null),
    setItem: vi.fn((key: string, value: string) => { store[key] = value; }),
    removeItem: vi.fn((key: string) => { delete store[key]; }),
    clear: vi.fn(() => { store = {}; }),
  };
})();
Object.defineProperty(global, 'localStorage', { value: localStorageMock });

describe('useAuthStore', () => {
  beforeEach(() => {
    useAuthStore.setState({ user: null, token: null, isLoading: true });
    localStorageMock.clear();
  });

  it('sets token', () => {
    useAuthStore.getState().setToken('test-token');
    expect(useAuthStore.getState().token).toBe('test-token');
  });

  it('sets user', () => {
    const user = {
      id: 'u1',
      firebase_uid: 'fb1',
      display_name: 'Test',
      is_guest: true,
      created_at: '2026-01-01',
      login_streak: 0,
      longest_streak: 0,
    };
    useAuthStore.getState().setUser(user);
    expect(useAuthStore.getState().user).toEqual(user);
  });

  it('signs out and clears localStorage', () => {
    useAuthStore.setState({ user: { id: 'u1', firebase_uid: 'fb1', display_name: 'Test', is_guest: true, created_at: '2026-01-01', login_streak: 0, longest_streak: 0 }, token: 'abc' });
    useAuthStore.getState().signOut();

    expect(useAuthStore.getState().user).toBeNull();
    expect(useAuthStore.getState().token).toBeNull();
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('mockstarket_token');
  });
});
